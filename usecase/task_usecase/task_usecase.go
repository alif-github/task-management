package task_usecase

import (
	"context"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/util"
	"time"
)

type taskUsecase struct {
	taskRepo       domain.TaskRepository
	contextTimeout time.Duration
}

func NewTaskUsecase(t domain.TaskRepository, timeout time.Duration) domain.TaskUsecase {
	return &taskUsecase{
		taskRepo:       t,
		contextTimeout: timeout,
	}
}

func (t taskUsecase) Fetch(ctx context.Context, task *domain.TaskRequest) ([]domain.ListTaskResponse, util.ErrorModel) {
	//var db = serverconfig.ServerAttribute.DBConnection
	//TODO implement me
	panic("implement me")
}
