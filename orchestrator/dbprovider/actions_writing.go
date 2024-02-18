package dbprovider

import (
	"errors"
	"fmt"

	"consts"
	"dbprovider/helpers"
	"dbprovider/models"
	"gorm.io/gorm"
)

var (
	ErrUpdateQueryWithoutID     = errors.New("query entity has no ID, so it cannot be updated")
	ErrGoodQueryCannotBeUpdated = fmt.Errorf(
		"query entity already has an empty `%s` in DB, thus it cannot be updated",
		consts.ModelQueryBadMessageField)

	ErrWorkerAlreadyCompleted = errors.New("worker entity already completed and has a valid result")
)

func (m *manager) InitDB() error {
	if err := helpers.MigrateSchemes(m.db); err != nil {
		return err
	}

	if err := helpers.DropTables(m.db); err != nil {
		return err
	}

	if err := helpers.CreateDefaults(m.db); err != nil {
		return err
	}

	return nil
}

func (m *manager) UpdateTimings(timings *models.Timings) error {
	return m.db.Save(timings).Error
}

// initQueryTasks
// inits tasks and their associations
func initQueryTasks(tx *gorm.DB, query *models.Query) error {
	if len(query.Tasks) == 0 {
		return nil
	}

	if err := tx.Omit(consts.ModelTaskSubtasksField).Create(query.Tasks).Error; err != nil {
		return err
	}

	for _, t := range query.Tasks {
		helpers.DefaultsForTask(t)
		t.Target = query
		for i, st := range t.Subtasks {
			st.Index = uint(i)
			st.Parent = t
		}
	}

	helpers.DefaultsForQuery(query)

	for _, t := range query.Tasks {
		if t.IsDone {
			err := helpers.FinishTask(nil, t, t.Result)
			if err != nil {
				return err
			}
		}

		{
			err := tx.Omit(consts.ModelTaskTargetField).Omit(consts.ModelTaskParentField).Save(t).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// NewQuery
// proceeds a creation only if query passes helpers.CheckQueryContract
//
// Also creates nested tasks
func (m *manager) NewQuery(query *models.Query) error {
	if err := helpers.CheckQueryContract(query); err != nil {
		return err
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := initQueryTasks(tx, query); err != nil {
			return err
		}

		return tx.Save(query).Error
	})
}

// UpdateQuery
// proceeds an update only if:
//   - an ID presented
//   - query passes helpers.CheckQueryContract
//   - query with such ID has a non-empty consts.ModelQueryBadMessageField in DB
func (m *manager) UpdateQuery(query *models.Query) error {
	if query.ID == 0 {
		return ErrUpdateQueryWithoutID
	}

	if err := helpers.CheckQueryContract(query); err != nil {
		return err
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		existed := &models.Query{}

		if err := tx.First(&existed, query.ID).Error; err != nil {
			return err
		}

		if existed.BadMessage == "" {
			return ErrGoodQueryCannotBeUpdated
		}

		if err := initQueryTasks(tx, query); err != nil {
			return err
		}

		return tx.Save(query).Error
	})
}

func (m *manager) CreateWorkers(amount uint, factor float32) (workers []*models.Worker, err error) {
	err = m.db.Transaction(func(tx *gorm.DB) error {

		return nil
	})
	return
}

func (m *manager) SetWorkResult(workerID uint, result float64) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		worker := &models.Worker{}

		if err := m.db.First(worker, workerID).Error; err != nil {
			return err
		}

		if worker.IsDone {
			return ErrWorkerAlreadyCompleted
		}

		return helpers.FinishWorker(tx, worker, result)
	})
}
