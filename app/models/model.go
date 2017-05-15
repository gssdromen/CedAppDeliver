package models

import (
	"bytes"

	"github.com/jinzhu/gorm"
)

// AppAlias ...
type AppAlias struct {
	AppID    string `gorm:"size:255"`
	AppAlias string `gorm:"size:255"`
}

func (appAlias AppAlias) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("应用ID：")
	buffer.WriteString(appAlias.AppID)
	buffer.WriteString(" | 应用别名：")
	buffer.WriteString(appAlias.AppAlias)
	return buffer.String()
}

// IPA Model of IPA
type IPA struct {
	gorm.Model
	RelativePlistPath        string `gorm:"size:255;unique;not null"`
	RelativePlistURL         string `gorm:"size:255;unique;not null"`
	RelativeIPAPath          string `gorm:"size:255;unique;not null"`
	RelativeIPAURL           string `gorm:"size:255;unique;not null"`
	RelativeIPADisplayImage  string `gorm:"size:255"`
	RelativeIPAFullSizeImage string `gorm:"size:255"`
	BundleIdentifier         string `gorm:"size:255;not null"`
	DisplayName              string `gorm:"size:255;not null"`
	Version                  string `gorm:"size:255;not null"`
	ShortVersion             string `gorm:"size:255"`
	AppID                    string `gorm:"size:255;not null"`
	RandomID                 string `gorm:"size:255;unique;not null"`
	IsMilestoneVersion       bool
	Comment                  string `gorm:"size:255"`
}

func (ipa IPA) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("应用名：")
	buffer.WriteString(ipa.DisplayName)
	buffer.WriteString(" | BundleID：")
	buffer.WriteString(ipa.BundleIdentifier)
	buffer.WriteString(" | Version：")
	buffer.WriteString(ipa.Version)
	buffer.WriteString(" | ShortVersion：")
	buffer.WriteString(ipa.ShortVersion)
	buffer.WriteString(" | 创建时间：")
	buffer.WriteString(ipa.CreatedAt.Format("2006-01-02 15:04:05"))
	buffer.WriteString(" | 修改时间：")
	buffer.WriteString(ipa.UpdatedAt.Format("2006-01-02 15:04:05"))

	return buffer.String()
}

// APK Model of APK
type APK struct {
	gorm.Model
	DisplayName        string `gorm:"size:255;not null"`
	PackageName        string `gorm:"size:255;not null"`
	Version            string `gorm:"size:255;not null"`
	ShortVersion       string `gorm:"size:255"`
	RelativePath       string `gorm:"size:255;unique;not null"`
	RelativeURL        string `gorm:"size:255;unique;not null"`
	AppID              string `gorm:"size:255;not null"`
	RandomID           string `gorm:"size:255;unique;not null"`
	IsMilestoneVersion bool
	Comment            string `gorm:"size:255"`
}

func (apk APK) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("应用名：")
	buffer.WriteString(apk.DisplayName)
	buffer.WriteString(" | PackageName：")
	buffer.WriteString(apk.PackageName)
	buffer.WriteString(" | Version：")
	buffer.WriteString(apk.Version)
	buffer.WriteString(" | 创建时间：")
	buffer.WriteString(apk.CreatedAt.Format("2006-01-02 15:04:05"))
	buffer.WriteString(" | 修改时间：")
	buffer.WriteString(apk.UpdatedAt.Format("2006-01-02 15:04:05"))

	return buffer.String()
}

// ResponseModel Json的返回模型
type ResponseModel struct {
	data interface{}
	code int
	msg  string
}
