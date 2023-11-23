package util

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound = errors.New(`your requested Item is not found`)
	ErrConflict = errors.New(`your Item already exist`)
)

type ErrorModel struct {
	Code     int
	FileName string
	FuncName string
	Err      error
}

func GenerateInternalDBServerError(fileName, funcName string, causedBy error) ErrorModel {
	return ErrorModel{
		Code:     http.StatusInternalServerError,
		FileName: fileName,
		FuncName: funcName,
		Err:      causedBy,
	}
}

func GenerateNotFoundError(fileName, funcName string) ErrorModel {
	return ErrorModel{
		Code:     http.StatusBadRequest,
		FileName: fileName,
		FuncName: funcName,
		Err:      ErrNotFound,
	}
}

func GenerateConflictError(fileName, funcName string) ErrorModel {
	return ErrorModel{
		Code:     http.StatusConflict,
		FileName: fileName,
		FuncName: funcName,
		Err:      ErrConflict,
	}
}

func GenerateNonError() ErrorModel {
	return ErrorModel{}
}
