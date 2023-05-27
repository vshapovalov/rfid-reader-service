package services

import "encoding/json"

const StatusOnline = "online"
const StatusOffline = "offline"

type BuzzerCount struct {
	Count int `json:"count"`
}

type StatusInfo struct {
	ReaderId  string `json:"readerId"`
	Status    string `json:"status"`
	ExtraInfo string `json:"extraInfo,omitempty"`
}

func NewStatusWithExtraInfoMessage(readerId string, status string, extraInfo string) *StatusInfo {
	return &StatusInfo{ReaderId: readerId, Status: status, ExtraInfo: extraInfo}
}

func NewStatusMessage(readerId string, status string) *StatusInfo {
	return &StatusInfo{ReaderId: readerId, Status: status}
}

func (m *StatusInfo) ToByteArray() []byte {
	marshal, _ := json.Marshal(m)
	return marshal
}

type CardReadInfo struct {
	CardNumber []string `json:"cardNumbers"`
	ReaderId   string   `json:"readerId"`
}

func NewCardReadInfo(cardNumbers []string, readerId string) *CardReadInfo {
	return &CardReadInfo{CardNumber: cardNumbers, ReaderId: readerId}
}

func (m *CardReadInfo) ToByteArray() []byte {
	marshal, _ := json.Marshal(m)
	return marshal
}
