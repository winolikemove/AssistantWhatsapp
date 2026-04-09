package whatsapp

import (
        "context"
        "fmt"
        "os"
        "path/filepath"
        "strings"

        "github.com/mdp/qrterminal/v3"
        "go.mau.fi/whatsmeow"
        "go.mau.fi/whatsmeow/store/sqlstore"
        waLog "go.mau.fi/whatsmeow/util/log"

        _ "github.com/mattn/go-sqlite3"
)

// LoginMethod defines how to authenticate with WhatsApp
type LoginMethod int

const (
        LoginMethodQR    LoginMethod = iota // Scan QR code (default)
        LoginMethodPairingCode               // Use pairing code (8-digit)
)

// LoginOptions contains configuration for WhatsApp login
type LoginOptions struct {
        Method      LoginMethod
        PhoneNumber string // Required for pairing code method (format: 628xxx)
}

// Connect initializes WhatsApp client with sqlite-backed session persistence.
// Default: Uses QR code authentication.
func Connect(dbPath string, log waLog.Logger) (*whatsmeow.Client, error) {
        return ConnectWithOptions(dbPath, log, LoginOptions{Method: LoginMethodQR})
}

// ConnectWithOptions initializes WhatsApp client with configurable login method.
// Supports both QR code and Pairing Code authentication.
//
// Pairing Code is recommended for better anti-ban protection.
// Note: Pairing code requires the phone number to be registered on WhatsApp.
func ConnectWithOptions(dbPath string, log waLog.Logger, opts LoginOptions) (*whatsmeow.Client, error) {
        if strings.TrimSpace(dbPath) == "" {
                return nil, fmt.Errorf("dbPath is required")
        }

        // Ensure parent directory exists.
        if dir := filepath.Dir(dbPath); dir != "." && dir != "" {
                if err := os.MkdirAll(dir, 0o755); err != nil {
                        return nil, fmt.Errorf("failed to create session db directory: %w", err)
                }
        }

        ctx := context.Background()

        container, err := sqlstore.New(
                ctx,
                "sqlite3",
                "file:"+dbPath+"?_foreign_keys=on&_busy_timeout=5000",
                log,
        )
        if err != nil {
                return nil, fmt.Errorf("failed to initialize sqlstore: %w", err)
        }

        deviceStore, err := container.GetFirstDevice(ctx)
        if err != nil {
                return nil, fmt.Errorf("failed to get device store: %w", err)
        }

        client := whatsmeow.NewClient(deviceStore, log)
        client.EnableAutoReconnect = true

        // First-time login requires authentication
        if client.Store.ID == nil {
                // Choose authentication method
                if opts.Method == LoginMethodPairingCode && opts.PhoneNumber != "" {
                        return loginWithPairingCode(ctx, client, opts.PhoneNumber)
                }
                // Default to QR code
                return loginWithQR(ctx, client)
        }

        // Existing session: connect directly.
        if err := client.Connect(); err != nil {
                return nil, fmt.Errorf("failed to connect with existing session: %w", err)
        }

        return client, nil
}

// loginWithQR handles QR code authentication
func loginWithQR(ctx context.Context, client *whatsmeow.Client) (*whatsmeow.Client, error) {
        qrChan, err := client.GetQRChannel(ctx)
        if err != nil {
                return nil, fmt.Errorf("failed to create QR channel: %w", err)
        }

        if err := client.Connect(); err != nil {
                return nil, fmt.Errorf("failed to connect WhatsApp client: %w", err)
        }

        fmt.Println("📱 Scan QR code berikut di WhatsApp (Linked Devices):")
        for evt := range qrChan {
                switch evt.Event {
                case "code":
                        qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
                        fmt.Println("Menunggu scan QR...")

                case "success":
                        fmt.Println("✅ Pairing berhasil.")
                        return client, nil

                case "timeout":
                        client.Disconnect()
                        return nil, fmt.Errorf("QR code timed out, please restart")

                case "error":
                        client.Disconnect()
                        if evt.Error != nil {
                                return nil, fmt.Errorf("QR pairing error: %w", evt.Error)
                        }
                        return nil, fmt.Errorf("QR pairing error")

                default:
                        // Non-success terminal events are treated as failures.
                        client.Disconnect()
                        if evt.Error != nil {
                                return nil, fmt.Errorf("pairing failed (%s): %w", evt.Event, evt.Error)
                        }
                        return nil, fmt.Errorf("pairing failed (%s)", evt.Event)
                }
        }

        client.Disconnect()
        return nil, fmt.Errorf("QR channel closed before pairing succeeded")
}

// loginWithPairingCode handles Pairing Code authentication
// This is the recommended method for better anti-ban protection
func loginWithPairingCode(ctx context.Context, client *whatsmeow.Client, phoneNumber string) (*whatsmeow.Client, error) {
        // Request pairing code
        pairingCode, err := client.PairPhone(phoneNumber, false, whatsmeow.PairClientDesktop, "AssistantWhatsapp")
        if err != nil {
                // Fallback to QR if pairing code fails
                fmt.Printf("⚠️  Pairing code gagal: %v\n", err)
                fmt.Println("📱 Beralih ke QR code...")
                return loginWithQR(ctx, client)
        }

        // Connect
        if err := client.Connect(); err != nil {
                return nil, fmt.Errorf("failed to connect WhatsApp client: %w", err)
        }

        fmt.Println("📱 Pairing Code Authentication")
        fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        fmt.Printf("🔑 Kode Pairing: %s\n", pairingCode)
        fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
        fmt.Println("")
        fmt.Println("Cara menggunakan:")
        fmt.Println("1. Buka WhatsApp di HP Anda")
        fmt.Println("2. Ketuk ⋮ → Perangkat tertaut → Tautkan perangkat")
        fmt.Println("3. Pilih 'Tautkan dengan nomor telepon'")
        fmt.Println("4. Masukkan kode di atas")
        fmt.Println("")

        // Keep connection alive and wait for events
        // The client will handle the pairing automatically
        // We just need to wait for Connected event or error
        for {
                select {
                case <-ctx.Done():
                        client.Disconnect()
                        return nil, fmt.Errorf("pairing timeout")
                default:
                        if client.IsConnected() && client.Store.ID != nil {
                                fmt.Println("✅ Pairing berhasil!")
                                return client, nil
                        }
                }
        }
}
