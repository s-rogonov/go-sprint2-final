package consts

const (
	DbEnvironmentKey = "DBNAME"
	DbProductionName = "sqlite.db"
	DbTestName       = "test_sqlite.db"

	EnvPort     = "PORT"
	EnvMaster   = "MASTER"
	EnvNWorkers = "NWORKERS"
	EnvBatch    = "BATCH"
	EnvDelay    = "DELAY"

	OrchestratorDefaultPort = "8181"

	AgentDefaultMaster  = "localhost:8181"
	AgentDefaultWorkers = 5
	AgentDefaultBatch   = 2
	AgentDefaultDelay   = 10

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
