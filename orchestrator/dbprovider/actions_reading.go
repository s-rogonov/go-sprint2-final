package dbprovider

import "dbprovider/models"

func (m *manager) GetTimings(timings *models.Timings) error {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	return m.db.First(timings).Error
}
