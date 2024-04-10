package domain

type Dialog struct {
	User1       *User            `json:"user1"`
	User2       *User            `json:"user2"`
	LastMessage *PersonalMessage `json:"lastMessage"`
}
