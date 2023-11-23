package task_usecase

import (
	"github.com/alif-github/task-management/domain"
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
