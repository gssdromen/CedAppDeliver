package utils

import (
	"CedAppDeliver/app/models"
	"path"
	"time"

	"fmt"

	"os"

	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

// IsFileExistForPath 判断文件是否存在
func IsFileExistForPath(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func findIPAsBefore(t time.Time, db *gorm.DB) ([]models.IPA, error) {
	var ipas []models.IPA

	db.Find(&ipas)

	var results []models.IPA
	for _, v := range ipas {
		if v.CreatedAt.Before(t) {
			results = append(results, v)
		}
	}

	return results, nil
}

func findAPKsBefore(t time.Time, db *gorm.DB) ([]models.APK, error) {
	var apks []models.APK

	db.Find(&apks)

	var results []models.APK
	for _, v := range apks {
		if v.CreatedAt.Before(t) {
			results = append(results, v)
		}
	}

	return results, nil
}

// CleanOutdateApps ...
func CleanOutdateApps() {
	db, err := gorm.Open("sqlite3", path.Join(revel.BasePath, "database.db"))
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	sevenDayBefore := time.Now().AddDate(0, 0, -7)

	// 清理IPA
	ipas, err := findIPAsBefore(sevenDayBefore, db)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range ipas {
		// 不删除关键历史版本IPA
		if v.IsMilestoneVersion {
			continue
		}
		dirPath, _ := path.Split(v.RelativePlistPath)
		var folderPath = path.Join(revel.BasePath, dirPath)
		_ = os.RemoveAll(folderPath)
		db.Unscoped().Delete(&v)
	}

	// 清理APK
	apks, err := findAPKsBefore(sevenDayBefore, db)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range apks {
		// 不删除关键历史版本APK
		if v.IsMilestoneVersion {
			continue
		}
		dirPath, _ := path.Split(v.RelativePath)
		var folderPath = path.Join(revel.BasePath, dirPath)
		_ = os.RemoveAll(folderPath)
		db.Unscoped().Delete(&v)
	}

}
