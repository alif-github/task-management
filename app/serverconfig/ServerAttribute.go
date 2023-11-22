package serverconfig

import (
	"database/sql"
	"fmt"
	"github.com/alif-github/task-management/config"
	"github.com/alif-github/task-management/util"
	_ "github.com/jackc/pgx/stdlib"
	"os"
	"sync"
)

var (
	ServerAttribute serverAttribute
	once            sync.Once
)

type dbInfo struct {
	instance      *sql.DB
	driver        string
	connectionStr string
	setParams     []string
}

type serverAttribute struct {
	DBConnection *sql.DB
}

func SetServerAttribute() {
	var err error
	dbInfoDB := dbInfo{
		driver:        config.ApplicationConfiguration.GetPostgresql().Driver,
		connectionStr: config.ApplicationConfiguration.GetPostgresql().Address,
		setParams: []string{
			fmt.Sprintf(`search_path='%s'`, config.ApplicationConfiguration.GetPostgresql().DefaultSchema),
		},
	}

	defer func() {
		if err != nil {
			logModel := util.GenerateLogModel(config.ApplicationConfiguration.GetServer().Version, config.ApplicationConfiguration.GetServer().Application)
			logModel.Code = 500
			logModel.Message = err.Error()
			util.LogError(logModel.LoggerZapFieldObject())
			os.Exit(1)
		}
	}()

	once.Do(func() {
		if dbInfoDB.setParams != nil && len(dbInfoDB.setParams) > 0 {
			for _, param := range dbInfoDB.setParams {
				dbInfoDB.connectionStr += fmt.Sprintf(` %s`, param)
			}
		}

		dbInfoDB.instance, err = sql.Open(dbInfoDB.driver, dbInfoDB.connectionStr)
		if err != nil {
			return
		}

		ServerAttribute.DBConnection = dbInfoDB.instance
	})
}
