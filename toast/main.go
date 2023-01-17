package main

import (
	"github.com/go-toast/toast"
	"log"
)

func main() {
	var notification = toast.Notification{
		AppID:   "Touch Fish",
		Title:   "Word",
		Message: "word",
		Icon:    "C:\\Users\\HELEN.WANG\\Pictures\\Camera Roll\\empty.png", // This file must exist (remove this line if it doesn't)
		Actions: []toast.Action{
			{"protocol", "记住了！", ""},
			{"protocol", "记错了！", ""},
			{"protocol", "不记得！", ""},
		},
	}
	err := notification.Push()
	if err != nil {
		log.Fatalln(err)
	}
}
