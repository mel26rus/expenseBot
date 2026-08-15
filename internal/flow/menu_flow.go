package flow

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/service"
	"log/slog"
	"strconv"
	"strings"
)

type MenuFlow struct {
	userService    *service.UserService
	accountService *service.AccountService
	reportFlow     *ReportFlow
}

func NewMenuFlow(
	u *service.UserService,
	a *service.AccountService,
	r *ReportFlow,
) *MenuFlow {
	return &MenuFlow{
		userService:    u,
		accountService: a,
		reportFlow:     r,
	}
}

func (m *MenuFlow) HandleCallback(ctx context.Context, session model.Session, data string) (Response, error) {
	slog.Debug("MenuFlow.HandleCallback:", "userId", session.UserID, "data", data)
	switch data {
	case constDataMenuSettings:
		userSetings, err := m.userService.GetUserSettings(ctx, session.UserID)
		return Response{Keyboard: buildMenuSettingsInline(userSetings.IsDailyReport, userSetings.IsMonthlyReport), IsSendMenuMessage: false}, err
	case constDataMenuMain:
		return Response{Keyboard: buildMenuMainInline()}, nil
	case constDailyReportChange:
		userSettings, err := m.userService.ChangeDailyReportConfig(ctx, session.UserID)
		return Response{Keyboard: buildMenuSettingsInline(userSettings.IsDailyReport, userSettings.IsMonthlyReport), IsSendMenuMessage: false}, err
	case constMonthlyReportChange:
		userSettings, err := m.userService.ChangeMonthlyReportConfig(ctx, session.UserID)
		return Response{Keyboard: buildMenuSettingsInline(userSettings.IsDailyReport, userSettings.IsMonthlyReport), IsSendMenuMessage: false}, err
	case constMenuReports:
		return Response{Keyboard: buildMenuReportsInline()}, nil
	case constMenuAccounts:
		accList, err := m.accountService.GetMenuAccountsByUserID(ctx, session.UserID)
		return Response{Keyboard: buildMenuUserAccountsInline(accList)}, err
	default:
		if strings.HasPrefix(data, constHideAccountChange) {
			accountIDStr := strings.TrimPrefix(data, constHideAccountChange)
			accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
			if err != nil {
				return Response{}, err
			}
			err = m.accountService.ChangeAccountIshidden(ctx, accountID)
			if err != nil {
				return Response{}, err
			}
			accList, err := m.accountService.GetMenuAccountsByUserID(ctx, session.UserID)
			return Response{Keyboard: buildMenuUserAccountsInline(accList)}, err
		}
		return Response{}, nil
	}
}
