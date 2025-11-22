package database

import (
	"time"

	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

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

type ProcessedSample struct {
	gorm.Model
	DeviceID   uint      `json:"device_id"`
	ScenarioID uint      `json:"scenario_id"`
	FrameID    uint      `json:"frame_id"`
	MetricName string    `json:"metric_name"`
	Value      float64   `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
}

type Float8Array []float64

func (x *Float8Array) Scan(value any) error {
	if value == nil {
		*x = nil
		return nil
	}
	str := value.(string)
	if str == "{}" {
		*x = []float64{}
		return nil
	}
	str, _ = strings.CutPrefix(str, "{")
	str, _ = strings.CutSuffix(str, "}")
	parts := strings.Split(str, ",")
	for _, s := range parts {
		if s == "" || s == `""` {
			panic("empty string is not a valid float64")
		}
		num, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic("not a valid float64")
		}
		*x = append(*x, num)
	}
	return nil
}

func (x Float8Array) Value() (driver.Value, error) {
	if len(x) == 0 {
		return "{}", nil
	}
	formatted := []string{}
	for _, s := range x {
		formatted = append(formatted, fmt.Sprint(s))
	}
	str := "{" + strings.Join(formatted, ",") + "}"
	return str, nil
}

type ProcessedChannel struct {
	gorm.Model
	Values     Float8Array `json:"values" gorm:"type:double precision[]"`
	Timestamp  time.Time   `json:"timestamp"`
	FrameID    uint        `json:"frame_id"`
	MetricName string      `json:"metric_name"`
	DeviceID   uint        `json:"device_id"`
	ScenarioID uint        `json:"scenario_id"`
}
