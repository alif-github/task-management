package usecase

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

func (t taskUsecase) Add(ctx context.Context, task *domain.TaskRequest) util.ErrorModel {
	//TODO implement me
	panic("implement me")
}

func (t taskUsecase) Update(ctx context.Context, task *domain.TaskRequest) util.ErrorModel {
	//TODO implement me
	panic("implement me")
}

func (t taskUsecase) Delete(ctx context.Context, task *domain.TaskRequest) util.ErrorModel {
	//TODO implement me
	panic("implement me")
}

func (t taskUsecase) Fetch(ctx context.Context, page, limit int) ([]domain.ListTaskResponse, util.ErrorModel) {
	//var db = serverconfig.ServerAttribute.DBConnection
	//TODO implement me
	panic("implement me")
}

func (t taskUsecase) GetByID(ctx context.Context, id int64) (domain.ViewTaskResponse, util.ErrorModel) {
	//TODO implement me
	panic("implement me")
}
