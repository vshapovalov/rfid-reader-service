package utils

import "strings"

const readerIdPlaceholder = "{reader_id}"

var (
	statusTopic   = "reader/" + readerIdPlaceholder + "/status"
	cardReadTopic = "reader/" + readerIdPlaceholder + "/card"
	buzzerTopic   = "reader/" + readerIdPlaceholder + "/buzzer"
)

func GetStatusTopic(readerId string) string {
	return strings.Replace(statusTopic, readerIdPlaceholder, readerId, 1)
}

func GetCardReadTopic(readerId string) string {
	return strings.Replace(cardReadTopic, readerIdPlaceholder, readerId, 1)
}

func GetBuzzerTopic(readerId string) string {
	return strings.Replace(buzzerTopic, readerIdPlaceholder, readerId, 1)
}
