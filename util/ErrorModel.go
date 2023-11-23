package util

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound  = errors.New(`your requested Item is not found`)
	ErrConflict  = errors.New(`your Item already exist`)
	ErrForbidden = errors.New(`access denied, forbidden`)
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

func GenerateUnknownServerError(fileName, funcName string, causedBy error) ErrorModel {
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

func GenerateForbiddenError(fileName, funcName string) ErrorModel {
	return ErrorModel{
		Code:     http.StatusForbidden,
		FileName: fileName,
		FuncName: funcName,
		Err:      ErrForbidden,
	}
}

func GenerateNonError() ErrorModel {
	return ErrorModel{}
}
