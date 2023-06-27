package zkhttp

import (
	"crypto/tls"
	"fmt"
	"github.com/zerok-ai/zk-utils-go/zkerrors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type zkApiUtils struct {
	Insecure         *bool
	Redirect         *bool
	TimeoutInSeconds time.Duration
	ContentType      *string
	Cookies          *[]http.Cookie
	Headers          *map[string]string
}

func Create() zkApiUtils {
	return CreateWithTimeout(10)
}

func CreateWithTimeout(duration time.Duration) zkApiUtils {
	return zkApiUtils{TimeoutInSeconds: duration}
}

func (zkApiUtils zkApiUtils) WithContentType(contentType *string) zkApiUtils {
	if contentType != nil {
		zkApiUtils.ContentType = contentType
	}
	return zkApiUtils
}

func (zkApiUtils zkApiUtils) WithInsecure(insecure *bool) zkApiUtils {
	if insecure != nil {
		zkApiUtils.Insecure = insecure
	}
	return zkApiUtils
}

func (zkApiUtils zkApiUtils) WithRedirect(redirect bool) zkApiUtils {
	zkApiUtils.Redirect = &redirect
	return zkApiUtils
}

func (zkApiUtils zkApiUtils) WithCookies(cookies *[]http.Cookie) zkApiUtils {
	if cookies != nil {
		zkApiUtils.Cookies = cookies
	}
	return zkApiUtils
}

func (zkApiUtils zkApiUtils) Header(key string, value string) zkApiUtils {
	if zkApiUtils.Headers == nil {
		zkApiUtils.Headers = &map[string]string{}
	}
	(*zkApiUtils.Headers)[key] = value
	return zkApiUtils
}

func (zkApiUtils zkApiUtils) Cookie(cookie http.Cookie) zkApiUtils {
	if zkApiUtils.Cookies == nil {
		zkApiUtils.Cookies = &[]http.Cookie{}
	}
	cookies := append(*zkApiUtils.Cookies, cookie)
	zkApiUtils.Cookies = &cookies
	return zkApiUtils
}

func (zkApiUtils zkApiUtils) Go(method string, urlToBeCalled string,
	requestBody io.Reader) (*http.Response, *zkerrors.ZkError) {
	return zkApiUtils.makeRawApiCallV3(method, urlToBeCalled, requestBody)
}

func (zkApiUtils zkApiUtils) Post(urlToBeCalled string,
	requestBody io.Reader) (*http.Response, *zkerrors.ZkError) {
	return zkApiUtils.Go("POST", urlToBeCalled, requestBody)
}

func (zkApiUtils zkApiUtils) Get(urlToBeCalled string) (*http.Response, *zkerrors.ZkError) {
	return zkApiUtils.Go("GET", urlToBeCalled, nil)
}

func (zkApiUtils zkApiUtils) Delete(urlToBeCalled string, requestBody io.Reader) (*http.Response, *zkerrors.ZkError) {
	return zkApiUtils.Go("DELETE", urlToBeCalled, requestBody)
}

func (zkApiUtils zkApiUtils) Put(urlToBeCalled string, requestBody io.Reader) (*http.Response, *zkerrors.ZkError) {
	return zkApiUtils.Go("PUT", urlToBeCalled, requestBody)
}

func (zkApiUtils zkApiUtils) makeRawApiCallV3(method string, urlToBeCalled string,
	requestBody io.Reader) (*http.Response, *zkerrors.ZkError) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: zkApiUtils.Insecure != nil && *zkApiUtils.Insecure},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * zkApiUtils.TimeoutInSeconds,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if zkApiUtils.Redirect == nil || *zkApiUtils.Redirect {
				return nil
			} else {
				return http.ErrUseLastResponse
			}
		},
	}

	return zkApiUtils.MakeRawApiCall(method, zkApiUtils.ContentType, *client, urlToBeCalled,
		zkApiUtils.Cookies, requestBody)
}

func (zkApiUtils zkApiUtils) MakeRawApiCall(method string, contentType *string, client http.Client,
	urlToBeCalled string, cookiesTobeAdded *[]http.Cookie,
	requestBody io.Reader) (*http.Response, *zkerrors.ZkError) {

	req, err := http.NewRequest(method, urlToBeCalled, requestBody)
	if err != nil {
		log.Fatalf("Got error %s", err.Error())
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, nil)
		return nil, &zkError
	}

	if cookiesTobeAdded != nil {
		for _, element := range *cookiesTobeAdded {
			req.AddCookie(&element)
		}
	}

	if contentType != nil {
		req.Header.Add("Content-Type", *contentType)
	}

	if zkApiUtils.Headers != nil {
		for key, value := range *zkApiUtils.Headers {
			req.Header.Add(key, value)
		}
	}

	response, err := client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorTimeout, "Request timed out - "+urlToBeCalled)
			return nil, &zkError
		}

		errorString := fmt.Sprintf("Unknown error occured while calling the API %v ", err)
		zkError := zkerrors.ZkErrorBuilder{}.Build(zkerrors.ZkErrorInternalServer, errorString)
		return nil, &zkError
	}

	return response, nil
}

func (zkApiUtils zkApiUtils) MakeApiCall(client http.Client, urlToBeCalled string, cookiesTobeAdded []http.Cookie, cookiesToBeExtracted []string,
	urlParamsToBeExtracted []string) (map[string]http.Cookie, map[string]string, *zkerrors.ZkError) {

	response, zkError := zkApiUtils.MakeRawApiCall("GET", nil, client, urlToBeCalled, &cookiesTobeAdded, nil)
	if zkError != nil {
		return nil, nil, zkError
	}
	var cookiesToBeReturned map[string]http.Cookie = map[string]http.Cookie{}
	var urlParamsToBeReturned map[string]string = map[string]string{}

	if urlParamsToBeExtracted != nil {
		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		val, err := url.QueryUnescape(string(responseData))
		for _, element := range urlParamsToBeExtracted {
			urlParamsToBeReturned[element] = zkApiUtils.ExtractUrlParam(element, val)
		}
	}

	cookiesToBeReturned = zkApiUtils.ExtractCookies(cookiesToBeExtracted, response.Cookies())
	return cookiesToBeReturned, urlParamsToBeReturned, nil
}

func (zkApiUtils zkApiUtils) ExtractCookies(cookiesToBeExtracted []string, cookies []*http.Cookie) map[string]http.Cookie {
	var cookiesToBeReturned map[string]http.Cookie = map[string]http.Cookie{}
	if cookiesToBeExtracted != nil {
		for _, element := range cookiesToBeExtracted {
			var foundCookie http.Cookie = zkApiUtils.ExtractCookie(element, cookies)
			cookiesToBeReturned[element] = foundCookie
		}
	}
	return cookiesToBeReturned
}

func (zkApiUtils zkApiUtils) ExtractCookie(name string, cookies []*http.Cookie) http.Cookie {
	var foundCookie http.Cookie
	for _, element := range cookies {
		if element == nil {
			continue
		}
		cookieName := element.Name
		// We are checking only against '*' as this is format that we defined
		// For example: The config yaml file can contain the following to use a regex:
		// 		  - cookiesExtract:
		//      		- csrf_token.*
		if strings.Contains(name, "*") {
			match, _ := regexp.MatchString(name, cookieName)
			if match {
				foundCookie = *element
				break
			}
		} else if cookieName == name {
			foundCookie = *element
			break
		}
	}

	return foundCookie
}

func (zkApiUtils zkApiUtils) ExtractUrlParam(name string, responseUnescaped string) string {
	var foundParam string
	result := strings.Split(responseUnescaped, name+"=")
	result = strings.Split(result[1], "\"")
	result = strings.Split(result[0], "&")
	foundParam = result[0]
	return foundParam
}

func (zkApiUtils zkApiUtils) ExtractUrlParamFromResponse(name string, response http.Response) string {
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	val, err := url.QueryUnescape(string(responseData))
	return zkApiUtils.ExtractUrlParam(name, val)
}
