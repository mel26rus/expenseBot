package flow

import (
	"context"
	"expense-bot/internal/model"
	"expense-bot/internal/service"
	"log/slog"
)

type MenuFlow struct {
	userService    *service.UserService
	accountService *service.AccountService
}

func NewMenuFlow(
	u *service.UserService,
	a *service.AccountService,
) *MenuFlow {
	return &MenuFlow{
		userService:    u,
		accountService: a,
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
		accList, err := m.accountService.GetAccountsByUserID(ctx, session.UserID)
		return Response{Keyboard: buildMenuUserAccountsInline(accList)}, err
	default:
		return Response{}, nil
	}
}
