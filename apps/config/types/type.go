package types

type CreateDeviceRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

type UpdateDeviceRequest struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required,min=1"`
}
