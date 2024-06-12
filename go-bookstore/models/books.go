package models

import "gorm.io/gorm"

type Book struct {
	Author *string `json:"author"`
	Title  *string `json:"title"`
	Year   *string `json:"year"`
	ISBN   *string `json:"isbn"`
	ID     uint    `gorm:"primary key; autoIncrement" json:"id"`
}

func MigrateBook(db *gorm.DB) error {
	err := db.AutoMigrate(&Book{})
	return err
}
