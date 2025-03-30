package jobs

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"

	"github.com/IgorAleksandroff/delivery/internal/core/application/usecases"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
)

var _ cron.Job = &AssignOrdersJob{}

type AssignOrdersJob struct {
	assignOrdersCommandHandler *usecases.AssignOrdersCommandHandler
}

func NewAssignOrdersJob(
	assignOrdersCommandHandler *usecases.AssignOrdersCommandHandler) (*AssignOrdersJob, error) {
	if assignOrdersCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("moveCouriersCommandHandler")
	}

	return &AssignOrdersJob{
		assignOrdersCommandHandler: assignOrdersCommandHandler}, nil
}

func (j *AssignOrdersJob) Run() {
	ctx := context.Background()
	command, err := usecases.NewAssignOrdersCommand()
	if err != nil {
		log.Error(err)
	}
	err = j.assignOrdersCommandHandler.Handle(ctx, command)
	if err != nil {
		log.Error(err)
	}
}
