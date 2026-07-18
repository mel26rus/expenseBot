package model

type Session struct {
	UserID        int64
	State         string
	Payload       []byte
	EditMessageId int64
}
