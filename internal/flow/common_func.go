package flow

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/userstate"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
)

const constTxTypeExpense = "Расход"
const constTxTypeIncome = "Доход"
const constTxTypeTransfer = "Перевод"
const constDataMenuMain = "menu:main"
const constDataMenuSettings = "menu:settings"
const constDailyReportChange = "menu:dailyreport"
const constMonthlyReportChange = "menu:monthlyreport"
const constMenuReports = "menu:reports"
const constMenuTodayReport = "rep:todayrep"
const constMenuCurrMonthReport = "rep:curmonrep"
const constMenuAccounts = "menu:accounts"
const constHideAccountChange = "menu:acchide"
const ConstCancel = "cancel"

func formatAmount(a float64) string {
	if a == float64(int64(a)) {
		return fmt.Sprintf("%d", int64(a))
	}
	return fmt.Sprintf("%.2f", a)
}

func parseAmount(text string) (float64, bool) {

	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\u00A0", "")
	text = strings.ReplaceAll(text, "\u00a0", "")
	text = strings.ReplaceAll(text, ",", ".")

	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, false
	}

	return val, true
}

func parseId(text string) (int64, string, error) {
	parts := strings.Split(text, ":")
	if len(parts) != 2 {
		slog.Error("Invalid callback data format for comment", "data", text)
		return 0, "", fmt.Errorf("invalid callback data format")
	}

	txID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		slog.Error("Invalid transaction ID in callback data", "data", text, "error", err)
		return 0, "", fmt.Errorf("invalid transaction ID")
	}
	return txID, parts[0], nil
}

func isTransactionState(state string) bool {
	switch state {
	case userstate.StateWaitingAccount,
		userstate.StateWaitingType,
		userstate.StateWaitingComment,
		userstate.StateWaitingTransferAmountTo,
		userstate.StateWaitingTransferToAccount,
		userstate.StateCancel:
		return true
	}
	return false
}

func isMenuFlow(data string) bool {
	return strings.Contains(data, "menu:")
}

func isReportFlow(data string) bool {
	return strings.Contains(data, "rep:")
}

func parseCurrencyName(text string) (string, bool) {
	text = strings.TrimSpace(text)
	text = strings.ToUpper(text)

	// проверка: только латиница и цифры
	re := regexp.MustCompile(`^[A-Z0-9]+$`)
	if !re.MatchString(text) {
		return "", false
	}

	// ограничение длины
	if len(text) > 10 {
		text = text[:10]
	}

	return text, true
}

func (f *TransactionFlow) generateFinalTxMessage(ctx context.Context, payload model.TxPayload, balance float64) string {
	//💰💸
	icon := "📥"
	if payload.Amount < 0 {
		icon = "📤"
	}
	messageText := fmt.Sprintf(`
%s <b>Транзакция выполнена</b>
💸 <b>Сумма:</b> <code>%s</code>
🏦 <b>Счёт:</b> %s
📝 <b>Описание:</b> %s
💰 <b>Баланс:</b> <code>%s</code>
`,
		icon,
		formatAmount(payload.Amount),
		f.accountService.GetAccountTitle(ctx, payload.AccountID),
		payload.Comment,
		formatAmount(balance),
	)
	return messageText
}

func (f *TransactionFlow) generateFinalTxMessageWithTransfer(ctx context.Context, payload model.TxPayload) string {
	balanceFrom := formatAmount(f.accountService.GetAccountBalance(ctx, payload.AccountID))
	balanceTo := formatAmount(f.accountService.GetAccountBalance(ctx, payload.AccountToID))
	exchangeRate := "1"
	if payload.AmountTo < payload.Amount {
		exchangeRate = fmt.Sprintf("%.2f", payload.Amount/payload.AmountTo)
	} else if payload.AmountTo > payload.Amount {
		exchangeRate = fmt.Sprintf("%.2f", payload.AmountTo/payload.Amount)
	}

	messageText := fmt.Sprintf(
		`
🔁 <b>Перевод выполнен</b>
🏦 <code>%s</code> → <code>%s</code>
💸 <b>%s</b> → <b>%s</b>
💱 <b>Курс:</b> <code>%s</code>
💰 Балансы:
• %s: <code>%s</code>
• %s: <code>%s</code>
			`,
		f.accountService.GetAccountTitle(ctx, payload.AccountID),
		f.accountService.GetAccountTitle(ctx, payload.AccountToID),
		formatAmount(payload.Amount),
		formatAmount(payload.AmountTo),
		exchangeRate,
		f.accountService.GetAccountTitle(ctx, payload.AccountID),
		balanceFrom,
		f.accountService.GetAccountTitle(ctx, payload.AccountToID),
		balanceTo,
	)
	return messageText
}
