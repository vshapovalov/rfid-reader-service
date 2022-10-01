package services

import "encoding/json"

const StatusOnline = "online"
const StatusOffline = "offline"

type StatusInfo struct {
	ReaderId string `json:"readerId"`
	Status   string `json:"status"`
}

func NewStatusInfo(readerId string, status string) *StatusInfo {
	return &StatusInfo{ReaderId: readerId, Status: status}
}

func (m *StatusInfo) ToByteArray() []byte {
	marshal, _ := json.Marshal(m)
	return marshal
}

type CardReadInfo struct {
	CardNumber string `json:"cardNumber"`
	ReaderId   string `json:"readerId"`
}

func NewCardReadInfo(cardNumber, readerId string) *CardReadInfo {
	return &CardReadInfo{CardNumber: cardNumber, ReaderId: readerId}
}

func (m *CardReadInfo) ToByteArray() []byte {
	marshal, _ := json.Marshal(m)
	return marshal
}
