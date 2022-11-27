package models

import "gorm.io/gorm"

type Book struct {
	ID        uint    `json:"id" gorm:"primary_key;auto_increment"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

func MigrateBooks(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}
