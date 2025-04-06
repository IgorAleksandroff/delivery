package jobs

import (
	"context"
	"github.com/IgorAleksandroff/delivery/internal/core/usecases/commands"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrdersJob{}

type AssignOrdersJob struct {
	assignOrdersCommandHandler *commands.AssignOrdersCommandHandler
}

func NewAssignOrdersJob(
	assignOrdersCommandHandler *commands.AssignOrdersCommandHandler) (*AssignOrdersJob, error) {
	if assignOrdersCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("moveCouriersCommandHandler")
	}

	return &AssignOrdersJob{
		assignOrdersCommandHandler: assignOrdersCommandHandler}, nil
}

func (j *AssignOrdersJob) Run() {
	ctx := context.Background()
	command, err := commands.NewAssignOrdersCommand()
	if err != nil {
		log.Error(err)
	}
	err = j.assignOrdersCommandHandler.Handle(ctx, command)
	if err != nil {
		log.Error(err)
	}
}
