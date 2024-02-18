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
