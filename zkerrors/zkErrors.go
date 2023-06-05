package zkerrors

import (
	"github.com/kataras/iris/v12"
	"net/http"
)

type ZkErrorType struct {
	Message  string `json:"message"`
	Type     string `json:"type"`
	Status   int    `json:"status"`
	Metadata any    `json:"metadata"`
}

type ZkError struct {
	Error    ZkErrorType `json:"error"`
	Metadata any         `json:"metadata"`
}

type ZkErrorBuilder struct {
}

var (
	ZkErrorInternalServer = ZkErrorType{Status: iris.StatusInternalServerError,
		Type:    StatusText(iris.StatusInternalServerError),
		Message: "Encountered an issue, contact support"}
	ZkErrorTimeout = ZkErrorType{Status: iris.StatusRequestTimeout,
		Type:    StatusText(iris.StatusRequestTimeout),
		Message: "Encountered an issue, contact support"}
	ZkErrorNotFound = ZkErrorType{Status: iris.StatusNotFound,
		Type:    StatusText(iris.StatusNotFound),
		Message: "Encountered an issue, contact support"}
	ZkErrorSessionExpired = ZkErrorType{Status: iris.StatusPageExpired,
		Type:    StatusText(iris.StatusPageExpired),
		Message: "The session has expired. Please login again to continue"}
	ZkErrorUnauthorized = ZkErrorType{Status: iris.StatusUnauthorized,
		Type:    StatusText(iris.StatusUnauthorized),
		Message: "You are unauthorized to perform this operation. Contact system admin"}
	ZkErrorBadRequest = ZkErrorType{Status: iris.StatusBadRequest,
		Type:    StatusText(iris.StatusBadRequest),
		Message: "The request was malformed or contained invalid parameters"}
	ZkErrorInterEmailServer = ZkErrorType{Status: iris.StatusInternalServerError,
		Type:    StatusText(iris.StatusInternalServerError),
		Message: "Encountered an issue while sending email, contact support"}
	ZkErrorDbError = ZkErrorType{Status: iris.StatusInternalServerError,
		Type:    StatusText(iris.StatusInternalServerError),
		Message: "Encountered an issue while executing db operation"}
	ZkErrorBadRequestLimitIsNotInteger  = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "LIMIT is not an integer"}
	ZkErrorBadRequestOffsetIsNotInteger = ZkErrorType{Status: iris.StatusBadRequest, Type: "BAD_REQUEST", Message: "OFFSET is not an integer"}
)

func (zkError ZkError) SetMetadata(metadata any) {
	zkError.Metadata = metadata
}

var localStatusText = map[int]string{
	iris.StatusPageExpired: "Session Expired",
}

func StatusText(status int) string {
	statusText := http.StatusText(status)
	if statusText == "" {
		statusText = localStatusText[status]
	}

	return statusText
}

func (zkErrorBuilder ZkErrorBuilder) Build(zkErrorType ZkErrorType, metadata any) ZkError {
	return ZkError{
		Error:    zkErrorType,
		Metadata: metadata,
	}
}

func (zkErrorBuilder ZkErrorBuilder) CreateZkErrorType(status int, message string) ZkErrorType {
	return ZkErrorType{
		Status:  status,
		Type:    StatusText(status),
		Message: message,
	}
}
