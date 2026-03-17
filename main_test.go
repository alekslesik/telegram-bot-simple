package main

import (
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

func TestHandleMessage_Echo(t *testing.T) {
	bot := &fakeBot{}

	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 123,
		},
		Text: "hello",
	}

	handleMessage(bot, msg)

	if bot.last == nil {
		t.Fatalf("expected bot to send a message, got nil")
	}

	cfg, ok := bot.last.(tgbotapi.MessageConfig)
	if !ok {
		t.Fatalf("expected MessageConfig, got %T", bot.last)
	}

	if cfg.ChatID != msg.Chat.ID {
		t.Errorf("expected ChatID %d, got %d", msg.Chat.ID, cfg.ChatID)
	}

	expectedText := "Ты написал: " + msg.Text
	if cfg.Text != expectedText {
		t.Errorf("expected text %q, got %q", expectedText, cfg.Text)
	}
}

func TestHandleCommand_Start(t *testing.T) {
	bot := &fakeBot{}

	text := "/start"
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 42,
		},
		Text: text,
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len(text),
			},
		},
	}

	handleCommand(bot, msg)

	cfg, ok := bot.last.(tgbotapi.MessageConfig)
	if !ok {
		t.Fatalf("expected MessageConfig, got %T", bot.last)
	}

	if cfg.ChatID != msg.Chat.ID {
		t.Errorf("expected ChatID %d, got %d", msg.Chat.ID, cfg.ChatID)
	}

	if cfg.Text == "" {
		t.Errorf("expected non-empty reply text for /start")
	}
}

func TestHandleCommand_Help(t *testing.T) {
	bot := &fakeBot{}

	text := "/help"
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 77,
		},
		Text: text,
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len(text),
			},
		},
	}

	handleCommand(bot, msg)

	cfg, ok := bot.last.(tgbotapi.MessageConfig)
	if !ok {
		t.Fatalf("expected MessageConfig, got %T", bot.last)
	}

	if cfg.ChatID != msg.Chat.ID {
		t.Errorf("expected ChatID %d, got %d", msg.Chat.ID, cfg.ChatID)
	}

	if cfg.Text == "" {
		t.Errorf("expected non-empty reply text for /help")
	}
}

func TestHandleCommand_Unknown(t *testing.T) {
	bot := &fakeBot{}

	text := "/unknown"
	msg := &tgbotapi.Message{
		Chat: &tgbotapi.Chat{
			ID: 99,
		},
		Text: text,
		Entities: []tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len(text),
			},
		},
	}

	handleCommand(bot, msg)

	cfg, ok := bot.last.(tgbotapi.MessageConfig)
	if !ok {
		t.Fatalf("expected MessageConfig, got %T", bot.last)
	}

	if cfg.ChatID != msg.Chat.ID {
		t.Errorf("expected ChatID %d, got %d", msg.Chat.ID, cfg.ChatID)
	}

	if cfg.Text == "" {
		t.Errorf("expected non-empty text for unknown command")
	}
}

