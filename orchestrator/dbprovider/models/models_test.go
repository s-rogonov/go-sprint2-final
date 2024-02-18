package models

import (
	"reflect"
	"testing"

	"consts"
)

func TestFieldsPresented(t *testing.T) {
	{
		rt := reflect.TypeOf(Query{})
		checkFieldPresence(t, rt, consts.ModelQueryTasksField)
		checkFieldPresence(t, rt, consts.ModelQueryBadMessageField)
	}

	{
		rt := reflect.TypeOf(Task{})

		checkFieldPresence(t, rt, consts.ModelTaskTargetField)
		checkFieldPresence(t, rt, consts.ModelTaskTargetField+"ID")

		checkFieldPresence(t, rt, consts.ModelTaskParentField)
		checkFieldPresence(t, rt, consts.ModelTaskParentField+"ID")

		checkFieldPresence(t, rt, consts.ModelTaskSubtasksField)
		checkFieldPresence(t, rt, consts.ModelTaskWorkersField)
		checkFieldPresence(t, rt, consts.ModelTaskIsDoneField)
	}

	{
		rt := reflect.TypeOf(Worker{})
		checkFieldPresence(t, rt, consts.ModelWorkerTargetField)
		checkFieldPresence(t, rt, consts.ModelWorkerTargetField+"ID")
	}
}

func checkFieldPresence(t *testing.T, rtype reflect.Type, fieldName string) {
	hasField := false
	for i := 0; i < rtype.NumField(); i++ {
		hasField = hasField || (rtype.Field(i).Name == fieldName)
	}
	if !hasField {
		t.Errorf("%v missing field `%s`", rtype, fieldName)
	}
}
