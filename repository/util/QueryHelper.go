package util

import (
	"context"
	"database/sql"
	"github.com/alif-github/task-management/util"
)

func UpdateRow(_ context.Context, tx *sql.Tx, query string, params []interface{}, fileName, funcName string) util.ErrorModel {
	stmt, errs := tx.Prepare(query)
	if errs != nil {
		return util.GenerateInternalDBServerError(fileName, funcName, errs)
	}

	_, errs = stmt.Exec(params...)
	if errs != nil {
		return util.GenerateInternalDBServerError(fileName, funcName, errs)
	}

	return util.GenerateNonError()
}
