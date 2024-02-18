package models

import (
	"time"

	"gorm.io/gorm"
)

type Timings struct {
	gorm.Model

	Factor         float32 // multiplier, applied to duration, defines worker timeout (deadline)
	Addition       time.Duration
	Multiplication time.Duration
	Subtraction    time.Duration
	Division       time.Duration
}

type Query struct {
	gorm.Model

	Expression   string
	BadMessage   string
	HasError     bool `gorm:"index"`
	PlainNumbers uint
	IsDone       bool `gorm:"index"`
	Result       float64

	Tasks []*Task `gorm:"foreignKey:TargetQueryID"`
}

type Task struct {
	gorm.Model

	Operation string `gorm:"size:1"`
	Duration  time.Duration
	Index     uint

	ParentTask    *Task   `gorm:"foreignKey:ParentTaskID"`
	ParentTaskID  uint    `gorm:"index"`
	TargetQuery   *Query  `gorm:"foreignKey:TargetQueryID"`
	TargetQueryID uint    `gorm:"index"`
	LastWorker    *Worker `gorm:"foreignKey:LastWorkerID"`
	LastWorkerID  uint    `gorm:"index"`

	TotalSubtasks    uint
	FinishedSubtasks uint
	IsDone           bool
	IsReady          bool `gorm:"index"`
	Result           float64

	Workers  []*Worker `gorm:"foreignKey:TargetTaskID"`
	Subtasks []*Task   `gorm:"foreignKey:ParentTaskID"`
}

type Worker struct {
	gorm.Model

	TargetTask   *Task `gorm:"foreignkey:TargetTaskID"`
	TargetTaskID uint  `gorm:"index"`
	IsDone       bool  `gorm:"index"`
	Result       float64
}
