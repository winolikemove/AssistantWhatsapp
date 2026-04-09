package whatsapp

import (
        "context"
        "math/rand"
        "strings"
        "time"
)

// Anti-ban delay constants (in milliseconds)
// Recommended: 3-7 seconds for human-like behavior
const (
        minTypingDelayMs = 3000  // 3 seconds (was 500ms)
        maxTypingDelayMs = 7000  // 7 seconds (was 1500ms)
)

// RandomTypingDelay returns an anti-ban friendly randomized typing delay.
// This mimics human typing behavior to avoid WhatsApp detection.
func RandomTypingDelay() time.Duration {
        ms := minTypingDelayMs + rand.Intn(maxTypingDelayMs-minTypingDelayMs+1)
        return time.Duration(ms) * time.Millisecond
}

// RandomHumanDelay returns a random delay for general anti-ban purposes.
// Use this before sending messages to appear more natural.
func RandomHumanDelay() time.Duration {
        // Random delay between 1-4 seconds
        ms := 1000 + rand.Intn(3000)
        return time.Duration(ms) * time.Millisecond
}

// SendPresenceBeforeReply emits typing presence before reply.
func SendPresenceBeforeReply(ctx context.Context, messenger Messenger, recipient string) {
        if messenger == nil {
                return
        }
        if strings.TrimSpace(recipient) == "" {
                return
        }

        // Best-effort presence; ignore error so chat flow still proceeds.
        _ = messenger.SendPresence(ctx, recipient)
}

// SendTextWithPresence sends typing presence first, then the text reply.
func SendTextWithPresence(ctx context.Context, messenger Messenger, recipient, text string) error {
        SendPresenceBeforeReply(ctx, messenger, recipient)
        return messenger.SendText(ctx, recipient, text)
}
