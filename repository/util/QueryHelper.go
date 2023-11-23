package util

import (
	"context"
	"database/sql"
	"github.com/alif-github/task-management/util"
)

type GetListParameterModel struct {
	Page   sql.NullInt64
	Limit  sql.NullInt64
	Filter []ListFilter
	Order  []sql.NullString
}

type ListFilter struct {
	Key      sql.NullString
	Operator sql.NullString
	Value    sql.NullString
}

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

func CountOffset(page, limit int) int {
	return (page - 1) * limit
}

func GetRows(rows *sql.Rows, wrap func(rows *sql.Rows) (interface{}, error)) (output []interface{}, errorModel util.ErrorModel) {
	var errs error
	if rows != nil {
		defer func() {
			errs = rows.Close()
			if errs != nil {
				return
			}
		}()
		for rows.Next() {
			temp, errs := wrap(rows)
			if errs != nil {
				errorModel = util.GenerateInternalDBServerError("QueryHelper.go", "GetRows", errs)
				return
			}
			output = append(output, temp)
		}
	} else {
		errorModel = util.GenerateInternalDBServerError("QueryHelper.go", "GetRows", errs)
		return
	}

	errorModel = util.GenerateNonError()
	return
}
