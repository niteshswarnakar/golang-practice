package models

type Permissions struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Subject  string `json:"subject"`
	Resource string `json:"resource"`
	Action   string `json:"action"`
}
