package controllers

import (
	"CedAppDeliver/app/models"

	"CedAppDeliver/utils"

	"fmt"

	"encoding/json"

	"github.com/revel/revel"
)

// Download ...
type Download struct {
	GormController
}

type ExtraConfigModel struct {
	Milestone bool
	AppID     string
	Comment   string
}

// Index ...
func (c Download) Index() revel.Result {
	var id = c.Params.Get("id")

	var ipa models.IPA
	Gdb.Where("random_id = ?", id).First(&ipa)
	// fmt.Println(ipa)

	var apk models.APK
	Gdb.Where("random_id = ?", id).First(&apk)
	// fmt.Println(apk)

	if ipa.RelativePlistURL != "" {
		var params = make(map[string]string)

		params["cerAddr"] = "public/myCer/WuYinJun-CA.crt"
		downloadURL := "https://" + utils.GetLocalIP() + "/" + ipa.RelativePlistURL
		params["downloadURL"] = downloadURL

		c.RenderArgs["model"] = ipa
		c.RenderArgs["params"] = params
		return c.RenderTemplate("Download/DownloadIPA.html")
	} else if apk.RelativeURL != "" {
		var params = make(map[string]string)

		params["cerAddr"] = "public/myCer/WuYinJun-CA.crt"
		downloadURL := "https://" + utils.GetLocalIP() + "/" + apk.RelativeURL
		params["downloadURL"] = downloadURL

		c.RenderArgs["model"] = apk
		c.RenderArgs["params"] = params
		return c.RenderTemplate("Download/DownloadAPK.html")
	} else {
		return c.NotFound("未找到此应用")
	}
}

func (c Download) ExtraConfig() revel.Result {
	data := c.Params.Form.Get("data")
	var model ExtraConfigModel

	json.Unmarshal([]byte(data), &model)

	fmt.Println(model)

	var ipa models.IPA
	Gdb.Where("random_id = ?", model.AppID).First(&ipa)
	// fmt.Println(ipa)

	var apk models.APK
	Gdb.Where("random_id = ?", model.AppID).First(&apk)
	// fmt.Println(apk)

	if ipa.RelativePlistURL != "" {
		ipa.IsMilestoneVersion = model.Milestone
		ipa.Comment = model.Comment
		Gdb.Save(&ipa)
	}
	if apk.RelativeURL != "" {
		apk.IsMilestoneVersion = model.Milestone
		apk.Comment = model.Comment
		Gdb.Save(&apk)
	}

	return c.Render()
}
