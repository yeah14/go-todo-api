package request

import "errors"

type UpdateProfileRequest struct {
	Username *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email,max=100"`
	Avatar   *string `json:"avatar,omitempty" binding:"omitempty,url,max=255"`
}

func (req *UpdateProfileRequest) HasUpdate() bool {
	if req.Username != nil || req.Email != nil || req.Avatar != nil {
		return true
	}
	return false
}

func (req *UpdateProfileRequest) ToMap() map[string]string {
	update := make(map[string]string)
	if req.Username != nil {
		update["username"] = *req.Username
	}
	if req.Email != nil {
		update["email"] = *req.Email
	}
	if req.Avatar != nil {
		update["avatar"] = *req.Avatar
	}
	return update
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" binding:"required,equal=NewPassword"`
}

func (req *ChangePasswordRequest) Validate() error {
	if req.OldPassword == req.ConfirmPassword {
		return errors.New("新密码不能与旧密码相同")
	}
	return nil
}
