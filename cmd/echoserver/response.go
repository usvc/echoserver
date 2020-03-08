package main

import (
	"time"
)

type Cookie struct {
	Name     string `json:"cookie"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Expires  string `json:"expires"`
	Secure   bool   `json:"secure"`
	HTTPOnly bool   `json:"http_only"`
}

type Metadata struct {
	DurationMs         int64     `json:"duration_ms"`
	ReceivedTimestamp  time.Time `json:"received_timestamp"`
	RespondedTimestamp time.Time `json:"responded_timestamp"`
}

func (m *Metadata) RequestReceived() {
	m.ReceivedTimestamp = time.Now()
}

func (m *Metadata) RequestProcessingCompleted() {
	m.RespondedTimestamp = time.Now()
	m.DurationMs = time.Since(m.ReceivedTimestamp).Microseconds()
}

type Response struct {
	ID       string   `json:"id"`
	Request  Request  `json:"request"`
	Errors   []string `json:"errors"`
	Metadata Metadata `json:"metadata"`
}

func NewResponse(id string) *Response {
	response := &Response{
		Errors:   []string{},
		ID:       id,
		Metadata: Metadata{},
		Request:  Request{},
	}
	response.Metadata.RequestReceived()
	return response
}

type Request struct {
	Body       interface{}            `json:"body"`
	Cookies    []Cookie               `json:"cookies"`
	Header     map[string]interface{} `json:"header"`
	Hostname   *string                `json:"hostname"`
	Form       map[string]interface{} `json:"form"`
	Method     *string                `json:"method"`
	Password   *string                `json:"password"`
	Path       *string                `json:"path"`
	Protocol   *string                `json:"protocol"`
	Query      map[string]interface{} `json:"query"`
	Referer    *string                `json:"referer"`
	RemoteAddr *string                `json:"remote_addr"`
	Size       int64                  `json:"size"`
	UserAgent  *string                `json:"user_agent"`
	Username   *string                `json:"username"`
}
