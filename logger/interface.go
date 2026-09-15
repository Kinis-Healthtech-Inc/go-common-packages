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

type RequestLogFieldsBuilder struct {
	fields RequestLogFields
}
type ApplicationLogFieldsBuilder struct {
	fields ApplicationLogFields
}

func NewRequestLogFields() *RequestLogFieldsBuilder {
	return &RequestLogFieldsBuilder{
		fields: RequestLogFields{
			LogFields: LogFields{
				LogType: RequestLogType,
			},
		},
	}
}
func NewApplicationLogFields() *ApplicationLogFieldsBuilder {
	return &ApplicationLogFieldsBuilder{
		fields: ApplicationLogFields{
			LogFields: LogFields{
				LogType: ApplicationLogType,
			},
		},
	}
}

func (b *RequestLogFieldsBuilder) WithContext(context string) *RequestLogFieldsBuilder {
	b.fields.Context = context
	return b
}

func (b *RequestLogFieldsBuilder) WithEvent(event string) *RequestLogFieldsBuilder {
	b.fields.Event = event
	return b
}

func (b *RequestLogFieldsBuilder) WithServiceName(serviceName string) *RequestLogFieldsBuilder {
	b.fields.ServiceName = serviceName
	return b
}

func (b *RequestLogFieldsBuilder) WithTimestamp(timestamp string) *RequestLogFieldsBuilder {
	b.fields.Timestamp = timestamp
	return b
}

func (b *RequestLogFieldsBuilder) WithEndpoint(endpoint string) *RequestLogFieldsBuilder {
	b.fields.Endpoint = endpoint
	return b
}

func (b *RequestLogFieldsBuilder) WithMethod(method string) *RequestLogFieldsBuilder {
	b.fields.Method = method
	return b
}

func (b *RequestLogFieldsBuilder) WithStatusCode(statusCode int) *RequestLogFieldsBuilder {
	b.fields.StatusCode = statusCode
	return b
}

func (b *RequestLogFieldsBuilder) WithQueryParams(queryParams any) *RequestLogFieldsBuilder {
	b.fields.QueryParams = queryParams
	return b
}

func (b *RequestLogFieldsBuilder) WithRequestID(requestID string) *RequestLogFieldsBuilder {
	b.fields.RequestID = requestID
	return b
}

func (b *RequestLogFieldsBuilder) WithURL(url string) *RequestLogFieldsBuilder {
	b.fields.URL = url
	return b
}

func (b *RequestLogFieldsBuilder) WithClientIP(clientIP string) *RequestLogFieldsBuilder {
	b.fields.ClientIP = clientIP
	return b
}

func (b *RequestLogFieldsBuilder) WithUserAgent(userAgent string) *RequestLogFieldsBuilder {
	b.fields.UserAgent = userAgent
	return b
}

func (b *RequestLogFieldsBuilder) WithUserID(userID string) *RequestLogFieldsBuilder {
	b.fields.UserID = userID
	return b
}

func (b *RequestLogFieldsBuilder) WithOrganizationID(organizationID string) *RequestLogFieldsBuilder {
	b.fields.OrganizationID = organizationID
	return b
}

func (b *RequestLogFieldsBuilder) WithLatency(latency time.Duration) *RequestLogFieldsBuilder {
	b.fields.Latency = latency
	return b
}

func (b *RequestLogFieldsBuilder) WithErrorMessage(errorMessage string) *RequestLogFieldsBuilder {
	b.fields.ErrorMessage = errorMessage
	return b
}

func (b *RequestLogFieldsBuilder) Build() RequestLogFields {
	return b.fields
}

func (b *ApplicationLogFieldsBuilder) WithContext(context string) *ApplicationLogFieldsBuilder {
	b.fields.Context = context
	return b
}

func (b *ApplicationLogFieldsBuilder) WithEvent(event string) *ApplicationLogFieldsBuilder {
	b.fields.Event = event
	return b
}

func (b *ApplicationLogFieldsBuilder) WithServiceName(serviceName string) *ApplicationLogFieldsBuilder {
	b.fields.ServiceName = serviceName
	return b
}

func (b *ApplicationLogFieldsBuilder) WithUserID(userID string) *ApplicationLogFieldsBuilder {
	b.fields.UserID = userID
	return b
}

func (b *ApplicationLogFieldsBuilder) WithRole(role string) *ApplicationLogFieldsBuilder {
	b.fields.Role = role
	return b
}

func (b *ApplicationLogFieldsBuilder) WithOrganizationID(organizationID string) *ApplicationLogFieldsBuilder {
	b.fields.OrganizationID = organizationID
	return b
}

func (b *ApplicationLogFieldsBuilder) WithOrganizationName(organizationName string) *ApplicationLogFieldsBuilder {
	b.fields.OrganizationName = organizationName
	return b
}

func (b *ApplicationLogFieldsBuilder) WithUserInfo(userInfo UserInfo) *ApplicationLogFieldsBuilder {
	b.fields.UserInfo = userInfo
	return b
}

func (b *ApplicationLogFieldsBuilder) WithErrorType(errorType string) *ApplicationLogFieldsBuilder {
	b.fields.ErrorType = errorType
	return b
}

func (b *ApplicationLogFieldsBuilder) WithErrorMessage(errorMessage string) *ApplicationLogFieldsBuilder {
	b.fields.ErrorMessage = errorMessage
	return b
}

func (b *ApplicationLogFieldsBuilder) Build() ApplicationLogFields {
	return b.fields
}
