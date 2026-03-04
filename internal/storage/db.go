package storage

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type URLMapping struct {
	ID          uint   `gorm:"column:id;primaryKey"`
	OriginalURL string `gorm:"column:original_url;not null"`
	ShortKey    string `gorm:"column:short_key;uniqueIndex;not null"`
}

type URLStats struct {
	ID              uint       `gorm:"column:id;primaryKey"`
	URLMappingID    uint       `gorm:"column:url_mapping_id;not null"`
	URLMapping      URLMapping `gorm:"column:url_mapping_id;foreignKey:URLMappingID"`
	RedirectedCount int        `gorm:"column:redirected_count;not null"`
	LastUpdated     int64      `gorm:"column:last_updated;not null"`
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	db.AutoMigrate(&URLMapping{}, &URLStats{})
	return db
}
