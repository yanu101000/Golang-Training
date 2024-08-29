package entity

import "time"

type Transaction struct {
	FromUserID User      `json:"fromUserID"`
	ToUserID   User      `json:"toUserID"`
	Nominal    int       `json:"nominal"`
	Timestamp  time.Time `json:"timestamp"`
}
