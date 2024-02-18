package dbprovider

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"consts"
	"dbprovider/helpers"
	"dbprovider/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrUpdateQueryWithoutID     = errors.New("query entity has no ID, so it cannot be updated")
	ErrGoodQueryCannotBeUpdated = fmt.Errorf(
		"query entity already has an empty `%s` in DB, thus it cannot be updated",
		consts.ModelQueryBadMessageField)

	ErrWorkerAlreadyCompleted = errors.New("worker entity already completed and has a valid result")
)

func (m *manager) WithContext(ctx context.Context) Actions {
	return &manager{db: m.db.WithContext(ctx)}
}

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

func (m *manager) CreateWorkers(amount uint) (workers []*models.Worker, err error) {
	err = m.db.Transaction(func(tx *gorm.DB) (err error) {
		var rows *sql.Rows
		{ // run rows selector
			rows, err =
				tx.Clauses(
					clause.Locking{
						Strength: clause.LockingStrengthUpdate,    // prevent rows from access from other threads
						Options:  clause.LockingOptionsSkipLocked, // skip locked rows
					},
				).Model(
					models.Task{},
				).Preload(
					consts.ModelTaskLastWorkerField,
				).Where(
					models.Task{
						IsDone:  false,
						IsReady: true,
					},
					consts.ModelTaskIsDoneField,
					consts.ModelTaskIsReadyField,
				).Rows()

			defer func() {
				err2 := rows.Close()
				if err == nil {
					err = err2
				}
			}()

			if err != nil {
				return err
			}
		}

		now := time.Now()
		timings := models.Timings{}
		if err = m.db.First(&timings).Error; err != nil {
			return err
		}

		var readyTasks []*models.Task

		for rows.Next() && uint(len(readyTasks)) < amount {
			task := &models.Task{}

			// ScanRows scans a row into a struct
			if err = tx.ScanRows(rows, task); err != nil {
				return err
			}

			// Perform operations on each task
			if task.LastWorker != nil {
				deadline := task.LastWorker.CreatedAt.Add(time.Duration(timings.Factor * float32(task.Duration)))
				if now.Before(deadline) {
					continue // task is still in progress
				}
			}

			readyTasks = append(readyTasks, task)
		}

		if err = rows.Err(); err != nil {
			return err
		}

		workers = make([]*models.Worker, len(readyTasks))
		for i, task := range readyTasks {
			workers[i] = &models.Worker{
				Target: task,
			}
			task.LastWorker = workers[i]
			task.LastWorkerID = 0

			if err = m.db.Save(task).Error; err != nil {
				return err
			}
		}

		return nil
	})

	// check for error
	if err != nil {
		return nil, err
	}

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
