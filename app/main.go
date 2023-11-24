package main

import (
	"fmt"
	"github.com/alif-github/task-management/app/serverconfig"
	"github.com/alif-github/task-management/config"
	controller "github.com/alif-github/task-management/delivery/http"
	"github.com/alif-github/task-management/util"
	"github.com/gobuffalo/packr/v2"
	migrate "github.com/rubenv/sql-migrate"
	"os"
)

func main() {
	environment := "local"
	args := os.Args
	if len(args) > 1 {
		environment = args[1]
		fmt.Println("Application Run In Environment : ", environment)
	}

	config.GenerateConfiguration(environment)
	util.ConfigZap(config.ApplicationConfiguration.GetLogFile())
	serverconfig.SetServerAttribute()

	autoCreateSchema()
	dbMigration()

	logModel := util.GenerateLogModel(config.ApplicationConfiguration.GetServer().Version, config.ApplicationConfiguration.GetServer().Application)
	logModel.Message = fmt.Sprintf(`Starting Port %s`, config.ApplicationConfiguration.GetServer().Port)
	util.LogInfo(logModel.LoggerZapFieldObject())

	controller.Controller()
}

func autoCreateSchema() {
	createSchema := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s;`, config.ApplicationConfiguration.GetPostgresql().DefaultSchema)
	_, errS := serverconfig.ServerAttribute.DBConnection.Exec(createSchema)
	if errS != nil {
		util.LogError(util.DefaultGenerateLogModel(500, fmt.Sprintf(`Error auto create schema -> %s`, errS.Error())).LoggerZapFieldObject())
		os.Exit(3)
	}
}

func dbMigration() {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "./app/migration"),
	}

	if serverconfig.ServerAttribute.DBConnection != nil {
		n, err := migrate.Exec(serverconfig.ServerAttribute.DBConnection, "postgres", migrations, migrate.Up)
		if err != nil {
			util.LogError(util.DefaultGenerateLogModel(500, fmt.Sprintf(`Error on migration -> %s`, err.Error())).LoggerZapFieldObject())
			os.Exit(3)
		}

		util.LogInfo(util.DefaultGenerateLogModel(200, fmt.Sprintf(`Has Applied %d Migrations`, n)).LoggerZapFieldObject())
		return
	}

	util.LogError(util.DefaultGenerateLogModel(500, fmt.Sprintf(`DB Connection Not Found`)).LoggerZapFieldObject())
	os.Exit(3)
}
