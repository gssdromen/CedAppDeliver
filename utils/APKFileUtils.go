package utils

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/revel/revel"

	"CedAppDeliver/app/models"
)

// APKGetInfoFromString 从String中获取APK模型信息
func APKGetInfoFromString(apk *models.APK, str string) {
	list := strings.Split(str, "\n")
	for _, v := range list {
		// 解析版本信息
		if strings.HasPrefix(v, "package: ") {
			tempString := v[len("package: "):]
			list := strings.Split(tempString, " ")
			for _, m := range list {
				if strings.HasPrefix(m, "name") {
					apk.PackageName = TrimQuote(m[len("name="):])
					apk.AppID = getMD5([]byte(apk.PackageName))[:4]
				} else if strings.HasPrefix(m, "versionCode") {
					apk.ShortVersion = TrimQuote(m[len("versionCode="):])
				} else if strings.HasPrefix(m, "versionName") {
					apk.Version = TrimQuote(m[len("versionName="):])
				}
			}
		} else if strings.HasPrefix(v, "application-label") {
			// 解析应用名称
			if strings.Contains(v, "zh-CN") {
				list := strings.Split(v, ":")
				apk.DisplayName = TrimQuote(list[len(list)-1])
			}
		}
	}
	// 判断是否发现应用名称
	if apk.DisplayName == "" {
		apk.DisplayName = "Unknown"
	}

	rand := RandSeq(4)
	apk.RandomID = rand
}

// APKMoveFilesToPathForModel 移动文件到对应文件夹
func APKMoveFilesToPathForModel(apk *models.APK, apkPath string) (string, error) {
	timeStr := time.Now().Format("2006-01-02_15-04-05")
	folderPath := path.Join("public", "apk", apk.PackageName, timeStr)
	apk.RelativePath = path.Join(folderPath, apk.PackageName+".apk")

	url := "public/apk/" + apk.PackageName + "/" + timeStr + "/"
	apk.RelativeURL = url + apk.PackageName + ".apk"

	folderPath = path.Join(revel.BasePath, folderPath)

	if !IsFileExistForPath(folderPath) {
		err := os.MkdirAll(folderPath, 0777)
		if err != nil {
			return "", err
		}
	}

	// 移动APK
	newFilePath := path.Join(folderPath, apk.PackageName+".apk")
	err := os.Rename(apkPath, newFilePath)
	if err != nil {
		return newFilePath, err
	}

	return newFilePath, nil
}
