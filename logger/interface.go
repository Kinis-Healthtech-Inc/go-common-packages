package logger

import "time"

type UserInfo struct {
	UserID           string
	Role             string
	OrganizationID   string
	OrganizationName string
}

type LogFields struct {
	Context     string  `json:"context"`
	Event       string  `json:"event"`
	LogType     LogType `json:"log_type"` // request_log or application_log
	ServiceName string  `json:"service_name"`
}
type RequestLogFields struct {
	LogFields
	Timestamp      string        `json:"timestamp"`
	Endpoint       string        `json:"endpoint"`
	Method         string        `json:"method"`
	StatusCode     int           `json:"status_code"`
	QueryParams    any           `json:"query_params,omitempty"`
	RequestID      string        `json:"request_id,omitempty"`
	URL            string        `json:"url"`
	ClientIP       string        `json:"client_ip,omitempty"`
	UserAgent      string        `json:"user_agent,omitempty"`
	UserID         string        `json:"user_id,omitempty"`
	OrganizationID string        `json:"organization_id,omitempty"`
	Latency        time.Duration `json:"latency"`
	ErrorMessage   string        `json:"error_message,omitempty"`
}

type ApplicationLogFields struct {
	UserInfo
	LogFields
	ErrorType    string `json:"error_type,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
