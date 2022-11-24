package base_info

type ApiUserInfo struct {
	UserID      string `json:"userID" binding:"required,min=1,max=64" swaggo:"true,用户ID,"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"omitempty,min=1,max=64" swaggo:"true,password,"`
	Nickname    string `json:"nickname" binding:"omitempty,min=1,max=64" swaggo:"true,my id,19"`
	FaceURL     string `json:"faceURL" binding:"omitempty,max=1024"`
	Gender      int32  `json:"gender" binding:"omitempty,oneof=-1 0 1 2 3"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,max=32"`
	Birth       uint32 `json:"birth" binding:"omitempty"`
	CreateTime  int64  `json:"createTime"`
	LoginLimit  int32  `json:"loginLimit" binding:"omitempty"`
	Ex          string `json:"ex" binding:"omitempty,max=1024"`
}
