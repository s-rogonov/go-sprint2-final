package dbprovider

import (
	"consts"
	"dbprovider/models"
)

func (m *manager) GetTimings(timings *models.Timings) error {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	return m.db.First(timings).Error
}

func (m *manager) GetQueries(limit uint) (qs []*models.Query, err error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	err = m.db.Preload(consts.ModelQueryTasksField).Limit(int(limit)).Order("id desc").Find(&qs).Error
	if err != nil {
		return nil, err
	}

	return
}

func (m *manager) GetQuery(id uint) (q *models.Query, err error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	q = &models.Query{}
	err = m.db.Preload(consts.ModelQueryTasksField).First(q, id).Error
	if err != nil {
		return nil, err
	}

	return
}

func (m *manager) GetWorkers(limit uint) (ws []*models.Worker, err error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	err = m.db.Preload(
		consts.ModelWorkerTargetField,
	).Preload(
		consts.ModelWorkerTargetField + "." + consts.ModelTaskSubtasksField,
	).Limit(int(limit)).Order("`workers`.`id` desc").Find(&ws).Error

	if err != nil {
		return nil, err
	}

	return
}

func (m *manager) GetWorker(id uint) (w *models.Worker, err error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	w = &models.Worker{}
	err = m.db.Joins(
		consts.ModelWorkerTargetField,
	).Preload(
		consts.ModelWorkerTargetField+"."+consts.ModelTaskSubtasksField,
	).First(w, id).Error

	if err != nil {
		return nil, err
	}

	return
}
