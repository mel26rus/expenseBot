package servicehost

import (
	"context"

	"expense-bot/internal/app"

	"github.com/kardianos/service"
)

type Program struct {
	App *app.App

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *Program) Start(s service.Service) error {

	p.ctx, p.cancel = context.WithCancel(context.Background())

	go p.App.Run(p.ctx)

	return nil
}

func (p *Program) Stop(s service.Service) error {

	if p.cancel != nil {
		p.cancel()
	}

	return nil
}
