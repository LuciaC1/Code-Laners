package utils

import (
	"backend/dto"
	"backend/models"
)

func ConvertUserModelToRegisterRequest(user models.User) dto.RegisterRequest {
	return dto.RegisterRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.PasswordHash,
	}
}

func ConverUserModelToLoginRequest(user models.User) dto.LoginRequest {
	return dto.LoginRequest{
		Email:    user.Email,
		Password: user.PasswordHash,
	}
}

func ConvertUserModelToUpdateUserRequest(user models.User) dto.UpdateUserRequest {
	return dto.UpdateUserRequest{
		Name:  user.Name,
		Email: user.Email,
	}
}

func ConvertUserModelToChangePasswordRequest(user models.User, newPassword string) dto.ChangePasswordRequest {
	return dto.ChangePasswordRequest{
		OldPassword: user.PasswordHash,
		NewPassword: newPassword,
	}
}


