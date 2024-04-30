package dtos

type MessageDtos struct {
	FromID  int    `json:"fromID"`
	ToID    int    `json:"toID"`
	Message string `json:"message"`
}
