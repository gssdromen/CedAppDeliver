package controllers

import (
	"CedAppDeliver/app/models"

	"encoding/json"

	"github.com/revel/revel"
)

// AliasConfig ...
type AliasConfig struct {
	GormController
}

// AliasConfig 显示别名页面
func (c AliasConfig) AliasConfig() revel.Result {
	var aliasModels []models.AppAlias
	Gdb.Find(&aliasModels)
	logger.Println("AliasConfig 显示别名页面")
	for i, v := range aliasModels {
		logger.Printf("%d  %s\n", i, v)
	}
	return c.Render(aliasModels)
}

// SaveAliasConfig 保存别名设置
func (c AliasConfig) SaveAliasConfig() revel.Result {
	logger.Println("SaveAliasConfig 保存别名设置")
	var aliasModels []models.AppAlias
	Gdb.Delete(&aliasModels)

	var params = c.Params.Form
	var modelArray []models.AppAlias
	for k, _ := range params {
		json.Unmarshal([]byte(k), &modelArray)
	}

	for i, v := range modelArray {
		logger.Printf("%d  %s\n", i, v)
		Gdb.Save(&v)
	}

	return c.Render()
}
