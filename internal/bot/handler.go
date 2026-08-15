package bot

import (
	"context"
	"expense-bot/internal/dateutil"
	"expense-bot/internal/flow"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	bot  *tgbotapi.BotAPI
	flow *flow.MainFlow
}

const timeLayout = "2006-01-02 15:04:05"

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
	const fn_name = "Handler.handleMessage"

	response, err := h.flow.HandleMessage(
		context.Background(),
		msg.From.ID,
		msg.Text,
	)
	if err != nil {
		slog.Error(fn_name, "error", err)
		response.Text = fmt.Sprintf("h.handleMessage Ошибка: %v", err)
	}
	slog.Debug(fn_name, "response.text", response.Text, "response.editmessageid", response.EditMessageId)

	if response.EditMessageId != 0 && !response.IsSendMenuMessage && !response.IsMainMenu {
		editMessage := tgbotapi.NewEditMessageText(msg.Chat.ID, int(response.EditMessageId), response.Text)
		if response.Keyboard != nil && !response.IsMainMenu {
			h.appendCancelButton(response.Keyboard)
			editMessage.ReplyMarkup = response.Keyboard
		}
		editMessage.ParseMode = tgbotapi.ModeHTML
		_, err := h.bot.Send(editMessage)
		if err != nil {
			slog.Error(fmt.Sprintf("%s error_1", fn_name), "error", err)
			newMessage := tgbotapi.NewMessage(editMessage.ChatID, editMessage.Text)
			if editMessage.ReplyMarkup != nil {
				newMessage.ReplyMarkup = editMessage.ReplyMarkup
			}
			newMessage.ParseMode = tgbotapi.ModeHTML
			_, err = h.bot.Send(newMessage)
			if err != nil {
				slog.Error(fmt.Sprintf("%s error_1_1", fn_name), "Error", err)
			}
		}

	} else {
		editMessage := tgbotapi.NewEditMessageText(msg.Chat.ID, int(response.EditMessageId), response.Text)
		if response.Keyboard != nil && !response.IsMainMenu {
			h.appendCancelButton(response.Keyboard)
			editMessage.ReplyMarkup = response.Keyboard
		}
		editMessage.ParseMode = tgbotapi.ModeHTML
		sendedMessage, err := h.bot.Send(editMessage)
		if err != nil {
			slog.Error(fn_name+" error_2", "error", err)
			newMessage := tgbotapi.NewMessage(msg.Chat.ID, response.Text)
			if editMessage.ReplyMarkup != nil {
				newMessage.ReplyMarkup = editMessage.ReplyMarkup
			}
			newMessage.ParseMode = tgbotapi.ModeHTML
			sendedMessage, err = h.bot.Send(newMessage)
			if err != nil {
				slog.Error(fn_name+" error_2_1", "error", err)
			}
		}
		slog.Debug("SendedMessage", "sendedMessageId", sendedMessage.MessageID)
		h.flow.SetUserSessionMessageId(context.Background(), msg.Chat.ID, int64(sendedMessage.MessageID))
	}
	slog.Debug(fn_name+" deleting message", "msg.Chat.ID", msg.Chat.ID, "msg.MessageID", msg.MessageID)
	delmsg := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	_, err = h.bot.Request(delmsg)
	if err != nil {
		slog.Error(fn_name+" failed to delete message", "error", err)
	}

	slog.Debug("response.IsSendMenuMessage", "value", response.IsSendMenuMessage)
	if response.IsSendMenuMessage {
		res, _ := h.flow.GenerateFirstMessage()
		menuMessagee := tgbotapi.NewMessage(msg.Chat.ID, res.Text)
		menuMessagee.ReplyMarkup = res.Keyboard
		menuMessagee.ParseMode = tgbotapi.ModeHTML
		sendedMessage, err := h.bot.Send(menuMessagee)
		if err != nil {
			slog.Error(fn_name+" error_3", "error", err)
		}
		slog.Debug("handlemessage", "response.IsSendMenuMessage", res.IsSendMenuMessage, "sendedMessage.MessageID", sendedMessage.MessageID)
		h.flow.SetUserSessionMessageId(context.Background(), msg.Chat.ID, int64(sendedMessage.MessageID))
	}

}

func (h *Handler) handleCallback(cb *tgbotapi.CallbackQuery) {

	slog.Debug("Handle callback", "userTgId", cb.From.ID, "callback data", cb.Data, "callback cb.Message.Chat.ID", cb.Message.Chat.ID)

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
	if response.Keyboard != nil && !response.IsMainMenu {
		h.appendCancelButton(response.Keyboard)
	}
	msg.ReplyMarkup = response.Keyboard
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
			slog.Error("MainFlow.handleMessage error_4", "error", err)
		}
		slog.Debug("handlecallback response.IsSendMenuMessage = true", "sendedMessage.MessageID", sendedMessage.MessageID)
		h.flow.SetUserSessionMessageId(context.Background(), cb.Message.Chat.ID, int64(sendedMessage.MessageID))
	}

}

func (h *Handler) HandleDailyReports(ctx context.Context) {
	start, end := dateutil.Yesterday()
	slog.Debug("HandleDailyReports", "start", start.Format(timeLayout), "end", end.Format(timeLayout))
	h.sendReports(ctx, start, end)
}
func (h *Handler) HandleMonthlyReports(ctx context.Context) {
	start, end := dateutil.PreviousMonth()
	slog.Debug("HandleMaonthlyReports", "start", start.Format(timeLayout), "end", end.Format(timeLayout))
	h.sendReports(ctx, start, end)
}

func (h *Handler) sendReports(ctx context.Context, start time.Time, end time.Time) {
	users, err := h.flow.ReportFlow.GetExistsTxUser(ctx, start, end)
	if err != nil {
		slog.Error("HandleDailyReports", "Error", err)
	}

	for _, user := range users {
		var res flow.Response
		res, err = h.flow.ReportFlow.BuildUserReport(ctx, user.ID, start, end)
		slog.Debug("SendReports_1 Deleting message", "user.ID", user.ID, "EditMessageId", res.EditMessageId)
		delMessage := tgbotapi.NewDeleteMessage(user.TelegramID, int(res.EditMessageId))
		_, err = h.bot.Request(delMessage)
		if err != nil {
			slog.Error("HandleDailyReports.Request delMessage", "Error", err)
		}
		msg := tgbotapi.NewMessage(user.TelegramID, res.Text)
		msg.ParseMode = tgbotapi.ModeHTML
		_, err = h.bot.Send(msg)
		if err != nil {
			slog.Error("HandleDailyReports.Send newMessage", "Error", err)
		}
		if res.IsSendMenuMessage {
			res, _ := h.flow.GenerateFirstMessage()
			menuMessage := tgbotapi.NewMessage(user.TelegramID, res.Text)
			menuMessage.ReplyMarkup = res.Keyboard
			menuMessage.ParseMode = tgbotapi.ModeHTML
			sendedMessage, err := h.bot.Send(menuMessage)
			if err != nil {
				slog.Error("MainFlow.SendReports_2 error_4", "error", err)
			}
			slog.Debug("HandleDailyReports end", "response.IsSendMenuMessage", res.IsSendMenuMessage, "sendedMessage.MessageID", sendedMessage.MessageID)
			h.flow.SetUserSessionMessageId(context.Background(), user.TelegramID, int64(sendedMessage.MessageID))
		}

	}
}

func (h *Handler) appendCancelButton(keyboard *tgbotapi.InlineKeyboardMarkup) {
	cancelButton := tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "cancel")
	cancelRow := []tgbotapi.InlineKeyboardButton{cancelButton}
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, cancelRow)
}
