package utils

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"CedAppDeliver/app/models"

	plist "github.com/DHowett/go-plist"
	"github.com/revel/revel"

	"github.com/disintegration/imaging"
)

// IPAUnzipInfoForIPA 从IPA中解压Info.plist, 返回plist和icon路径
func IPAUnzipInfoForIPA(ipaPath string, toPath string) (string, string, error) {
	_, tempName := path.Split(ipaPath)
	reader, err := zip.OpenReader(ipaPath)
	defer reader.Close()
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	var plistFile *zip.File
	var iconFile *zip.File
	var size int64

	for _, k := range reader.Reader.File {
		if k.FileInfo().IsDir() == false {
			dir, fileName := filepath.Split(k.Name)
			list := strings.Split(dir, string(filepath.Separator))

			// 上一级目录的文件夹名字最后是.app
			if strings.HasSuffix(list[len(list)-2], ".app") {
				if fileName == "Info.plist" {
					plistFile = k
				}
				if strings.HasPrefix(fileName, "AppIcon") {
					s := k.FileInfo().Size()
					if s > size {
						size = s
						iconFile = k
					}
				}
			}
		}
	}

	if plistFile != nil && iconFile != nil {
		// 先处理plist文件
		reader1, _ := plistFile.Open()
		newPlistFile, err := os.Create(filepath.Join(toPath, tempName+".plist"))
		defer reader1.Close()
		defer newPlistFile.Close()
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}
		io.Copy(newPlistFile, reader1)

		// 再处理icon文件
		reader2, _ := iconFile.Open()
		newIconFile, err := os.Create(filepath.Join(toPath, tempName+".png"))
		defer reader2.Close()
		defer newIconFile.Close()
		if err != nil {
			fmt.Println(err)
			return "", "", err
		}
		io.Copy(newIconFile, reader2)

		return filepath.Join(toPath, tempName+".plist"), filepath.Join(toPath, tempName+".png"), nil
	}
	return "", "", errors.New("Can not find Info.plist in app folder")
}

// IPAGetInfoFromPlist 从Info.plist中取得应用信息
func IPAGetInfoFromPlist(ipa *models.IPA, plistPath string) error {
	buf, err := ioutil.ReadFile(plistPath)
	if err != nil {
		return err
	}

	var xval = make(map[string]interface{})
	buf1 := bytes.NewReader(buf)

	decoder := plist.NewDecoder(buf1)
	err = decoder.Decode(&xval)
	if err != nil {
		return err
	}

	identifier := xval["CFBundleIdentifier"]
	if i, ok := identifier.(string); ok {
		ipa.BundleIdentifier = i
		ipa.AppID = getMD5([]byte(i))[:4]
	}
	displayName := xval["CFBundleDisplayName"]
	if i, ok := displayName.(string); ok {
		ipa.DisplayName = i
	}
	version := xval["CFBundleVersion"]
	if i, ok := version.(string); ok {
		ipa.Version = i
	}
	shortVersion := xval["CFBundleShortVersionString"]
	if i, ok := shortVersion.(string); ok {
		ipa.ShortVersion = i
	}

	rand := RandSeq(4)
	ipa.RandomID = rand

	return nil
}

// RenameIconFile 修改分辨率并加Icon
func RenameIconFile(iconPath string) error {
	dirPath, _ := path.Split(iconPath)
	src, err := imaging.Open(iconPath)
	if err != nil {
		fmt.Println(err.Error())
		displayImagePath := path.Join(dirPath, "displaySize.png")
		fullImagePath := path.Join(dirPath, "fullSize.png")

		imageFile, _ := os.Open(iconPath)
		newDisplayImageFile, _ := os.Create(displayImagePath)
		newFullImageFile, _ := os.Create(fullImagePath)
		defer newDisplayImageFile.Close()
		defer newFullImageFile.Close()
		defer imageFile.Close()

		io.Copy(newDisplayImageFile, imageFile)
		io.Copy(newFullImageFile, imageFile)

		return nil
	} else {
		img := imaging.Resize(src, 57, 57, imaging.Lanczos)
		displayImagePath := path.Join(dirPath, "displaySize.png")
		imaging.Save(img, displayImagePath)
		img2 := imaging.Resize(src, 512, 512, imaging.Lanczos)
		fullImagePath := path.Join(dirPath, "fullSize.png")
		imaging.Save(img2, fullImagePath)
		return nil
	}
}

// IPAGenInstallPlistFromIPAModel 生成安装所需plist
func IPAGenInstallPlistFromIPAModel(model models.IPA) error {
	var baseURL = "https://" + GetLocalIP()

	metadataDict := map[string]string{
		"bundle-identifier": model.BundleIdentifier,
		"bundle-version":    model.ShortVersion,
		"kind":              "software",
		"title":             model.DisplayName}

	temp := map[string]interface{}{
		"kind": "software-package",
		"url":  baseURL + "/" + model.RelativeIPAURL}
	temp2 := map[string]interface{}{
		"kind": "display-image",
		"url":  baseURL + "/" + model.RelativeIPADisplayImage}
	temp3 := map[string]interface{}{
		"kind": "full-size-image",
		"url":  baseURL + "/" + model.RelativeIPAFullSizeImage}

	temp4 := map[string]interface{}{
		"assets":   [3]interface{}{temp, temp2, temp3},
		"metadata": metadataDict}
	pl := map[string]interface{}{
		"items": [1]interface{}{temp4}}

	writer, err := os.Create(path.Join(revel.BasePath, model.RelativePlistPath))
	defer writer.Close()
	if err != nil {
		return err
	}

	encoder := plist.NewEncoder(writer)
	encoder.Encode(pl)
	return nil
}

// IPAMoveFilesToPathForModel 移动文件到对应文件夹
func IPAMoveFilesToPathForModel(model *models.IPA, ipaPath string, plistPath string, iconPath string) ([]string, error) {
	var newFilePaths []string

	timeStr := time.Now().Format("2006-01-02_15-04-05")
	folderPath := path.Join("public", "ipa", model.BundleIdentifier, timeStr)
	model.RelativeIPAPath = path.Join(folderPath, model.BundleIdentifier+".ipa")
	model.RelativePlistPath = path.Join(folderPath, model.BundleIdentifier+".plist")

	model.RelativeIPADisplayImage = "public/ipa/" + model.BundleIdentifier + "/" + "displaySize.png"
	model.RelativeIPAFullSizeImage = "public/ipa/" + model.BundleIdentifier + "/" + "fullSize.png"
	url := "public/ipa/" + model.BundleIdentifier + "/" + timeStr + "/"
	model.RelativeIPAURL = url + model.BundleIdentifier + ".ipa"
	model.RelativePlistURL = url + model.BundleIdentifier + ".plist"

	folderPath = path.Join(revel.BasePath, folderPath)

	if !IsFileExistForPath(folderPath) {
		err := os.MkdirAll(folderPath, 0777)
		if err != nil {
			return newFilePaths, err
		}
	}

	err := IPAGenInstallPlistFromIPAModel(*model)
	if err != nil {
		fmt.Println(err.Error())
		return newFilePaths, nil
	}

	// 移动Icon
	tempDir, tempFileName := path.Split(iconPath)
	newIconPath := path.Join(tempDir, model.BundleIdentifier, tempFileName)
	err = os.Rename(iconPath, newIconPath)
	if err != nil {
		return newFilePaths, err
	}

	// Icon分辨率修改
	err = RenameIconFile(newIconPath)
	if err != nil {
		return newFilePaths, err
	}

	// 删除无用图片
	err = os.Remove(newIconPath)
	if err != nil {
		return newFilePaths, err
	}

	// 移动IPA
	newFilePath := path.Join(folderPath, model.BundleIdentifier+".ipa")
	err = os.Rename(ipaPath, newFilePath)
	if err != nil {
		return newFilePaths, err
	}
	newFilePaths = append(newFilePaths, newFilePath)

	// 删除无用plist
	err = os.Remove(plistPath)
	if err != nil {
		return newFilePaths, err
	}

	return newFilePaths, nil
}
