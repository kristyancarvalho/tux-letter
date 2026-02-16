package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Article struct {
	ID    uint   `gorm:"primaryKey"`
	URL   string `gorm:"uniqueIndex;not null"`
	Title string
}

var DB *gorm.DB

func Init(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&Article{})
}

func IsURLSeen(url string) bool {
	var count int64
	DB.Model(&Article{}).Where("url = ?", url).Count(&count)
	return count > 0
}

func SaveArticle(url, title string) error {
	article := Article{URL: url, Title: title}
	return DB.Create(&article).Error
}
