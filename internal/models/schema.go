package models

type User struct {
	ID   int    `column:"id" json:"id"`
	Name string `column:"name" json:"name"`
}

type Notification struct {
	FromID  int    `column:"from_id" json:"fromID"`
	ToID    int    `column:"to_id" json:"toID"`
	Message string `column:"message" json:"message"`
}
