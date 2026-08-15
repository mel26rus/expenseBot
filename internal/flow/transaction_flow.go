package flow

import (
	"context"
	"encoding/json"
	"expense-bot/internal/model"
	"expense-bot/internal/service"
	"expense-bot/internal/userstate"
	"fmt"
	"log/slog"
	"math"

	"github.com/google/uuid"
)

type TransactionFlow struct {
	userService        *service.UserService
	accountService     *service.AccountService
	transactionService *service.TransactionService
	sessionService     *service.SessionService
}

func NewTransactionFlow(
	u *service.UserService,
	a *service.AccountService,
	t *service.TransactionService,
	s *service.SessionService,
) *TransactionFlow {
	return &TransactionFlow{
		userService:        u,
		accountService:     a,
		transactionService: t,
		sessionService:     s,
	}
}

func (f *TransactionFlow) Start(ctx context.Context, session model.Session, amount float64) (Response, error) {

	payload := model.TxPayload{
		Amount: amount,
	}

	f.sessionService.Set(ctx, session.UserID, userstate.StateWaitingAccount, payload)

	accs, _ := f.accountService.GetAccountsByUserID(ctx, session.UserID)

	// slog.Debug("Starting transaction flow", "userID", session.UserID, "amount", amount, "accountsdata",)

	return Response{
		Text:          fmt.Sprintf("Сумма: %s \nВыберите счет:", formatAmount(amount)),
		Keyboard:      buildUserAccountsInline(accs),
		EditMessageId: session.EditMessageId,
	}, nil
}

func (f *TransactionFlow) HandleMessage(ctx context.Context, session model.Session, text string) (Response, error) {

	var payload model.TxPayload
	json.Unmarshal(session.Payload, &payload)

	switch session.State {

	case userstate.StateWaitingTransferAmountTo:
		// ..получаем сумму зачисления на счет с дргого счета
		amount, ok := parseAmount(text)
		if !ok {
			slog.Debug("MainFlow.HandleMessage: Parsed amount", "UserID", session.UserID, "amount", text)
			return Response{
				Text:          fmt.Sprintf("Неверный формат числа '%s', повторите ввод", text),
				EditMessageId: session.EditMessageId,
			}, nil
		}

		payload.AmountTo = amount
		guid := uuid.New().String()
		//списываем
		f.transactionService.CreateTx(ctx, session.UserID, payload.AccountID, constTxTypeTransfer, -math.Abs(payload.Amount), constTxTypeTransfer, &guid)
		//пополняем
		f.transactionService.CreateTx(ctx, session.UserID, payload.AccountToID, constTxTypeTransfer, math.Abs(payload.AmountTo), constTxTypeTransfer, &guid)
		//если валюты одинаковые Комментарий не ждем просто пишем две транзакции и гуид группы
		f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, payload)

		messageText := f.generateFinalTxMessageWithTransfer(ctx, payload)

		return Response{
			Text:              messageText,
			Keyboard:          nil,
			EditMessageId:     session.EditMessageId,
			IsSendMenuMessage: true,
		}, nil

	case userstate.StateWaitingComment:
		payload.Comment = text

		txType := constTxTypeIncome
		if payload.ExpenseType < 0 {
			txType = constTxTypeExpense
		}

		f.transactionService.CreateTx(ctx, session.UserID, payload.AccountID, txType, payload.Amount, payload.Comment, nil)

		balance := f.accountService.GetAccountBalance(ctx, payload.AccountID)

		f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)

		slog.Debug("TransactionFlow.HandleMessage: Transaction created", "userId", session.UserID, "accountId", payload.AccountID, "txType", txType, "amount", payload.Amount, "comment", payload.Comment)

		messageText := f.generateFinalTxMessage(ctx, payload, balance)

		return Response{
			Text:              messageText,
			Keyboard:          nil,
			EditMessageId:     session.EditMessageId,
			IsSendMenuMessage: true,
		}, nil

	default:
		f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
		return Response{Text: "😔😔😔 Неизвестная команда, попробуйте снова"}, nil
	}
}

func (f *TransactionFlow) HandleCallback(ctx context.Context, session model.Session, data string) (Response, error) {

	slog.Debug("TransactionFlow.HandleCallback: Got callback", "UserID", session.UserID, "data", data)

	var payload model.TxPayload
	json.Unmarshal(session.Payload, &payload)

	switch session.State {

	case userstate.StateCancel:
		f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
		return Response{
			Text: "Действие отменено, начните заново",
		}, nil

	case userstate.StateWaitingAccount:
		slog.Debug("TransactionFlow.HandleCallback: StateWaitingAccount", "userId", session.UserID, "editMessageId", session.EditMessageId, "data", data)
		accountID, command, err := parseId(data)
		if command != "account" {
			slog.Error("Wrong inline command", "command", command)
			return Response{Text: "Неверная команда"}, nil
		}
		if err != nil {
			slog.Error("Error parsing callback data", "data", data, "error", err)
			return Response{Text: "1_Ошибка данных"}, nil
		}
		payload.AccountID = accountID

		accountTitle := f.accountService.GetAccountTitle(ctx, payload.AccountID)

		f.sessionService.Set(ctx, session.UserID, userstate.StateWaitingType, payload)

		slog.Debug("TransactionFlow.HandleCallback: StateWaitingType", "userId", session.UserID, "data", data)

		return Response{
			Text:          fmt.Sprintf("Сумма: %s \n Счет: %s \n Выберите тип операции:", formatAmount(payload.Amount), accountTitle),
			Keyboard:      buildTypeInline(),
			EditMessageId: session.EditMessageId,
		}, nil

	case userstate.StateWaitingType:
		slog.Debug("TransactionFlow.HandleCallback: StateWaitingType", "userId", session.UserID, "data", data)
		expenseType, command, _ := parseId(data)
		if command != "txtype" {
			slog.Error("Wrong expence callback data", "data", data)
			return Response{Text: "😱 Неверная команда"}, nil
		}
		payload.ExpenseType = expenseType

		if payload.ExpenseType == 0 { //значит перевод
			f.sessionService.Set(ctx, session.UserID, userstate.StateWaitingTransferToAccount, payload)
			accs, _ := f.accountService.GetAccountsByUserID(ctx, session.UserID)
			kb := buildUserAccountsInline(accs)
			return Response{
				Text:          "Выберите счёт зачисления:",
				Keyboard:      kb,
				EditMessageId: session.EditMessageId,
			}, nil
		} else if payload.ExpenseType < 0 {
			payload.Amount = -math.Abs(payload.Amount)
		} else {
			payload.Amount = math.Abs(payload.Amount)
		}

		accountTitle := f.accountService.GetAccountTitle(ctx, payload.AccountID)

		txType := constTxTypeIncome
		if payload.ExpenseType < 0 {
			txType = constTxTypeExpense
		}
		comments, _ := f.transactionService.GetLastComments(ctx, session.UserID, txType)

		f.sessionService.Set(ctx, session.UserID, userstate.StateWaitingComment, payload)

		return Response{
			Text:          fmt.Sprintf("Сумма: %s \n Счет: %s \n Тип операции: %s \n Выберите или напишите комментарий:", formatAmount(payload.Amount), accountTitle, txType),
			Keyboard:      buildCommentInline(comments),
			EditMessageId: session.EditMessageId,
		}, nil

	case userstate.StateWaitingTransferToAccount:

		accountToID, command, _ := parseId(data)
		if command != "account" {
			slog.Error("Wrong inline command", "command", command)
			return Response{Text: "Неверная команда"}, nil
		}
		if accountToID == payload.AccountID {
			slog.Error("Can't transfer to yourself")
			f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
			return Response{
				Text:          "Нельзя переводить на тот же самый счет, начните заново",
				EditMessageId: session.EditMessageId,
			}, nil
		}

		currencyFrom := f.accountService.GetAccountCurrencyByID(ctx, payload.AccountID)
		currencyTo := f.accountService.GetAccountCurrencyByID(ctx, accountToID)
		accountTitle := f.accountService.GetAccountTitle(ctx, payload.AccountID)
		accountToTitle := f.accountService.GetAccountTitle(ctx, accountToID)
		payload.AccountToID = accountToID
		//тут проверить валюты аккаунтов
		//разделить спросить сумму прихода если разные
		if currencyFrom != currencyTo {
			f.sessionService.Set(ctx, session.UserID, userstate.StateWaitingTransferAmountTo, payload)
			return Response{
				Text: fmt.Sprintf("Перевод.\nСумма списания со счета %s: %s\nСчет пополнения:%s\nВведите сумму пополнения в валюте счета: ",
					accountTitle, formatAmount(payload.Amount), accountToTitle),
				EditMessageId: session.EditMessageId,
				Keyboard:      nil,
			}, nil

		}
		guid := uuid.New().String()
		//списываем
		f.transactionService.CreateTx(ctx, session.UserID, payload.AccountID, constTxTypeTransfer, -math.Abs(payload.Amount), constTxTypeTransfer, &guid)
		//пополняем
		f.transactionService.CreateTx(ctx, session.UserID, payload.AccountToID, constTxTypeTransfer, math.Abs(payload.Amount), constTxTypeTransfer, &guid)
		//если валюты одинаковые Комментарий не ждем просто пишем две транзакции и гуид группы
		f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, payload)

		messageText := f.generateFinalTxMessage(ctx, payload, f.accountService.GetAccountBalance(ctx, payload.AccountID))

		return Response{
			Text:              messageText,
			Keyboard:          nil,
			EditMessageId:     session.EditMessageId,
			IsSendMenuMessage: true,
		}, nil

	case userstate.StateWaitingComment:

		payload.Comment = data
		txCommentId, txType, err := parseId(data)
		if err != nil {
			slog.Error("userstate.StateWaitingComment: Error parsing callback data", "data", data, "error", err)
			return Response{Text: "2_Ошибка получения данных"}, nil
		}

		txType = constTxTypeIncome
		if payload.ExpenseType < 0 {
			txType = constTxTypeExpense
		}

		payload.Comment, err = f.transactionService.GetCommentByTxID(ctx, txCommentId)
		if err != nil {
			slog.Error("userstate.StateWaitingComment: Error getting comment by txID", "txID", txCommentId, "error", err)
			return Response{Text: "3_Ошибка получения комментария"}, nil
		}

		err = f.transactionService.CreateTx(
			ctx,
			session.UserID,
			payload.AccountID,
			txType,
			payload.Amount,
			payload.Comment,
			nil,
		)

		if err != nil {
			slog.Error("TransactionFlow.HandleCallback: Error creating transaction", "userId", session.UserID, "error", err)
			return Response{
				Text:              "4_Ошибка при создании транзакции",
				IsSendMenuMessage: true,
			}, nil
		}

		balance := f.accountService.GetAccountBalance(ctx, payload.AccountID)

		f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)

		slog.Debug("TransactionFlow.HandleCallback: Transaction created", "userId", session.UserID, "accountId", payload.AccountID, "amount", payload.Amount, "comment", payload.Comment)

		messageText := f.generateFinalTxMessage(ctx, payload, balance)
		return Response{
			Text:              messageText,
			Keyboard:          nil,
			EditMessageId:     session.EditMessageId,
			IsSendMenuMessage: true,
		}, nil
	}

	f.sessionService.Set(ctx, session.UserID, userstate.StateIdle, nil)
	return Response{Text: "😱😔🙈 Неизвестная команда, попробуйте снова"}, nil
}
