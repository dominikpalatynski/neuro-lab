package database

import (
	"gorm.io/gorm"
)

type Status string

const (
	StatusCompleted Status = "COMPLETE"
	StatusActive    Status = "ACTIVE"
	StatusInactive  Status = "INACTIVE"
)

type Device struct {
	gorm.Model
	Name string `json:"name" validate:"required,min=1"`
}

type TestSession struct {
	gorm.Model
	Name string `json:"name" validate:"required,min=1"`

	DeviceID uint    `json:"device_id" validate:"required"`
	Device   *Device `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"Device,omitempty"`
}

type Condition struct {
	gorm.Model
	Name string `json:"name" validate:"required,min=1"`
}

type ConditionValue struct {
	gorm.Model
	Value       string     `json:"value" validate:"required"`
	ConditionID uint       `json:"condition_id" validate:"required"`
	Condition   *Condition `gorm:"foreignKey:ConditionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"Condition,omitempty"`
}

type Scenario struct {
	gorm.Model
	Name               string              `json:"name" validate:"required,min=1"`
	Status             Status              `json:"status" gorm:"default:INACTIVE"`
	TestSessionID      uint                `json:"test_session_id" validate:"required"`
	TestSession        *TestSession        `gorm:"foreignKey:TestSessionID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"TestSession,omitempty"`
	ScenarioConditions []ScenarioCondition `gorm:"foreignKey:ScenarioID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"ScenarioConditions,omitempty"`
}

type ScenarioCondition struct {
	gorm.Model
	ScenarioID       uint            `json:"scenario_id" validate:"required"`
	Scenario         *Scenario       `gorm:"foreignKey:ScenarioID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"Scenario,omitempty"`
	ConditionValueID uint            `json:"condition_value_id" validate:"required"`
	ConditionValue   *ConditionValue `gorm:"foreignKey:ConditionValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"ConditionValue,omitempty"`
}
