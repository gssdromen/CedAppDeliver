package controllers

import (
	"CedAppDeliver/app/models"

	"path"

	"github.com/jinzhu/gorm"
	// import sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/revel/revel"
)

// GormController ...
type GormController struct {
	BaseController
}

// Gdb 用此变量操作数据库
var Gdb *gorm.DB

// InitDB 初始化数据库
func InitDB() {
	var err error
	Gdb, err = gorm.Open("sqlite3", path.Join(revel.BasePath, "database.db"))
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	Gdb.AutoMigrate(&models.IPA{}, &models.APK{}, &models.AppAlias{})
}
