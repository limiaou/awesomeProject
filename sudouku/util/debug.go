package util

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type OneErrorStruct struct {
	Domain  string `json:"domain"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type ErrStruct struct {
	Errors  []OneErrorStruct `json:"errors"`
	Code    int              `json:"code"`
	Message string           `json:"message"`
}

type RequestLog struct {
	Time        time.Time
	ClientIP    string
	UserAgent   string
	Cookies     []*http.Cookie
	URL         *url.URL
	HandlerName string
	Method      string
	Header      http.Header
	RequestBody interface{}
	Error       ErrStruct
}

func APIErr(c *gin.Context, errStruct ErrStruct, requestBody interface{}) {
	c.JSON(errStruct.Code, gin.H{
		"error": errStruct,
	})
	errAPILog := RequestLog{
		Time:        time.Now(),
		ClientIP:    c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		Cookies:     c.Request.Cookies(),
		URL:         c.Request.URL,
		HandlerName: c.HandlerName(),
		Method:      c.Request.Method,
		Header:      c.Request.Header,
		RequestBody: requestBody,
		Error:       errStruct,
	}
	json, err := MarshalIndent(errAPILog)
	if err != nil {
		LogUnexpectedErr(err)
		return
	}
	LogAPIErr(c.HandlerName(), json)
	return
}

func JSONStatusOK(c *gin.Context, responseBody interface{}) {
	c.JSON(http.StatusOK, responseBody)
	return
}

func JPEGStatusOK(c *gin.Context, imgBytes []byte) {
	c.Data(http.StatusOK, "image/jpeg", imgBytes)
	return
}

func RequestInfo(c *gin.Context, errStruct ErrStruct, requestBody interface{}) {
	requestAPILog := RequestLog{
		Time:        time.Now(),
		ClientIP:    c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		Cookies:     c.Request.Cookies(),
		URL:         c.Request.URL,
		HandlerName: c.HandlerName(),
		Method:      c.Request.Method,
		Header:      c.Request.Header,
		RequestBody: requestBody,
		Error:       errStruct,
	}
	json, err := MarshalIndent(requestAPILog)
	if err != nil {
		LogUnexpectedErr(err)
		return
	}
	LogDebugMsg(json)
	return
}
