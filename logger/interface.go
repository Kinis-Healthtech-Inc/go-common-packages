package logger

import "time"

// Example of How to use
//handler.logger.Logger(ctx).
//	With(logger.LogKeyContext, "patient calls GetUserAssessmentSessionByID").
//	With(logger.LogKeyErrorType, "database error").
//	With(logger.LogKeyEvent, "assessment_session.get").
//	With(logger.LogKeyErrorMessage, "error getting assessment session: "+err.Error()).
//	InfoContext(ctx.Context(), "Failed to get assessment session")

const (
	// --- Service Keys --- Almost all services will have these keys
	LogKeyLogType     = "log_type"
	LogKeyServiceName = "service_name"
	// --- UserInfo Keys --- Almost all services will have these keys
	LogKeyUserID           = "user_id"
	LogKeyRole             = "role"
	LogKeyOrganizationID   = "organization_id"
	LogKeyOrganizationName = "organization_name"

	// --- LogFields Keys ---
	LogKeyContext      = "context"       // explain what the log is for, e.g. "patient calls GetUserAssessmentSessionByID"
	LogKeyEvent        = "event"         // explain what the log is about, <some_entity>.<action> e.g. "assessment_session.get"
	LogKeyErrorType    = "error_type"    // e.g. "database error", "validation error", "third party api error"
	LogKeyErrorMessage = "error_message" // says what the error is, e.g. "error getting assessment session: invalid id"

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
	LogKeySourceFile = "custom_source_file"
	LogKeySourceLine = "custom_source_line"
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
