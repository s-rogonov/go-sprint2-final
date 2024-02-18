package dbprovider

import (
	"context"
	"os"
	"sync"

	"consts"
	"dbprovider/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Actions interface {
	getDB() *gorm.DB // for test purposes only
	WithContext(ctx context.Context) Actions
	InitDB() error

	UpdateTimings(timings *models.Timings) error

	NewQuery(query *models.Query) error
	UpdateQuery(query *models.Query) error

	CreateWorkers(amount uint) ([]*models.Worker, error)
	SetWorkResult(workerID uint, result float64) error
}

type manager struct {
	db          *gorm.DB
	lock4update sync.Mutex
}

func (m *manager) getDB() *gorm.DB {
	return m.db
}

var Manager Actions

func InitConnection() {
	dbname, ok := os.LookupEnv(consts.DbEnvironmentKey)
	if !ok {
		dbname = consts.DbProductionName
	}

	db, err := gorm.Open(sqlite.Open(dbname), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	Manager = &manager{
		db: db,
	}
}
