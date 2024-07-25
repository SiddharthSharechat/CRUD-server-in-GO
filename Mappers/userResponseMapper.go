package Mappers

import "github.com/SiddharthSharechat/CRUDGo/Models"

func UserResponseMapper(user Models.User) Models.UserResponse {
	return Models.UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Location:  user.Location,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
