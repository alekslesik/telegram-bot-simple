package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

var logger = newLogger()

func newLogger() *slog.Logger {
	// Уровень логирования: LOG_LEVEL=debug включает debug, всё остальное — info.
	level := slog.LevelInfo
	if strings.EqualFold(strings.TrimSpace(os.Getenv("LOG_LEVEL")), "debug") {
		level = slog.LevelDebug
	}

	// Формат времени: 16.06.2026 15:31:30
	timeReplacer := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			if t, ok := a.Value.Any().(time.Time); ok {
				a.Value = slog.StringValue(t.Format("02.01.2006 15:04:05"))
			}
		}
		return a
	}

	opts := &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: timeReplacer,
	}

	// Формат логов: LOG_FORMAT=json → JSON, иначе человекочитаемый текст.
	format := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))
	if format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("env var TOKEN is not set (see .env)")
	}

	username := os.Getenv("USERNAME")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	if username != "" {
		logger.Info("authorized",
			"username", bot.Self.UserName,
			"expected_username", username,
		)
	} else {
		logger.Info("authorized",
			"username", bot.Self.UserName,
		)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Graceful shutdown: stop on SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	logger.Info("bot started with long polling, press Ctrl+C to stop")

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			handleMessage(bot, update.Message)

		case sig := <-stop:
			logger.Info("received signal, shutting down", "signal", sig.String())
			return
		}
	}
}

func handleMessage(bot Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	if msg.IsCommand() {
		handleCommand(bot, msg)
		return
	}

	// Simple echo for non-command messages
	reply := tgbotapi.NewMessage(chatID, "Ты написал: "+msg.Text)
	if _, err := bot.Send(reply); err != nil {
		logger.Error("failed to send message", "err", err)
	}
}

func handleCommand(bot Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	switch msg.Command() {
	case "start":
		text := "Привет! Я демонстрационный Telegram‑бот для бизнеса.\n\n" +
			"Я показываю, как может выглядеть живой продукт для заказчика:\n" +
			"- приветствие новых клиентов\n" +
			"- ответы на типовые вопросы\n" +
			"- сбор заявок прямо в чат\n" +
			"- простая обратная связь.\n\n" +
			"Напиши /help, чтобы увидеть, что я уже умею."
		reply := tgbotapi.NewMessage(chatID, text)
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send /start reply", "err", err)
		}

	case "help":
		text := "Я бот, который помогает автоматизировать общение с клиентами.\n\n" +
			"*Что я умею прямо сейчас:*\n" +
			"/start — приветствие и краткое объяснение\n" +
			"/help — это сообщение с возможностями\n" +
			"/about — чем полезен такой бот для бизнеса\n" +
			"/usecases — примеры задач, которые можно решить ботом\n" +
			"/ping — проверка, что бот онлайн\n" +
			"/echo <текст> — повторить ваш текст (пример простой команды)\n\n" +
			"Если просто написать сообщение — я отвечу тем же текстом. Это демонстрирует, как бот может принимать и обрабатывать любые обращения клиентов."
		reply := tgbotapi.NewMessage(chatID, text)
		reply.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send /help reply", "err", err)
		}

	case "about":
		text := "Этот бот — пример того, что вы можете получить как продукт.\n\n" +
			"Он подходит, если вам нужно:\n" +
			"- быстро отвечать клиентам 24/7\n" +
			"- разгрузить менеджеров от типовых вопросов\n" +
			"- собирать заявки и контакты прямо в Telegram\n" +
			"- аккуратно подводить людей к покупке или записи.\n\n" +
			"На основе этого бота можно добавить меню, оплату, интеграцию с CRM, базу знаний и любые сценарии под ваш бизнес."
		reply := tgbotapi.NewMessage(chatID, text)
		reply.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send /about reply", "err", err)
		}

	case "features":
		text := "*Какие возможности можно добавить в такого бота:*\n\n" +
			"- Меню с разделами (услуги, цены, контакты)\n" +
			"- Приём заявок: имя, телефон, комментарий → вам в чат или CRM\n" +
			"- Запись на услуги по времени (простое расписание)\n" +
			"- Опросы и быстрый сбор обратной связи\n" +
			"- Отправка файлов, инструкций, прайсов\n" +
			"- Ограниченный доступ по списку клиентов или ролям.\n\n" +
			"Текущая версия — минимальный живой пример. Все перечисленное можно добавить в этот же бот под ваши задачи."
		reply := tgbotapi.NewMessage(chatID, text)
		reply.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send /features reply", "err", err)
		}

	case "ping":
		reply := tgbotapi.NewMessage(chatID, "pong ✅ Бот запущен и готов работать с клиентами.")
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send /ping reply", "err", err)
		}

	case "echo":
		// /echo используется, если хочется явную команду вместо обычного сообщения
		args := strings.TrimSpace(msg.CommandArguments())
		if args == "" {
			reply := tgbotapi.NewMessage(chatID, "Использование: /echo <текст, который нужно повторить>")
			if _, err := bot.Send(reply); err != nil {
				logger.Error("failed to send /echo usage reply", "err", err)
			}
			return
		}
		reply := tgbotapi.NewMessage(chatID, args)
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send /echo reply", "err", err)
		}

	case "usecases":
		text := "*Примеры задач, для которых подходит такой бот:*\n\n" +
			"1. Салон / студия / услуги:\n" +
			"   — рассказать про услуги и цены\n" +
			"   — принять заявку или запись\n" +
			"   — отправить напоминание перед визитом\n\n" +
			"2. Онлайн‑курсы / эксперты:\n" +
			"   — выдать материалы и инструкции\n" +
			"   — собрать вопросы от учеников\n" +
			"   — аккуратно предлагать доп. продукты.\n\n" +
			"3. Малый бизнес:\n" +
			"   — ответы на частые вопросы\n" +
			"   — получение контакта для звонка\n" +
			"   — быстрые опросы клиентов.\n\n" +
			"Идея простая: всё, что менеджер делает руками в переписке, можно постепенно перенести в бота."
		reply := tgbotapi.NewMessage(chatID, text)
		reply.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send /usecases reply", "err", err)
		}

	default:
		reply := tgbotapi.NewMessage(chatID, "Неизвестная команда. Напиши /help, чтобы узнать, что я умею.")
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send unknown command reply", "err", err)
		}
	}
}
