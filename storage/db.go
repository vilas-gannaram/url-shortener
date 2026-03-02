package storage

import (
	"github.com/glebarez/sqlite" // <--- Change this import
	"gorm.io/gorm"
)

type URLMapping struct {
	gorm.Model         // Adds ID, CreatedAt, UpdatedAt, DeletedAt automatically
	OriginalURL string `gorm:"not null"`
	ShortKey    string `gorm:"uniqueIndex;not null"`
}

func InitDB() *gorm.DB {
	// It works exactly the same way, just no C compiler needed
	db, err := gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error()) // Added err.Error() to see WHY
	}
	db.AutoMigrate(&URLMapping{})
	return db
}
