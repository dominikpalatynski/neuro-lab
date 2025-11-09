package database

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	Name string `json:"name" validate:"required,min=1"`
}

type TestSession struct {
	gorm.Model
	Name string `json:"name" validate:"required,min=1"`

	DeviceID uint   `json:"device_id" validate:"required"`
	Device   Device `gorm:"foreignKey:DeviceID"`
}

type Condition struct {
	gorm.Model
	Name string `json:"name" validate:"required,min=1"`
}

type ConditionValue struct {
	gorm.Model
	Value       string    `json:"value" validate:"required"`
	ConditionID uint      `json:"condition_id" validate:"required"`
	Condition   Condition `gorm:"foreignKey:ConditionID"`
}

type Scenario struct {
	gorm.Model
	Name          string      `json:"name" validate:"required,min=1"`
	TestSessionID uint        `json:"test_session_id" validate:"required"`
	TestSession   TestSession `gorm:"foreignKey:TestSessionID"`
}

type ScenarioCondition struct {
	gorm.Model
	ScenarioID       uint           `json:"scenario_id" validate:"required"`
	Scenario         Scenario       `gorm:"foreignKey:ScenarioID"`
	ConditionValueID uint           `json:"condition_value_id" validate:"required"`
	ConditionValue   ConditionValue `gorm:"foreignKey:ConditionValueID"`
}
