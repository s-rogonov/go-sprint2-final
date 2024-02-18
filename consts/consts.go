package consts

const (
	DbEnvironmentKey = "DBNAME"
	DbProductionName = "sqlite.db"
	DbTestName       = "test_sqlite.db"

	ModelTaskParentField     = "Parent"
	ModelTaskTargetField     = "Target"
	ModelWorkerTargetField   = "Target"
	ModelTaskLastWorkerField = "LastWorker"

	ModelQueryTasksField      = "Tasks"
	ModelQueryBadMessageField = "BadMessage"
	ModelTaskSubtasksField    = "Subtasks"
	ModelTaskWorkersField     = "Workers"
	ModelTaskIsDoneField      = "IsDone"
	ModelTaskIsReadyField     = "IsReady"
)
