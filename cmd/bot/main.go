package main

import (
        "context"
        "encoding/base64"
        "fmt"
        "log"
        "os"
        "os/signal"
        "strings"
        "syscall"
        "time"

        "github.com/winolikemove/AssistantWhatsapp/internal/ai"
        "github.com/winolikemove/AssistantWhatsapp/internal/app"
        "github.com/winolikemove/AssistantWhatsapp/internal/commands"
        "github.com/winolikemove/AssistantWhatsapp/internal/config"
        "github.com/winolikemove/AssistantWhatsapp/internal/finance"
        "github.com/winolikemove/AssistantWhatsapp/internal/notes"
        "github.com/winolikemove/AssistantWhatsapp/internal/reminder"
        "github.com/winolikemove/AssistantWhatsapp/internal/sales"
        "github.com/winolikemove/AssistantWhatsapp/internal/sheets"
        "github.com/winolikemove/AssistantWhatsapp/internal/whatsapp"
        waLog "go.mau.fi/whatsmeow/util/log"
)

func verifyLLMConnectivity(ctx context.Context, llmClient *ai.LLMClient) error {
        checkCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
        defer cancel()

        const systemPrompt = "You are a healthcheck assistant. Reply with one short word only."
        resp, err := llmClient.Chat(checkCtx, systemPrompt, "reply with: READY")
        if err != nil {
                return fmt.Errorf("LLM connectivity check failed: %w", err)
        }
        if strings.TrimSpace(resp.Content) == "" && len(resp.ToolCalls) == 0 {
                return fmt.Errorf("LLM connectivity check failed: empty response")
        }
        return nil
}

func main() {
        // 0) Setup credentials from base64 (for cloud deployment)
        if credsBase64 := os.Getenv("GOOGLE_CREDENTIALS_BASE64"); credsBase64 != "" {
                credsJSON, err := base64.StdEncoding.DecodeString(credsBase64)
                if err != nil {
                        log.Fatalf("❌ Failed to decode GOOGLE_CREDENTIALS_BASE64: %v", err)
                }
                if err := os.WriteFile("credentials.json", credsJSON, 0644); err != nil {
                        log.Fatalf("❌ Failed to write credentials.json: %v", err)
                }
                log.Println("✅ Google credentials loaded from environment")
        }

        // 1) Load config
        cfg, err := config.Load()
        if err != nil {
                log.Fatalf("❌ Config error: %v", err)
        }
        log.Println("✅ Config loaded")

        // 2) Connect WhatsApp
        waLogger := waLog.Stdout("WhatsApp", "INFO", true)
        client, err := whatsapp.Connect(cfg.WASessionDBPath, waLogger)
        if err != nil {
                log.Fatalf("❌ WhatsApp connection error: %v", err)
        }
        log.Println("✅ WhatsApp connected")

        // 3) Initialize Google Sheets repository
        repo, err := sheets.NewGoogleSheetRepository(cfg.GoogleCredsPath, cfg.SheetsID)
        if err != nil {
                log.Fatalf("❌ Sheets error: %v", err)
        }
        log.Println("✅ Google Sheets connected")

        ctx := context.Background()

        // 4) Initialize core tabs
        if err := repo.InitDashboard(ctx); err != nil {
                log.Fatalf("❌ Failed to init Dashboard tab: %v", err)
        }
        if err := repo.InitBudgetTab(ctx); err != nil {
                log.Fatalf("❌ Failed to init Budget tab: %v", err)
        }
        if err := repo.InitNotesTab(ctx); err != nil {
                log.Fatalf("❌ Failed to init Notes tab: %v", err)
        }
        if err := repo.InitReminderTab(ctx); err != nil {
                log.Fatalf("❌ Failed to init Reminder tab: %v", err)
        }
        log.Println("✅ Sheets tabs initialized")

        // 5) Initialize sales configuration
        salesConfig := &sheets.SalesConfig{
                SupplierName:              cfg.SupplierName,
                SupplierPayDay:            cfg.SupplierPayDay,
                DefaultCreditDays:         14,
                ReminderTime:              cfg.ReminderTime,
                ReminderEnabled:           true,
                WhatsAppToCustomerEnabled: cfg.WAReminderToCust,
                WhatsAppReminderDays:      1,
        }

        // 6) Initialize services
        llmClient := ai.NewLLMClient(cfg.LLMBaseURL, cfg.LLMApiKey, cfg.LLMModel)
        financeService := finance.NewFinanceService(repo, llmClient)
        notesService := notes.NewNotesService(repo)
        messenger := whatsapp.NewWhatsAppClient(client)
        reminderService := reminder.NewService(repo, messenger, cfg.OwnerPhoneNumber)

        // 7) Initialize sales services
        salesRepo := sheets.NewSalesRepository(repo)
        salesService := sales.NewService(salesRepo, salesConfig)
        salesScheduler := sales.NewScheduler(salesService, salesConfig)

        // 8) Initialize sales tabs
        if err := salesRepo.InitSalesTabs(ctx); err != nil {
                log.Fatalf("❌ Failed to init Sales tabs: %v", err)
        }
        log.Println("✅ Sales tabs initialized")

        // 9) LLM preflight (fail fast)
        if err := verifyLLMConnectivity(ctx, llmClient); err != nil {
                log.Fatalf("❌ %v", err)
        }
        log.Println("✅ LLM preflight check passed")

        // 10) Start reminder scheduler
        if err := reminderService.Start(ctx); err != nil {
                log.Fatalf("❌ Failed to start reminder service: %v", err)
        }
        log.Println("✅ Reminder scheduler started")

        // 11) Start sales scheduler
        notifier := sales.NewWhatsAppNotifier(
                func(msg string) error { return messenger.SendText(ctx, cfg.OwnerPhoneNumber, msg) },
                nil, // customer sender not implemented yet
        )
        if err := salesScheduler.Start(ctx, notifier); err != nil {
                log.Fatalf("❌ Failed to start sales scheduler: %v", err)
        }
        log.Println("✅ Sales scheduler started")

        // 12) Command router
        cmdRouter := commands.NewRouter()
        cmdRouter.Register("/start", commands.StartHandler)
        cmdRouter.Register("/help", commands.HelpHandler)
        cmdRouter.Register("/menu", commands.MenuHandler)
        cmdRouter.Register("/kategori", commands.CategoryHandler)
        cmdRouter.Register("/export", commands.NewExportHandlerFactory(cfg.SheetsID).Handler)
        cmdRouter.Register("/laporan", commands.NewReportHandlerFactory(financeService).Handler)
        cmdRouter.Register("/budget", commands.NewBudgetHandlerFactory(financeService).Handler)
        cmdRouter.Register("/notes", commands.NewNotesHandlerFactory(notesService).Handler)
        cmdRouter.Register("/edit", commands.NewEditHandlerFactory(financeService).Handler)
        cmdRouter.Register("/hapus", commands.NewDeleteHandlerFactory(financeService).Handler)
        cmdRouter.Register("/reminder", commands.NewReminderHandlerFactory(reminderService).Handler)
        cmdRouter.Register("/done", commands.NewDoneHandlerFactory(reminderService).Handler)

        // 13) App router with sales support
        appRouter := app.NewAppRouter(cmdRouter, llmClient, financeService, notesService, reminderService, salesService)

        // 14) WhatsApp message handler registration
        handler := whatsapp.NewHandler(messenger, cfg.OwnerPhoneNumber, appRouter.HandleMessage)
        handler.Register(client)

        log.Println("✅ Bot is running! Waiting for messages...")
        log.Println("📦 Sales features enabled:")
        log.Printf("   - Supplier: %s", salesConfig.SupplierName)
        log.Printf("   - Pay day: tanggal %d setiap bulan", salesConfig.SupplierPayDay)
        log.Printf("   - Reminder time: %s", salesConfig.ReminderTime)
        log.Printf("   - WA to customer: %v", salesConfig.WhatsAppToCustomerEnabled)

        // 15) Graceful shutdown
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan

        log.Println("⏳ Shutting down...")
        salesScheduler.Stop()
        reminderService.Stop()
        client.Disconnect()
        log.Println("✅ Bot stopped. Goodbye!")
}
