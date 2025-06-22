package app

type UserWebInfo struct {
	UserUUID string `json:"user_uuid"`
	Method   string `json:"method"`
	AlertMsg string `json:"alert_msg"`
}
