package types

type CreateDeviceRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type UpdateDeviceRequest struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=1"`
}

type CreateTestSessionRequest struct {
	Name     string `json:"name" validate:"required,min=1"`
	DeviceID uint   `json:"device_id" validate:"required"`
}

type UpdateTestSessionRequest struct {
	ID       uint   `json:"id" validate:"required"`
	Name     string `json:"name" validate:"required,min=1"`
	DeviceID uint   `json:"device_id" validate:"required"`
}

type CreateConditionRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type UpdateConditionRequest struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=1"`
}

type CreateConditionValueRequest struct {
	Value       string `json:"value" validate:"required"`
	ConditionID uint   `json:"condition_id" validate:"required"`
}

type UpdateConditionValueRequest struct {
	ID          uint   `json:"id" validate:"required"`
	Value       string `json:"value" validate:"required"`
	ConditionID uint   `json:"condition_id" validate:"required"`
}

type CreateScenarioRequest struct {
	Name          string `json:"name" validate:"required,min=1"`
	TestSessionID uint   `json:"test_session_id" validate:"required"`
}

type CreateScenarioWithConditionValuesRequest struct {
	Name              string `json:"name" validate:"required,min=1"`
	TestSessionID     uint   `json:"test_session_id" validate:"required"`
	ConditionValueIDs []uint `json:"condition_value_ids" validate:"required"`
}

type UpdateScenarioRequest struct {
	ID            uint   `json:"id" validate:"required"`
	Name          string `json:"name" validate:"required,min=1"`
	TestSessionID uint   `json:"test_session_id" validate:"required"`
}

type CreateScenarioConditionRequest struct {
	ScenarioID       uint `json:"scenario_id" validate:"required"`
	ConditionValueID uint `json:"condition_value_id" validate:"required"`
}

type UpdateScenarioConditionRequest struct {
	ID               uint `json:"id" validate:"required"`
	ScenarioID       uint `json:"scenario_id" validate:"required"`
	ConditionValueID uint `json:"condition_value_id" validate:"required"`
}

type ValidationRequest struct {
	ScenarioID int `json:"scenario_id" validate:"required,gt=0"`
}
