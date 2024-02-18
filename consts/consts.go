package consts

const (
	DbEnvironmentKey = "DBNAME"
	DbProductionName = "sqlite.db"
	DbTestName       = "test_sqlite.db"

	ModelTaskParentField     = "ParentTask"
	ModelTaskTargetField     = "TargetQuery"
	ModelWorkerTargetField   = "TargetTask"
	ModelTaskLastWorkerField = "LastWorker"

	ModelQueryTasksField      = "Tasks"
	ModelQueryBadMessageField = "BadMessage"
	ModelTaskSubtasksField    = "Subtasks"
	ModelTaskWorkersField     = "Workers"
	ModelTaskIsDoneField      = "IsDone"
	ModelTaskIsReadyField     = "IsReady"
)
