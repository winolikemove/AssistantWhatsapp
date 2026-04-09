package sales

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/verssache/AssistantWhatsapp/internal/sheets"
	"github.com/verssache/AssistantWhatsapp/pkg/formatter"
)

// Scheduler handles periodic reminder tasks for sales
type Scheduler struct {
	service *Service
	config  *sheets.SalesConfig
	cron    chan struct{}
}

// Notifier interface for sending messages (matches existing pattern)
type Notifier interface {
	SendMessageToOwner(msg string) error
	SendMessageToCustomer(customerID, msg string) error
}

// WhatsAppNotifier implements Notifier using WhatsApp client
type WhatsAppNotifier struct {
	sendToOwner    func(msg string) error
	sendToCustomer func(customerID, msg string) error
}

// NewWhatsAppNotifier creates a new WhatsApp notifier
func NewWhatsAppNotifier(ownerSender func(msg string) error, customerSender func(customerID, msg string) error) *WhatsAppNotifier {
	return &WhatsAppNotifier{
		sendToOwner:    ownerSender,
		sendToCustomer: customerSender,
	}
}

func (n *WhatsAppNotifier) SendMessageToOwner(msg string) error {
	if n.sendToOwner != nil {
		return n.sendToOwner(msg)
	}
	return nil
}

func (n *WhatsAppNotifier) SendMessageToCustomer(customerID, msg string) error {
	if n.sendToCustomer != nil {
		return n.sendToCustomer(customerID, msg)
	}
	return nil
}

// NewScheduler creates a new sales scheduler
func NewScheduler(service *Service, config *sheets.SalesConfig) *Scheduler {
	if config == nil {
		config = sheets.DefaultSalesConfig()
	}
	return &Scheduler{
		service: service,
		config:  config,
		cron:    make(chan struct{}),
	}
}

// Start begins the scheduler loop
func (s *Scheduler) Start(ctx context.Context, notifier Notifier) error {
	if s == nil {
		return fmt.Errorf("scheduler is nil")
	}

	go s.run(ctx, notifier)
	log.Println("✅ Sales reminder scheduler started")
	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	if s != nil && s.cron != nil {
		close(s.cron)
	}
}

func (s *Scheduler) run(ctx context.Context, notifier Notifier) {
	// Calculate next reminder time
	for {
		now := time.Now().In(sheets.WIB)

		// Parse reminder time (default 08:00)
		hour, minute := 8, 0
		if parts := strings.Split(s.config.ReminderTime, ":"); len(parts) == 2 {
			fmt.Sscanf(parts[0], "%d", &hour)
			fmt.Sscanf(parts[1], "%d", &minute)
		}

		// Calculate next trigger time
		nextTrigger := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
		if nextTrigger.Before(now) {
			nextTrigger = nextTrigger.Add(24 * time.Hour)
		}

		duration := nextTrigger.Sub(now)

		select {
		case <-ctx.Done():
			return
		case <-s.cron:
			return
		case <-time.After(duration):
			if s.config.ReminderEnabled {
				s.sendDailyReminder(ctx, notifier)
			}
			// Also check for supplier reminder on 22nd of each month
			if now.Day() == 22 {
				s.sendSupplierReminder(ctx, notifier)
			}
		}
	}
}

func (s *Scheduler) sendDailyReminder(ctx context.Context, notifier Notifier) {
	if s.service == nil || notifier == nil {
		return
	}

	summary, err := s.service.GetReceivableSummary(ctx)
	if err != nil {
		log.Printf("Error getting receivable summary: %v", err)
		return
	}

	// Build reminder message
	now := time.Now().In(sheets.WIB)
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("🌅 *PENGINGAT HARIAN - %s*\n\n", now.Format("02/01/2006")))

	// Due today
	if len(summary.DueToday) > 0 {
		msg.WriteString("📌 *PIUTANG JATUH TEMPO HARI INI:*\n")
		for _, r := range summary.DueToday {
			daysLate := int(time.Since(r.JatuhTempo).Hours() / 24)
			msg.WriteString(fmt.Sprintf("├─ %s: %s (%s)\n",
				r.CustomerNama, formatter.FormatIDR(r.Jumlah), r.ItemNama))
		}
		msg.WriteString(fmt.Sprintf("   *Total: %s*\n\n", formatter.FormatIDR(summary.TotalDueToday)))
	}

	// Overdue
	if len(summary.Overdue) > 0 {
		msg.WriteString("⚠️ *PIUTANG OVERDUE (TERLAMBAT):*\n")
		for _, r := range summary.Overdue {
			daysLate := int(time.Since(r.JatuhTempo).Hours() / 24)
			msg.WriteString(fmt.Sprintf("├─ %s: %s (telat %d hari)\n",
				r.CustomerNama, formatter.FormatIDR(r.Jumlah), daysLate))
		}
		msg.WriteString(fmt.Sprintf("   *Total: %s*\n\n", formatter.FormatIDR(summary.TotalOverdue)))
	}

	// Total
	msg.WriteString(fmt.Sprintf("💰 *TOTAL PIUTANG: %s*", formatter.FormatIDR(summary.TotalPending)))

	// Send to owner
	if err := notifier.SendMessageToOwner(msg.String()); err != nil {
		log.Printf("Error sending reminder to owner: %v", err)
	}

	// Send reminder to customers if enabled
	if s.config.WhatsAppToCustomerEnabled {
		for _, r := range summary.DueToday {
			customerMsg := fmt.Sprintf(
				"📊 *PENGINGAT PEMBAYARAN*\n\nHalo %s 👋\n\nIni pengingat bahwa pembayaran Anda sebesar %s untuk %s jatuh tempo pada %s.\n\nTerima kasih! 🙏",
				r.CustomerNama,
				formatter.FormatIDR(r.Jumlah),
				r.ItemNama,
				r.JatuhTempo.Format("02/01/2006"),
			)
			if err := notifier.SendMessageToCustomer(r.CustomerID, customerMsg); err != nil {
				log.Printf("Error sending reminder to customer %s: %v", r.CustomerNama, err)
			}
		}
	}
}

func (s *Scheduler) sendSupplierReminder(ctx context.Context, notifier Notifier) {
	if s.service == nil || notifier == nil {
		return
	}

	summary, err := s.service.GetPayableSummary(ctx)
	if err != nil {
		log.Printf("Error getting payable summary: %v", err)
		return
	}

	if summary.TotalPending <= 0 {
		return
	}

	now := time.Now().In(sheets.WIB)
	msg := fmt.Sprintf(
		"📅 *PENGINGAT HUTANG SUPPLIER*\n\nHutang bulan %s:\n%s (dari %d transaksi)\n\nJatuh tempo pembayaran: %s\nSupplier: %s\n\nKetik \"bayar hutang [jumlah]\" setelah pembayaran.",
		now.Format("January 2006"),
		formatter.FormatIDR(summary.TotalPending),
		len(summary.Payables),
		summary.DueDate.Format("02/01/2006"),
		s.config.SupplierName,
	)

	if err := notifier.SendMessageToOwner(msg); err != nil {
		log.Printf("Error sending supplier reminder: %v", err)
	}
}

// GetReminderTime returns configured reminder time
func (s *Scheduler) GetReminderTime() string {
	if s == nil || s.config == nil {
		return "08:00"
	}
	return s.config.ReminderTime
}

// IsEnabled returns whether reminders are enabled
func (s *Scheduler) IsEnabled() bool {
	if s == nil || s.config == nil {
		return false
	}
	return s.config.ReminderEnabled
}
