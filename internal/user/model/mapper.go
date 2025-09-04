package model

// Map request register
func ToUserModel(r RegisterRequest) User {
	return User{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password, 
	}
}

func ToUserResponse(u User) UserResponse {
	return UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func ToUserProfileResponse(u User) UserProfileResponse {
	return UserProfileResponse{
		Name:  u.Name,
		Email: u.Email,
	}
}

