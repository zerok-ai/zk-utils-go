package zkhttp

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	zklogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"io/ioutil"
	"net/http"
	"reflect"
)

var LOG_TAG = "zkhttp"

const (
	HTTP_UTILS_ZK_RESPONSE         = "zkHttpResponse"
	HTTP_UTILS_AUTH_TOKEN          = "auth_token"
	HTTP_UTILS_OPERATOR_AUTH_TOKEN = "Operator_auth_token"
	HTTP_UTILS_TOKEN               = "Token"
	HTTP_UTILS_CLUSTER_ID          = "Cluster_id"
	HTTP_UTILS_REQUEST_IN_BYTES    = "requestBodyInBytes"
	HTTP_UTILS_REQUEST_INTERFACE   = "requestBodyInterface"
)

func GenerateHttpCookie(name string, value string) *http.Cookie {
	return &http.Cookie{
		Name:  name,
		Value: value,
	}
}

func ReadRequestBody(ctx iris.Context) {
	var zkHttpResponse ZkHttpResponse[any]
	request := ctx.Request()
	requestBodyInBytes, err := ioutil.ReadAll(request.Body)

	if err != nil {

		zklogger.Debug(LOG_TAG, "Error while reading request body")
		zklogger.Error(LOG_TAG, err)
		zkHttpResponse = ZkHttpResponseBuilder[any]{}.WithZkErrorType(zkerrors.ZkErrorBadRequest).
			Debug("actual", "Error while reading request body").
			Build()
		ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
		// c.StopWithProblem(iris.StatusInternalServerError, iris.NewProblem().
		// 	Title(err.Error()))
		return
	}
	ctx.Values().Set(HTTP_UTILS_REQUEST_IN_BYTES, requestBodyInBytes)

	ctx.Next()
}

func ValidateRequestBody(ctx iris.Context, bodyType reflect.Type) {
	if bodyType.Kind() == reflect.Ptr {
		bodyType = bodyType.Elem()
	}
	bodyInterface := reflect.New(bodyType).Interface()

	err := ctx.ReadJSON(bodyInterface)
	if err != nil {
		zklogger.Error(LOG_TAG, err)
		zkHttpResponse := ZkHttpResponseBuilder[any]{}.WithZkErrorType(zkerrors.ZkErrorBadRequest).
			Debug("actual", "Error while reading request body").
			Build()
		ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
		return
	} else {
		var validate = validator.New()
		err := validate.Struct(bodyInterface)
		var zkErrorParamToMessage zkerrors.ZkErrorParamToMessage = map[string]string{}
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				message := zkerrors.MessageFromValidation(err)
				// log.Println("err.Field()==", err.Field())
				// log.Println("err.Param()==", err.Tag())
				zkErrorParamToMessage[err.Field()] = message
			}
		}
		if len(zkErrorParamToMessage) > 0 {
			zkHttpResponse := ZkHttpResponseBuilder[any]{}.WithZkErrorType(zkerrors.ZkErrorBadRequest).
				ErrorInfo("validations", zkErrorParamToMessage).
				Build()
			ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
			return
		}
		ctx.Values().Set(HTTP_UTILS_REQUEST_INTERFACE, bodyInterface)
	}

	ctx.Next()
}

func ValidateObject(ctx iris.Context, s interface{}) {
	var validate = validator.New()
	err := validate.Struct(s)
	var zkErrorParamToMessage zkerrors.ZkErrorParamToMessage = map[string]string{}
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			message := zkerrors.MessageFromValidation(err)
			zkErrorParamToMessage[err.Field()] = message
		}
	}
	if len(zkErrorParamToMessage) > 0 {
		zkHttpResponse := ZkHttpResponseBuilder[any]{}.WithZkErrorType(zkerrors.ZkErrorBadRequest).
			ErrorInfo("validations", zkErrorParamToMessage).
			Build()
		ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
		return
	}
	// ctx.Values().Set(HTTP_UTILS_REQUEST_INTERFACE, s)

	ctx.Next()
}

func Validate(ctx iris.Context, bodyType reflect.Type) {
	if bodyType.Kind() == reflect.Ptr {
		bodyType = bodyType.Elem()
	}
	bodyInterface := reflect.New(bodyType).Interface()

	err := ctx.ReadJSON(bodyInterface)
	if err != nil {
		zklogger.Error(LOG_TAG, err)
		zkHttpResponse := ZkHttpResponseBuilder[any]{}.WithZkErrorType(zkerrors.ZkErrorBadRequest).
			Debug("actual", "Error while reading request body").
			Build()
		ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
		return
	} else {
		var validate = validator.New()
		err := validate.Struct(bodyInterface)
		var zkErrorParamToMessage zkerrors.ZkErrorParamToMessage = map[string]string{}
		if err != nil {
			for _, err := range err.(validator.ValidationErrors) {
				message := zkerrors.MessageFromValidation(err)
				// log.Println("err.Field()==", err.Field())
				// log.Println("err.Param()==", err.Tag())
				zkErrorParamToMessage[err.Field()] = message
			}
		}
		if len(zkErrorParamToMessage) > 0 {
			zkHttpResponse := ZkHttpResponseBuilder[any]{}.WithZkErrorType(zkerrors.ZkErrorBadRequest).
				ErrorInfo("validations", zkErrorParamToMessage).
				Build()
			ctx.StopWithJSON(zkHttpResponse.Status, zkHttpResponse)
			return
		}
		ctx.Values().Set(HTTP_UTILS_REQUEST_INTERFACE, bodyInterface)
	}

	ctx.Next()
}

func ToZkResponse[T any](status int, payload T, rawResponse any, zkError *zkerrors.ZkError) ZkHttpResponse[T] {
	zkHttpResponseBuilder := ZkHttpResponseBuilder[T]{}
	var zkHttpResponse ZkHttpResponse[T]
	if zkError != nil {
		zkHttpResponseBuilder := zkHttpResponseBuilder.WithZkErrorType(zkError.Error).
			// Message(&zkError.Error.Message).
			Debug("rawResponse", rawResponse)
		if zkError.Metadata != nil {
			//TODO: Move zkHttpResponseBuilder.Debug to use pointer
			zkHttpResponseBuilder = zkHttpResponseBuilder.ErrorInfo("error", zkError.Metadata)
		}
		zkHttpResponse = zkHttpResponseBuilder.Build()
	} else {
		zkHttpResponse = zkHttpResponseBuilder.WithStatus(status).
			// Message(nil).
			Data(&payload).
			Debug("rawResponse", rawResponse).
			Build()
	}
	return zkHttpResponse
}
