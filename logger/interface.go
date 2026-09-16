package logger

import "time"

const (
	// --- UserInfo Keys ---
	LogKeyUserID           = "user_id"
	LogKeyRole             = "role"
	LogKeyOrganizationID   = "organization_id"
	LogKeyOrganizationName = "organization_name"

	// --- LogFields Keys ---
	LogKeyContext     = "context"
	LogKeyEvent       = "event"
	LogKeyLogType     = "log_type"
	LogKeyServiceName = "service_name"

	// --- RequestLogFields Unique Keys ---
	LogKeyTimestamp   = "timestamp"
	LogKeyEndpoint    = "endpoint"
	LogKeyMethod      = "method"
	LogKeyStatusCode  = "status_code"
	LogKeyQueryParams = "query_params"
	LogKeyRequestID   = "request_id"
	LogKeyURL         = "url"
	LogKeyClientIP    = "client_ip"
	LogKeyUserAgent   = "user_agent"
	LogKeyLatency     = "latency"

	// --- ApplicationLogFields / Shared Keys ---
	LogKeyErrorType    = "error_type"
	LogKeyErrorMessage = "error_message"
	LogKeySourceFile   = "custom_source_file"
	LogKeySourceLine   = "custom_source_line"
)

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
