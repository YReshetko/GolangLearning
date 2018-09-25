package domain

import (
	"bytes"
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
	Buffer    *bytes.Buffer
}
