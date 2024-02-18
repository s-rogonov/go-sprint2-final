package helpers

import (
	"fmt"

	"consts"
	"dbprovider/models"
)

var (
	ErrQueryContractBothPresented = fmt.Errorf(
		"both `%s` and `%s` are presented in a query",
		consts.ModelQueryTasksField, consts.ModelQueryBadMessageField)

	ErrQueryContractBothMissed = fmt.Errorf(
		"neither `%s` nor `%s` are presented in a query",
		consts.ModelQueryTasksField, consts.ModelQueryBadMessageField)
)

// CheckQueryContract
// checks if query meets one of two cases:
//  1. query has a non-empty consts.ModelQueryTasksField and an empty consts.ModelQueryBadMessageField
//  2. query has an empty consts.ModelQueryTasksField and a non-empty consts.ModelQueryBadMessageField
func CheckQueryContract(query *models.Query) error {
	hasTasks := len(query.Tasks) > 0
	hasMessage := query.BadMessage != ""

	if hasTasks && hasMessage {
		return ErrQueryContractBothPresented
	}

	if (!hasTasks) && (!hasMessage) {
		return ErrQueryContractBothMissed
	}

	return nil
}
