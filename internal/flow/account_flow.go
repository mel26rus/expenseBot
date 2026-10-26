package flow

import (
	"context"
	"encoding/json"
	"errors"
	"expense-bot/internal/model"
	"expense-bot/internal/service"
	"expense-bot/internal/userstate"
	"fmt"
	"log/slog"
)

type AccountFlow struct {
	userService     *service.UserService
	accountService  *service.AccountService
	currencyService *service.CurrencyService
	sessionService  *service.SessionService
}

func NewAccountFlow(
	u *service.UserService,
	a *service.AccountService,
	c *service.CurrencyService,
	s *service.SessionService,
) *AccountFlow {
	return &AccountFlow{
		userService:     u,
		accountService:  a,
		currencyService: c,
		sessionService:  s,
	}
}

func (f *AccountFlow) Start(ctx context.Context, session model.Session, text string) (Response, error) {

	accs, _ := f.accountService.GetAccountsByUserID(ctx, session.UserID)

	if len(accs) == 0 {
		f.sessionService.Set(ctx, session.UserID, userstate.StateWaitingAccountName, nil)

		return Response{
			Text:          "У вас нет счетов. Введите название:",
			EditMessageId: session.EditMessageId,
		}, nil
	}

	currencies, err := f.currencyService.GetAll(ctx)
	if err != nil {
		slog.Error("Error getting currencies", "error", err)
		return Response{Text: "8_Ошибка получения валют"}, nil
	}

	payload := model.AccountPayload{
		Name: text,
	}

	f.sessionService.Set(ctx, session.UserID, userstate.StateWaitingCurrency, payload)

	return Response{
		Text:          fmt.Sprintf("Новый счет: %s\n Выберите валюту:", text),
		Keyboard:      buildCurrencyInline(currencies),
		EditMessageId: session.EditMessageId,
	}, nil
}

func (f *AccountFlow) HandleMessage(ctx context.Context, session model.Session, text string) (Response, error) {

	switch session.State {

	case userstate.StateWaitingAccountName:

		payload := model.AccountPayload{
			Name: text,
		}

		f.sessionService.Set(ctx, session.UserID, userstate.StateWaitingCurrency, payload)

		currencies, _ := f.currencyService.GetAll(ctx)

		return Response{
			Text:          "Выберите валюту:",
			Keyboard:      buildCurrencyInline(currencies),
			EditMessageId: session.EditMessageId,
		}, nil
	case userstate.StateWaitingCurrency:
		slog.Debug("Recived currency by message", "text", text)
		currencyCode, ok := parseCurrencyName(text)
		slog.Debug("Parsed currency from message", "currencyCode", currencyCode, "session.EditMessageId", session.EditMessageId)
		if !ok {
			f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
			return Response{
				Text:          "😔😔😔 Неверная валюта, начните заново.",
				EditMessageId: session.EditMessageId,
			}, nil
		}
		currencyId, err := f.currencyService.Create(ctx, currencyCode)
		if err != nil {
			slog.Error("Create currency", "Error", err)
			return f.currencyError(ctx, session, currencyCode, err), nil
		}
		var payload model.AccountPayload
		json.Unmarshal(session.Payload, &payload)
		// ++IDR
		accountId, err := f.accountService.CreateAccount(ctx, session.UserID, payload.Name, currencyId)
		if err != nil {
			slog.Error("Create account", "error", err)
			f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
			return Response{Text: "Не удалось создать счёт. Попробуйте ещё раз.", EditMessageId: session.EditMessageId}, nil
		}
		accountTitle := f.accountService.GetAccountTitle(ctx, accountId)
		f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
		return Response{
			Text:              fmt.Sprintf("✅ %s \n Счёт создан", accountTitle),
			EditMessageId:     session.EditMessageId,
			IsSendMenuMessage: true,
		}, nil

	}
	f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
	return Response{
		Text:          "😔😔😔 Неверная команда, начните заново.",
		EditMessageId: session.EditMessageId,
	}, nil
}

func (f *AccountFlow) HandleCallback(ctx context.Context, session model.Session, data string) (Response, error) {

	switch session.State {

	case userstate.StateWaitingCurrency:

		var payload model.AccountPayload
		json.Unmarshal(session.Payload, &payload)

		currencyID, command, _ := parseId(data)
		if command != "currency" {
			slog.Error("Wrong command", "command", command)
			return Response{Text: "Неверная команда"}, nil
		}
		currencies, err := f.currencyService.GetAll(ctx)
		if err != nil {
			return f.currencyError(ctx, session, "", err), nil
		}
		currencyCode := ""
		for _, currency := range currencies {
			if currency.ID == currencyID {
				currencyCode = currency.Code
				break
			}
		}
		if currencyCode == "" {
			return f.currencyError(ctx, session, "", service.ErrUnsupportedCurrency), nil
		}
		currencyID, err = f.currencyService.Create(ctx, currencyCode)
		if err != nil {
			return f.currencyError(ctx, session, currencyCode, err), nil
		}

		accoountId, err := f.accountService.CreateAccount(ctx, session.UserID, payload.Name, currencyID)
		if err != nil {
			slog.Error("Error creating account", "userID", session.UserID, "accountName", payload.Name, "currencyID", currencyID, "error", err)
			return Response{Text: "10_Ошибка создания счета"}, nil
		}

		accountTitle := f.accountService.GetAccountTitle(ctx, accoountId)

		f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)

		return Response{
			Text:              fmt.Sprintf("✅ %s \n Счёт создан", accountTitle),
			IsSendMenuMessage: true,
		}, nil
	}

	f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
	return Response{
		Text:              "11_Ошибка Неизвестный колбек",
		IsSendMenuMessage: true}, nil
}

func (f *AccountFlow) currencyError(ctx context.Context, session model.Session, code string, err error) Response {
	f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
	message := "😔😔😔 9_Ошибка создания валюты"
	if errors.Is(err, service.ErrUnsupportedCurrency) {
		message = fmt.Sprintf("Валюта %s не найдена в источнике курсов. Проверьте код или введите другую валюту", code)
	} else if errors.Is(err, service.ErrRatesUnavailable) {
		message = "Не удалось проверить валюту. Попробуйте ещё раз позже"
	}
	return Response{Text: message, EditMessageId: session.EditMessageId}
}
