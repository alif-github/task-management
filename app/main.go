package main

import (
	"fmt"
	"github.com/alif-github/task-management/config"
	"github.com/alif-github/task-management/util"
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
}
