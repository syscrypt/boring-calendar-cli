package main

type Calendar struct {
	UserToken string   `json:"user_token"`
	Date      int64    `json:"date"`
	Events    []*Event `json:"events"`
}

type Event struct {
	Title     string `json:"title"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
}
