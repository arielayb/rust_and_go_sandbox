package app

const (
	USER_INFO = "USER_INFO"
)

type UserWebInfo struct {
	UserUUID string `json:"user_uuid"`
	Method   string `json:"method"`
	Message  string `json:"msg"`
}
