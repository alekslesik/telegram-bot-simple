package bot

import (
	"io"
	"log/slog"
	"strings"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type fakeBot struct {
	last tgbotapi.Chattable
}

func (f *fakeBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	f.last = c
	return tgbotapi.Message{}, nil
}

func newTestHandlers(bot Sender) Handlers {
	return Handlers{
		Bot:    bot,
		Logger: slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
	}
}

func TestHandlers_HandleMessage_Echo(t *testing.T) {
	fb := &fakeBot{}
	h := newTestHandlers(fb)

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 123},
		Text: "hello",
	}

	h.HandleMessage(msg)

	cfg, ok := fb.last.(tgbotapi.MessageConfig)
	if !ok {
		t.Fatalf("expected MessageConfig, got %T", fb.last)
	}

	if cfg.ChatID != msg.Chat.ID {
		t.Errorf("expected ChatID %d, got %d", msg.Chat.ID, cfg.ChatID)
	}
	if cfg.Text != "Ты написал: "+msg.Text {
		t.Errorf("unexpected text: %q", cfg.Text)
	}
}

func TestHandlers_HandleCommand_Start(t *testing.T) {
	fb := &fakeBot{}
	h := newTestHandlers(fb)

	text := "/start"
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 42},
		Text: text,
		Entities: []tgbotapi.MessageEntity{{
			Type:   "bot_command",
			Offset: 0,
			Length: len(text),
		}},
	}

	h.HandleCommand(msg)

	cfg, ok := fb.last.(tgbotapi.MessageConfig)
	if !ok {
		t.Fatalf("expected MessageConfig, got %T", fb.last)
	}
	if cfg.Text == "" {
		t.Fatalf("expected non-empty /start reply")
	}
	if _, ok := cfg.ReplyMarkup.(tgbotapi.ReplyKeyboardMarkup); !ok {
		t.Fatalf("expected reply keyboard markup for /start")
	}
}

func TestHandlers_HandleCommand_Echo_Args(t *testing.T) {
	fb := &fakeBot{}
	h := newTestHandlers(fb)

	text := "/echo hello"
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 77},
		Text: text,
		Entities: []tgbotapi.MessageEntity{{
			Type:   "bot_command",
			Offset: 0,
			Length: len("/echo"),
		}},
	}

	h.HandleCommand(msg)

	cfg, ok := fb.last.(tgbotapi.MessageConfig)
	if !ok {
		t.Fatalf("expected MessageConfig, got %T", fb.last)
	}
	if cfg.Text != "hello" {
		t.Fatalf("expected echoed args %q, got %q", "hello", cfg.Text)
	}
}

func TestHandlers_HandleMessage_ButtonMappedToCommand(t *testing.T) {
	fb := &fakeBot{}
	h := newTestHandlers(fb)

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{ID: 555},
		Text: "✅ Проверка статуса",
	}

	h.HandleMessage(msg)

	cfg, ok := fb.last.(tgbotapi.MessageConfig)
	if !ok {
		t.Fatalf("expected MessageConfig, got %T", fb.last)
	}
	if !strings.Contains(cfg.Text, "pong") {
		t.Fatalf("expected status response, got %q", cfg.Text)
	}
}
