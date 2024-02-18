package dbprovider

import "dbprovider/models"

func (m *manager) GetTimings(timings *models.Timings) error {
	return m.db.First(timings).Error
}
