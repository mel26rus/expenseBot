package bot

import (
	"context"
	"expense-bot/internal/flow"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot  *tgbotapi.BotAPI
	flow *flow.MainFlow
}

func NewHandler(bot *tgbotapi.BotAPI, flow *flow.MainFlow) *Handler {
	return &Handler{
		bot:  bot,
		flow: flow,
	}
}

func (h *Handler) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := h.bot.GetUpdatesChan(u)

	for {
		select {

		case <-ctx.Done():
			slog.Info("Stopping telegram handler...")
			h.bot.StopReceivingUpdates()
			return

		case update, ok := <-updates:
			if !ok {
				return
			}

			if update.CallbackQuery != nil {
				h.handleCallback(update.CallbackQuery)
			} else if update.Message != nil {
				h.handleMessage(update.Message)
			}
		}
	}
}

func (h *Handler) handleMessage(msg *tgbotapi.Message) {

	response, err := h.flow.HandleMessage(
		context.Background(),
		msg.From.ID,
		msg.Text,
	)
	if err != nil {
		slog.Error("MainFlow.handleMessage error", "error", err)
		response.Text = fmt.Sprintf("h.handleMessage Ошибка: %v", err)
	}
	slog.Debug("MainFlow.handleMessage success", "response.text", response.Text, "response.editmessageid", response.EditMessageId)

	if response.EditMessageId != 0 && !response.IsSendMenuMessage {
		editMessage := tgbotapi.NewEditMessageText(msg.Chat.ID, int(response.EditMessageId), response.Text)
		if response.Keyboard != nil && !response.IsSendMenuMessage {
			h.appendCancelButton(response.Keyboard)
			editMessage.ReplyMarkup = response.Keyboard
		}
		editMessage.ParseMode = "HTML"
		_, err := h.bot.Send(editMessage)
		if err != nil {
			slog.Error("MainFlow.handleMessage error_1", "error", err)
		}

	} else {
		editMessage := tgbotapi.NewEditMessageText(msg.Chat.ID, int(response.EditMessageId), response.Text)
		if response.Keyboard != nil && !response.IsSendMenuMessage {
			h.appendCancelButton(response.Keyboard)
			editMessage.ReplyMarkup = response.Keyboard
		}
		editMessage.ParseMode = "HTML"
		sendedMessage, err := h.bot.Send(editMessage)
		if err != nil {
			slog.Error("MainFlow.handleMessage error_2", "error", err)
		}
		h.flow.SetUserSessionMessageId(context.Background(), msg.Chat.ID, int64(sendedMessage.MessageID))
	}
	slog.Debug("MainFlow.handleMessage deleting message", "msg.Chat.ID", msg.Chat.ID, "msg.MessageID", msg.MessageID)
	delmsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	_, err = h.bot.Request(delmsg)
	if err != nil {
		slog.Error("MainFlow.handleMessage failed to delete message", "error", err)
	}

	slog.Debug("response.IsSendMenuMessage", "value", response.IsSendMenuMessage)
	if response.IsSendMenuMessage {
		res, _ := h.flow.GenerateFirstMessage()
		menuMessagee := tgbotapi.NewMessage(msg.Chat.ID, res.Text)
		menuMessagee.ReplyMarkup = res.Keyboard
		menuMessagee.ParseMode = "HTML"
		sendedMessage, err := h.bot.Send(menuMessagee)
		if err != nil {
			slog.Error("MainFlow.handleMessage error_2", "error", err)
		}
		slog.Debug("handlemessage", "response.IsSendMenuMessage", res.IsSendMenuMessage, "sendedMessage.MessageID", sendedMessage.MessageID)
		h.flow.SetUserSessionMessageId(context.Background(), msg.Chat.ID, int64(sendedMessage.MessageID))
	}

}

func (h *Handler) handleCallback(cb *tgbotapi.CallbackQuery) {

	slog.Debug("Handle callback", "userTgId", cb.From.ID, "callback data", cb.Data)

	response, err := h.flow.HandleCallback(
		context.Background(),
		cb.From.ID,
		cb.Data,
	)
	if err != nil {
		slog.Error("flow handle callback error", "error", err)
		response.Text = fmt.Sprintf("h.handleCallback Ошибка: %v", err)
	}
	slog.Debug("flow handle callback success", "response.text", response.Text, " cb.Message.MessageID", cb.Message.MessageID, "response.keyboard", response.Keyboard)
	h.flow.SetUserSessionMessageId(context.Background(), cb.From.ID, int64(cb.Message.MessageID))
	msg := tgbotapi.NewEditMessageText(cb.Message.Chat.ID, cb.Message.MessageID, response.Text)

	if response.Keyboard != nil && !response.IsSendMenuMessage {
		h.appendCancelButton(response.Keyboard)
		msg.ReplyMarkup = response.Keyboard
	}

	msg.ParseMode = "HTML"
	_, err = h.bot.Request(msg)
	if err != nil {
		slog.Error("failed to edit message", "error", err)
	}

	slog.Debug("response.IsSendMenuMessage", "value", response.IsSendMenuMessage, "EditMessageId", response.EditMessageId)
	if response.IsSendMenuMessage {
		res, _ := h.flow.GenerateFirstMessage()
		menuMessagee := tgbotapi.NewMessage(cb.Message.Chat.ID, res.Text)
		menuMessagee.ReplyMarkup = res.Keyboard
		menuMessagee.ParseMode = "HTML"
		sendedMessage, err := h.bot.Send(menuMessagee)
		if err != nil {
			slog.Error("MainFlow.handleMessage error_2", "error", err)
		}
		slog.Debug("handlecallback response.IsSendMenuMessage = true", "sendedMessage.MessageID", sendedMessage.MessageID)
		h.flow.SetUserSessionMessageId(context.Background(), cb.Message.Chat.ID, int64(sendedMessage.MessageID))
	}

}

func (h *Handler) appendCancelButton(keyboard *tgbotapi.InlineKeyboardMarkup) {
	cancelButton := tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "cancel")
	cancelRow := []tgbotapi.InlineKeyboardButton{cancelButton}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, cancelRow)
}
