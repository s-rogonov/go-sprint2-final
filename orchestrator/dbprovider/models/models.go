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

	Tasks []*Task `gorm:"foreignKey:TargetID"`
}

type Task struct {
	gorm.Model

	Operation string `gorm:"size:1"`
	Duration  time.Duration
	Index     uint

	Parent       *Task   `gorm:"foreignkey:ParentID;association_foreignkey:ID"`
	ParentID     uint    `gorm:"index"`
	Target       *Query  `gorm:"foreignkey:TargetID;association_foreignkey:ID"`
	TargetID     uint    `gorm:"index"`
	LastWorker   *Worker `gorm:"foreignkey:LastWorkerID;association_foreignkey:ID"`
	LastWorkerID uint    `gorm:"index"`

	TotalSubtasks    uint
	FinishedSubtasks uint
	IsDone           bool
	IsReady          bool `gorm:"index"`
	Result           float64

	Workers  []*Worker `gorm:"foreignKey:TargetID"`
	Subtasks []*Task   `gorm:"foreignKey:ParentID"`
}

type Worker struct {
	gorm.Model

	Target   *Task `gorm:"foreignkey:TargetID;association_foreignkey:ID"`
	TargetID uint  `gorm:"index"`
	IsDone   bool  `gorm:"index"`
	Result   float64
}
