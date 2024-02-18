package helpers

import (
	"consts"
	"dbprovider/models"
	"gorm.io/gorm"
)

func FinishWorker(db *gorm.DB, worker *models.Worker, result float64) error {
	if db != nil && worker.TargetTask == nil {
		err := db.Preload(consts.ModelWorkerTargetField).First(worker).Error
		if err != nil {
			return err
		}
	}

	worker.IsDone = true
	worker.Result = result

	if db != nil {
		err := db.Save(worker).Error
		if err != nil {
			return err
		}
	}

	return FinishTask(db, worker.TargetTask, worker.Result)
}

func FinishTask(db *gorm.DB, target *models.Task, result float64) error {
	if db != nil && (target.ParentTask == nil || target.TargetQuery == nil) {
		err := db.Preload(consts.ModelTaskTargetField).Preload(consts.ModelTaskParentField).First(target).Error
		if err != nil {
			return err
		}
	}

	target.IsDone = true
	target.Result = result

	if db != nil {
		err := db.Save(target).Error
		if err != nil {
			return err
		}
	}

	if target.ParentTask == nil {
		return FinishQuery(db, target.TargetQuery, result)
	} else {
		return IncrementFinishedSubtasks(db, target.ParentTask)
	}
}

func IncrementFinishedSubtasks(db *gorm.DB, parent *models.Task) error {
	parent.FinishedSubtasks += 1
	parent.IsReady = parent.TotalSubtasks == parent.FinishedSubtasks
	if db == nil {
		return nil
	} else {
		return db.Save(parent).Error
	}
}

func FinishQuery(db *gorm.DB, target *models.Query, result float64) error {
	target.Result = result
	target.IsDone = true
	if db == nil {
		return nil
	} else {
		return db.Save(target).Error
	}
}
