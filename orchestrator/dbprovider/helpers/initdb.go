package helpers

import (
	"time"

	"dbprovider/models"
	"gorm.io/gorm"
)

func CreateDefaults(db *gorm.DB) error {
	result := db.Create(&models.Timings{
		Factor:         2.0,
		Addition:       1 * time.Second,
		Multiplication: 1 * time.Second,
		Subtraction:    1 * time.Second,
		Division:       1 * time.Second,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DropTables(db *gorm.DB) (err error) {
	err = db.Unscoped().Where("1 = 1").Delete(&models.Timings{}).Error
	if err != nil {
		return
	}
	err = db.Unscoped().Where("1 = 1").Delete(&models.Worker{}).Error
	if err != nil {
		return
	}
	err = db.Unscoped().Where("1 = 1").Delete(&models.Task{}).Error
	if err != nil {
		return
	}
	err = db.Unscoped().Where("1 = 1").Delete(&models.Query{}).Error
	if err != nil {
		return
	}
	return nil
}

func MigrateSchemes(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.Timings{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Query{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Task{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Worker{}); err != nil {
		return err
	}
	return nil
}
