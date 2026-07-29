package flow

import (
	"context"
	"expense-bot/internal/service"
	"expense-bot/internal/userstate"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MainFlow struct {
	sessionService  *service.SessionService
	accountFlow     *AccountFlow
	transactionFlow *TransactionFlow
	MenuFlow        *MenuFlow
	ReportFlow      *ReportFlow
}

type Response struct {
	Text              string
	EditMessageId     int64
	Keyboard          *tgbotapi.InlineKeyboardMarkup
	IsSendMenuMessage bool
}

func NewMainFlow(
	session *service.SessionService,
	account *AccountFlow,
	tx *TransactionFlow,
	menuFlow *MenuFlow,
	reportFlow *ReportFlow,
) *MainFlow {
	return &MainFlow{
		sessionService:  session,
		accountFlow:     account,
		transactionFlow: tx,
		MenuFlow:        menuFlow,
		ReportFlow:      reportFlow,
	}
}

func (f *MainFlow) HandleMessage(ctx context.Context, tgUserID int64, text string) (Response, error) {

	user, _ := f.accountFlow.userService.GetOrCreate(ctx, tgUserID)
	session, _ := f.sessionService.GetUserSession(ctx, user.ID)
	slog.Debug("MainFlow.HandleMessage: Got message", "tgUserID", tgUserID, "text", text, "sessionState", session.State, "session.EditMessageId", session.EditMessageId)

	switch session.State {

	case userstate.StateIdle:
		slog.Debug("MainFlow.HandleMessage: StateIdle", "tgUserID", tgUserID)
		if amount, ok := parseAmount(text); ok {
			slog.Debug("MainFlow.HandleMessage: Parsed amount", "tgUserID", tgUserID, "amount", amount)
			return f.transactionFlow.Start(ctx, session, amount)
		}
		if text == "/start" {
			return f.GenerateFirstMessage()
		}
		slog.Debug("MainFlow.HandleMessage: StateIdle but not amount", "tgUserID", tgUserID, "text", text)
		return f.accountFlow.Start(ctx, session, text)

	default:
		slog.Debug("MainFlow.HandleMessage: State is default", "tgUserID", tgUserID, "sessionState", session.State)
		if isTransactionState(session.State) {
			return f.transactionFlow.HandleMessage(ctx, session, text)
		}
		return f.accountFlow.HandleMessage(ctx, session, text)
	}
}

func (f *MainFlow) HandleCallback(ctx context.Context, tgUserID int64, data string) (Response, error) {
	user, _ := f.accountFlow.userService.GetOrCreate(ctx, tgUserID)
	session, _ := f.sessionService.GetUserSession(ctx, user.ID)
	slog.Debug("MainFlow.HandleCallback: Got callback", "tgUserID", tgUserID, "data", data, "sessionState", session.State)
	if data == userstate.StateCancel {
		f.sessionService.Set(ctx, user.ID, userstate.StateIdle, nil)
		return f.GenerateFirstMessage()
	}
	if isTransactionState(session.State) {
		return f.transactionFlow.HandleCallback(ctx, session, data)
	}
	if isMenuFlow(data) {
		restext, _ := f.GenerateFirstMessage()
		res, err := f.MenuFlow.HandleCallback(ctx, session, data)
		res.Text = restext.Text
		return res, err
	}
	return f.accountFlow.HandleCallback(ctx, session, data)
}

// да да мне не нравится что тут tgUserID, но пока так, потом можно будет юзера по id доставать и юзать его tgId там, а не передавать его везде
func (f *MainFlow) SetUserSessionMessageId(ctx context.Context, tgUserID int64, editMessageId int64) (Response, error) {
	user, err := f.accountFlow.userService.GetOrCreate(ctx, tgUserID)
	if err != nil {
		return Response{Text: "6_Ошибка пользователя"}, err
	}
	err = f.sessionService.SetEditMessageId(ctx, user.ID, editMessageId)
	if err != nil {
		return Response{Text: "7_Ошибка при установке ID сообщения"}, err
	}
	slog.Debug("SetUserSessionMessageId", "user.ID", user.ID, "editMessageId", editMessageId)
	return Response{}, nil
}

func (f *MainFlow) GenerateFirstMessage() (Response, error) {
	messageText := `
💰 <b>Учёт личных финансов</b>
🏠 <b>Главное меню</b>

✏️ <b>Введите:</b>
• <code>250</code> — новая транзакция
• <code>Зарплатная карта</code> — создать новый счёт

💡 <i>Просто отправьте сумму или название нового счёта.</i>
`
	return Response{Text: messageText, Keyboard: f.buildMenuInline()}, nil

}

func (f *MainFlow) buildMenuInline() *tgbotapi.InlineKeyboardMarkup {
	ikb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📊 Меню", constDataMenuMain),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🔄 Настройки", constDataMenuSettings),
		),
	)
	return &ikb
}
