package controllers

import (
	"CedAppDeliver/app/models"
	"CedAppDeliver/utils"

	"github.com/revel/revel"
)

// App ...
type App struct {
	GormController
}

// AllForID ...
func (c App) AllForID() revel.Result {
	var id = c.Params.Get("id")
	var alias models.AppAlias
	Gdb.First(&alias, "app_alias = ?", id)
	if alias.AppID != "" {
		id = alias.AppID
	}

	var params = make(map[string]string)

	var ipas []models.IPA
	Gdb.Order("created_at desc").Where(models.IPA{AppID: id}).Find(&ipas)

	if len(ipas) > 0 {
		params["DisplayName"] = ipas[0].DisplayName

		params["cerAddr"] = "public/myCer/WuYinJun-CA.crt"
		downloadURL := "https://" + utils.GetLocalIP() + "/" + ipas[0].RelativePlistURL
		params["downloadURL"] = downloadURL

		c.RenderArgs["apps"] = ipas
		c.RenderArgs["model"] = ipas[0]
		c.RenderArgs["params"] = params
	}

	return c.RenderTemplate("App/AllForID.html")
}

// ShowForID ...
func (c App) ShowForID() revel.Result {
	var id = c.Params.Get("id")

	var ipas []models.IPA
	Gdb.Order("created_at desc").Where(models.IPA{AppID: id}).Find(&ipas)
	var apks []models.APK
	Gdb.Order("created_at desc").Where(models.APK{AppID: id}).Find(&apks)

	var params = make(map[string]string)

	if len(ipas) > 0 {
		params["DisplayName"] = ipas[0].DisplayName
		c.RenderArgs["apps"] = ipas
		c.RenderArgs["params"] = params
		return c.RenderTemplate("App/ShowForIDIPA.html")
	} else if len(apks) > 0 {
		params["DisplayName"] = apks[0].DisplayName
		c.RenderArgs["apps"] = apks
		c.RenderArgs["params"] = params
		return c.RenderTemplate("App/ShowForIDAPK.html")
	}

	return c.NotFound("未找到应用")
}

// Show ...
func (c App) Show() revel.Result {
	var ipas []models.IPA
	var apks []models.APK

	// 读取IPA
	rows, err := Gdb.Raw("SELECT DISTINCT app_id, display_name, bundle_identifier FROM ipas", "").Rows()
	defer rows.Close()
	if err != nil {
		handleError(err)
		return c.RenderError(err)
	}

	for rows.Next() {
		var ipa models.IPA
		rows.Scan(&ipa.AppID, &ipa.DisplayName, &ipa.BundleIdentifier)
		ipas = append(ipas, ipa)
	}

	// 读取APK
	rows, err = Gdb.Raw("SELECT DISTINCT app_id, display_name, package_name FROM apks", "").Rows()
	defer rows.Close()
	if err != nil {
		handleError(err)
		return c.RenderError(err)
	}

	for rows.Next() {
		var apk models.APK
		rows.Scan(&apk.AppID, &apk.DisplayName, &apk.PackageName)
		apks = append(apks, apk)
	}

	return c.Render(ipas, apks)
}
