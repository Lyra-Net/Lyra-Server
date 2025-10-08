package dto

type UserRegistered struct {
	UserID    string `json:"user_id"`
	DeviceID  string `json:"device_id"`
	Browser   string `json:"browser"`
	OS        string `json:"os"`
	UserIP    string `json:"user_ip"`
	Timestamp int64  `json:"timestamp"`
}
