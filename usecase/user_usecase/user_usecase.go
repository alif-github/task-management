package user_usecase

import (
	"context"
	"fmt"
	"github.com/alif-github/task-management/app/serverconfig"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/util"
	"time"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUsersUsecase(u domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (u userUsecase) Store(_ context.Context, users *domain.User) (err util.ErrorModel) {
	var (
		fileName = "user_usecase.go"
		funcName = "Store"
		db       = serverconfig.ServerAttribute.DBConnection
		timeNow  = time.Now()
	)

	existedUser, _ := u.userRepo.GetByName(nil, db, fmt.Sprintf(`%s %s`, users.FirstName, users.LastName))
	if existedUser != (domain.User{}) {
		err = util.GenerateConflictError(fileName, funcName)
		return
	}

	tx, errs := db.Begin()
	if errs != nil {
		return
	}

	defer func() {
		if err.Err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	users.Password, errs = util.HashPassword(users.Password)
	if errs != nil {
		err = util.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	users.CreatedBy = 1
	users.CreatedAt = timeNow
	users.UpdatedBy = 1
	users.UpdatedAt = timeNow

	err = u.userRepo.Store(nil, tx, users)
	if err.Err != nil {
		return
	}

	err = util.GenerateNonError()
	return
}

func (u userUsecase) Update(_ context.Context, users *domain.User) (err util.ErrorModel) {
	var (
		fileName = "user_usecase.go"
		funcName = "Update"
		db       = serverconfig.ServerAttribute.DBConnection
		timeNow  = time.Now()
	)

	existedUser, _ := u.userRepo.GetByID(nil, db, users.ID)
	if existedUser == (domain.User{}) {
		err = util.GenerateNotFoundError(fileName, funcName)
		return
	}

	tx, errs := db.Begin()
	if errs != nil {
		return
	}

	defer func() {
		if err.Err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	users.Password, errs = util.HashPassword(users.Password)
	if errs != nil {
		err = util.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	users.UpdatedBy = 1
	users.UpdatedAt = timeNow

	err = u.userRepo.Update(nil, tx, users)
	if err.Err != nil {
		_ = tx.Rollback()
		return
	}

	err = util.GenerateNonError()
	return
}

func (u userUsecase) Delete(_ context.Context, users *domain.User) (err util.ErrorModel) {
	var (
		fileName = "user_usecase.go"
		funcName = "Delete"
		db       = serverconfig.ServerAttribute.DBConnection
		timeNow  = time.Now()
	)

	existedUser, _ := u.userRepo.GetByID(nil, db, users.ID)
	if existedUser == (domain.User{}) {
		err = util.GenerateNotFoundError(fileName, funcName)
		return
	}

	tx, errs := db.Begin()
	if errs != nil {
		return
	}

	defer func() {
		if err.Err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	users.UpdatedBy = 1
	users.UpdatedAt = timeNow

	err = u.userRepo.Delete(nil, tx, users)
	if err.Err != nil {
		return
	}

	err = util.GenerateNonError()
	return
}

func (u userUsecase) Fetch(_ context.Context, page, limit int) (output []domain.User, err util.ErrorModel) {
	db := serverconfig.ServerAttribute.DBConnection
	output, err = u.userRepo.Fetch(nil, db, page, limit)
	if err.Err != nil {
		return
	}

	return
}

func (u userUsecase) GetByID(_ context.Context, id int64) (output domain.User, err util.ErrorModel) {
	db := serverconfig.ServerAttribute.DBConnection
	output, err = u.userRepo.GetByID(nil, db, id)
	if err.Err != nil {
		return
	}

	return
}

func (u userUsecase) GetLogin(_ context.Context, users *domain.User) (output domain.User, err util.ErrorModel) {
	db := serverconfig.ServerAttribute.DBConnection
	output, err = u.userRepo.GetLogin(nil, db, users)
	if err.Err != nil {
		return
	}

	return
}
