package zkhttp

import (
	"fmt"
	zkCommon "github.com/zerok-ai/zk-utils-go/common"
	zkHttp "github.com/zerok-ai/zk-utils-go/http"
	zkLogger "github.com/zerok-ai/zk-utils-go/logs"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// JsonExtractor This interface is used to extract key-values from the response
// Any custom value extraction logic can be provided here. The returned key-value map can be used by apiCalFlow
// further to execute the next set of APIs
type JsonExtractor interface {
	ExtractIdentifier(string, any) map[string]string
}

// ApiOperationValidator For now we don't need a validation as such
// We can throw an error if an expected cookie or key doesn't appear after an operation
// Can be extended further to add proper validations, if required
type ApiOperationValidator interface {
	Validate(map[string]http.Cookie, map[string]string) bool
}

type ApiCallOperation struct {
	Debug                           bool               `yaml:"debug"`
	ApiCallUnit                     ApiCallUnit        `yaml:"apiCall"`
	CookiesToBeExtracted            *[]string          `yaml:"cookiesExtract"`
	ResponseUrlParamsToBeExtracted  *[]string          `yaml:"responseParamsExtract"`
	ResponseJsonParamsToBeExtracted *[]string          `yaml:"responseJsonExtract"`
	HeaderUrlParamsToBeExtracted    *map[string]string `yaml:"headersParamExtract"`
}

type ApiCallFlow struct {
	ApiCallOperations *[]ApiCallOperation      `yaml:"apiCallOperations"`
	JsonExtractors    map[string]JsonExtractor `yaml:"-"`
	HardStopOnMiss    bool                     `yaml:"stopOnMiss"`
}

func (apiCallFlow ApiCallFlow) RegisterJsonExtractor(identifier string, jsonExtractor JsonExtractor) ApiCallFlow {
	if apiCallFlow.JsonExtractors == nil {
		apiCallFlow.JsonExtractors = map[string]JsonExtractor{}
	}
	apiCallFlow.JsonExtractors[identifier] = jsonExtractor

	return apiCallFlow
}

func (apiCallFlow ApiCallFlow) Execute(initParams map[string]any) (map[string]http.Cookie, map[string]string, *zkerrors.ZkError) {
	apiCallOperations := apiCallFlow.ApiCallOperations
	if apiCallOperations == nil {
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, "apiCallFlow.ApiCallOperations is null")
		return nil, nil, &zkError
	}

	var processedCookies = make(map[string]http.Cookie)
	var processedKeys = make(map[string]string)

	for index, apiCallOperation := range *apiCallOperations {
		zkLogger.Info(zkHttp.LOG_TAG, "$$$$$$$$$$$ STARTING OPERATION "+fmt.Sprint(index+1)+" $$$$$$$$$$$$")
		zkError := apiCallOperation.Execute(initParams, apiCallFlow.JsonExtractors,
			processedCookies, processedKeys,
			apiCallFlow.HardStopOnMiss)
		if zkError != nil {
			return nil, nil, zkError
		}
		zkLogger.Debug(zkHttp.LOG_TAG, "processedCookies", processedCookies)
		zkLogger.Debug(zkHttp.LOG_TAG, "processedKeys", processedKeys)
	}

	return processedCookies, processedKeys, nil
}

func (apiCallOperation ApiCallOperation) Execute(initParams map[string]any,
	jsonExtractors map[string]JsonExtractor,
	processedCookies map[string]http.Cookie,
	processedKeys map[string]string,
	hardStopOnMiss bool) *zkerrors.ZkError {
	zkApiUtils := zkHttp.Create()
	response, zkError := apiCallOperation.ApiCallUnit.Execute(initParams, processedCookies, processedKeys)
	if zkError != nil {
		return zkError
	}
	zkLogger.Debug(zkHttp.LOG_TAG, "response", response)
	// responseData2, err2 := ioutil.ReadAll(response.Body)
	// if err2 == nil {
	// 	zkLogger.Debug(zkhttpUtils.LOG_TAG, "responseData2", string(responseData2))
	// }else {
	// 	zkLogger.Debug(zkhttpUtils.LOG_TAG, "err2", err2)
	// }
	var cookiesToBeReturned *map[string]http.Cookie
	var urlParamsToBeReturned *map[string]string
	var jsonParamsToBeReturned *map[string]string
	var headerParamsToBeReturned *map[string]string

	if apiCallOperation.CookiesToBeExtracted != nil && len(*apiCallOperation.CookiesToBeExtracted) > 0 {
		if len(response.Cookies()) == 0 && hardStopOnMiss {
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, "Got empty cookies in response for "+
				apiCallOperation.ApiCallUnit.Url)
			return &zkError
		}

		extractedCookies := zkApiUtils.ExtractCookies(*apiCallOperation.CookiesToBeExtracted, response.Cookies())
		cookiesToBeReturned = &extractedCookies
		for k, v := range *cookiesToBeReturned {
			processedCookies[k] = v
		}
	}

	if apiCallOperation.ResponseUrlParamsToBeExtracted != nil {
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		val, err := url.QueryUnescape(string(responseData))
		for _, element := range *apiCallOperation.ResponseUrlParamsToBeExtracted {
			processedKeys[element] = zkApiUtils.ExtractUrlParam(element, val)
		}
	}

	if apiCallOperation.HeaderUrlParamsToBeExtracted != nil {
		for headerKey, keyToBeExtracted := range *apiCallOperation.HeaderUrlParamsToBeExtracted {
			headerValue := response.Header[headerKey]
			processedKeys[keyToBeExtracted] = zkApiUtils.ExtractUrlParam(keyToBeExtracted, headerValue[0])
		}
	}

	if apiCallOperation.ResponseJsonParamsToBeExtracted != nil {
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		for _, element := range *apiCallOperation.ResponseJsonParamsToBeExtracted {
			jsonExtractor := (jsonExtractors)[element]
			keyValueMap := jsonExtractor.ExtractIdentifier(element, responseData)

			for k, v := range keyValueMap {
				processedKeys[k] = v
			}
		}
	}

	return nil
}

func (apiCallUnit ApiCallUnit) Execute(initParams map[string]any, processedCookies map[string]http.Cookie,
	processedKeys map[string]string) (*http.Response, *zkerrors.ZkError) {
	request, url, cookies := apiCallUnit.Initialize(initParams, processedCookies, processedKeys)
	zkLogger.Debug(zkHttp.LOG_TAG, "url", url)
	zkLogger.Debug(zkHttp.LOG_TAG, "cookies", cookies)
	if apiCallUnit.ContentType != nil {
		zkLogger.Debug(zkHttp.LOG_TAG, "apiCallUnit.ContentType", apiCallUnit.ContentType)
	}
	var bodyReader io.Reader
	if request != "" {
		bodyReader = strings.NewReader(request)
	}

	httpResponse, zkErr := zkHttp.Create().
		WithInsecure(zkCommon.ToPtr[bool](false)).
		WithRedirect(false).
		WithContentType(apiCallUnit.ContentType).
		WithCookies(&cookies).
		Go(*apiCallUnit.Method, url, bodyReader)
	return httpResponse, zkErr
}

// Initialize The urls, request params and request body - everything can be parameterised
// In this function, we sanitise everything by replacing the parameters with the appropriate values
// before starting the execution
func (apiCallUnit ApiCallUnit) Initialize(initParams map[string]any, processedCookies map[string]http.Cookie,
	processedKeys map[string]string) (string, string, []http.Cookie) {
	rawCookies := apiCallUnit.Cookies
	existingCookies := apiCallUnit.ExistingCookies
	lenRawCookies := 0
	if rawCookies != nil && len(*rawCookies) > 0 {
		lenRawCookies = len(*rawCookies)
	}
	cookies := make([]http.Cookie, lenRawCookies)

	if rawCookies != nil && len(*rawCookies) > 0 {
		for index, rawCookie := range *rawCookies {
			httpCookie := http.Cookie{}
			httpCookie.Name = rawCookie.Name
			httpCookie.Value = rawCookie.Value
			httpCookie.Expires = rawCookie.Expires
			cookies[index] = httpCookie
		}
	}

	if existingCookies != nil && len(*existingCookies) > 0 {
		zkLogger.Debug(zkHttp.LOG_TAG, "existingCookies", existingCookies)
		for _, existingCookie := range *existingCookies {
			cookies = append(cookies, processedCookies[existingCookie])
			// cookies[index] = processedCookies[existingCookie]
		}
	}

	//Process Url
	url := apiCallUnit.Url
	url = apiCallUnit.SanitiseString(url, initParams, processedKeys)

	//Process Request Body
	var request string = ""
	if apiCallUnit.RequestBody != nil {
		request = *apiCallUnit.RequestBody
		request = apiCallUnit.SanitiseString(request, initParams, processedKeys)
	}

	return request, url, cookies
}

func (apiCallUnit ApiCallUnit) SanitiseString(stringToBeSanitised string, initParams map[string]any, 
		processedKeys map[string]string) string {
	//https://stackoverflow.com/a/40586418/4666116
	rex := regexp.MustCompile(`{{[^}]+}}`)
	allParams := rex.FindAllStringSubmatch(stringToBeSanitised, -1)

	for _, param := range allParams {
		paramName := param[0]
		paramName = strings.Replace(paramName, "{", "", -1)
		paramName = strings.Replace(paramName, "}", "", -1)

		paramValue, paramValuePresent := processedKeys[paramName]
		if !paramValuePresent {
			paramValue, _ = initParams[paramName].(string)
		}
		stringToBeSanitised = strings.Replace(stringToBeSanitised, "{{"+paramName+"}}", paramValue, -1)
	}
	return stringToBeSanitised;
}

