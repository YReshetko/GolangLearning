package domain

import (
	"io"
	"time"
)

type EmailData struct {
	Uid            uint32
	Date           time.Time
	CallType       string
	WhoCalls       string
	InputNumber    string
	Participant    string
	Duration       string
	RecordFileName string
}

type EmailToSave struct {
	EmailData EmailData
	Reader    io.Reader
}
