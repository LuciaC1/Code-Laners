package dto

type UpdateUserDTO struct {
	Name   *string  `json:"name,omitempty"`
	Email  *string  `json:"email,omitempty"` // validar en servicio
	Weight *float64 `json:"weight,omitempty"`
	Height *float64 `json:"height,omitempty"`
	Level  *string  `json:"level,omitempty"`
	Goals  []string `json:"goals,omitempty"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
