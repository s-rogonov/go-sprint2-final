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
	ClearDB() error

	UpdateTimings(timings *models.Timings) error

	NewQuery(query *models.Query) error
	UpdateQuery(query *models.Query) error

	CreateWorkers(amount uint) ([]*models.Worker, error)
	SetWorkResult(workerID uint, result float64) error

	GetTimings(timings *models.Timings) error
	GetQueries(limit uint) (qs []*models.Query, err error)
	GetQuery(id uint) (q *models.Query, err error)
}

type manager struct {
	db      *gorm.DB
	rwMutex *sync.RWMutex
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
		db:      db,
		rwMutex: &sync.RWMutex{},
	}

	if err := Manager.InitDB(); err != nil {
		panic(err)
	}

	count := int64(0)
	if err := db.Model(models.Timings{}).Count(&count).Error; err != nil {
		panic(err)
	}

	if count == 0 {
		err := Manager.ClearDB()
		if err != nil {
			panic(err)
		}
	}
}
