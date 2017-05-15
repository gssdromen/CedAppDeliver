package controllers

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/revel/revel"
)

// BaseController ...
type BaseController struct {
	*revel.Controller
}

// it can be used for jobs
var logger *log.Logger

// InitLogger 初始化Logger
func InitLogger() {
	logFilePath := path.Join(revel.BasePath, "server.log")

	logfile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Printf("%s\r\n", err.Error())
		os.Exit(-1)
	}
	logger = log.New(logfile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)
}

func handleError(err error) {
	logger.Println(err.Error())
}
