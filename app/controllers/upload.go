package controllers

import (
	"CedAppDeliver/app/models"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"strings"

	"CedAppDeliver/utils"

	"fmt"

	"github.com/revel/revel"
)

// Upload ...
type Upload struct {
	GormController
}

// Index ...
func (c Upload) Index() revel.Result {
	return c.Render()
}

func handleAPK(fsrc multipart.File) (models.APK, error) {
	var model models.APK
	tempFileName := utils.RandSeq(5) + ".apk"
	folderPath := path.Join(revel.BasePath, "public", "apk")
	if !utils.IsFileExistForPath(folderPath) {
		os.MkdirAll(folderPath, 0777)
	}
	fdstPath := path.Join(folderPath, tempFileName)
	fdst, err := os.OpenFile(fdstPath, os.O_CREATE|os.O_WRONLY, 0777)
	defer fdst.Close()
	defer fsrc.Close()
	if err != nil {
		handleError(err)
		return model, nil
	}
	// Write file field from file to upload
	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		handleError(err)
		return model, nil
	}

	// 解析apk中的信息
	cmd := exec.Command("aapt", "dump", "badging", fdstPath)
	infoString, err := utils.RunShellCommand(cmd)
	if err != nil {
		handleError(err)
		return model, nil
	}

	utils.APKGetInfoFromString(&model, infoString)

	_, err = utils.APKMoveFilesToPathForModel(&model, fdstPath)
	if err != nil {
		handleError(err)
		return model, nil
	}
	fmt.Println(model)

	return model, nil
}

func handleIPA(fsrc multipart.File) (models.IPA, error) {
	var model models.IPA
	tempFileName := utils.RandSeq(5) + ".ipa"
	folderPath := path.Join(revel.BasePath, "public", "ipa")
	if !utils.IsFileExistForPath(folderPath) {
		os.MkdirAll(folderPath, 0777)
	}
	fdstPath := path.Join(folderPath, tempFileName)
	fdst, err := os.OpenFile(fdstPath, os.O_CREATE|os.O_WRONLY, 0777)
	defer fdst.Close()
	defer fsrc.Close()
	if err != nil {
		handleError(err)
		return model, err
	}

	// Write file field from file to upload
	_, err = io.Copy(fdst, fsrc)
	if err != nil {
		handleError(err)
		return model, err
	}

	// 从ipa中解析出Info.plist
	plistPath, iconPath, err := utils.IPAUnzipInfoForIPA(fdstPath, path.Join(revel.BasePath, "public", "ipa"))
	if err != nil {
		handleError(err)
		return model, err
	}
	// 解析plist中的信息
	err = utils.IPAGetInfoFromPlist(&model, plistPath)
	if err != nil {
		handleError(err)
		return model, nil
	}

	// 移动文件到各自的目录
	_, err = utils.IPAMoveFilesToPathForModel(&model, fdstPath, plistPath, iconPath)
	if err != nil {
		handleError(err)
		return model, err
	}

	return model, nil
}

// Upload ...
func (c Upload) Upload() revel.Result {
	// Create file field
	fsrc, fsrcHeader, err := c.Request.FormFile("file")
	defer fsrc.Close()
	if err != nil {
		handleError(err)
		return c.RenderError(err)
	}

	if strings.HasSuffix(fsrcHeader.Filename, "ipa") {
		model, err := handleIPA(fsrc)
		if err != nil {
			handleError(err)
			return c.RenderError(err)
		}

		Gdb.Create(&model)
		logger.Println("IPA Web上传")
		logger.Println(model)

		c.RenderArgs["model"] = model

		return c.RenderTemplate("Upload/UploadIPA.html")
	} else if strings.HasSuffix(fsrcHeader.Filename, "apk") {
		model, err := handleAPK(fsrc)
		if err != nil {
			handleError(err)
			return c.RenderError(err)
		}

		Gdb.Create(&model)
		logger.Println("APK Web上传")
		logger.Println(model)

		c.RenderArgs["model"] = model
		return c.RenderTemplate("Upload/UploadAPK.html")
	} else {
		return c.RenderError(errors.New("不支持的附件类型"))
	}
}

// APIUpload ...
func (c Upload) APIUpload() revel.Result {
	// Create file field
	fsrc, fsrcHeader, err := c.Request.FormFile("file")
	defer fsrc.Close()
	if err != nil {
		handleError(err)
		return c.RenderError(err)
	}

	if strings.HasSuffix(fsrcHeader.Filename, "ipa") {
		model, err := handleIPA(fsrc)
		if err != nil {
			handleError(err)
			return c.RenderError(err)
		}

		Gdb.Create(&model)
		logger.Println("IPA API上传")
		logger.Println(model)

		return c.RenderJson(model)
	} else if strings.HasSuffix(fsrcHeader.Filename, "apk") {
		model, err := handleAPK(fsrc)
		if err != nil {
			handleError(err)
			return c.RenderError(err)
		}

		Gdb.Create(&model)
		logger.Println("APK API上传")
		logger.Println(model)

		return c.RenderJson(model)
	} else {
		return c.RenderError(errors.New("不支持的附件类型"))
	}
}
