package http

import (
	"crypto/tls"
	"io"
	"net/http"
)

//NO_CONTENT_TYPE -
const NO_CONTENT_TYPE string = ""

//HttpRequestEntity -
type HttpRequestEntity struct {
	Url           string
	Username      string
	Password      string
	ContentType   string
	Authorization string
}

//RequestFunc -
type RequestFunc func(HttpRequestEntity, string, io.Reader) (*http.Response, error)

//Request -
func Request(entity HttpRequestEntity, method string, body io.Reader) (response *http.Response, err error) {
	transport := NewRoundTripper()
	req, err := http.NewRequest(method, entity.Url, body)
	if err != nil {
		return
	}

	if entity.Authorization == "" {
		req.SetBasicAuth(entity.Username, entity.Password)

	} else {
		req.Header.Add("Authorization", entity.Authorization)
	}

	if entity.ContentType != NO_CONTENT_TYPE {
		req.Header.Add("Content-Type", entity.ContentType)
	}
	return transport.RoundTrip(req)
}

//RequestAdaptor -
type RequestAdaptor func() (*http.Response, error)

//HttpGateway -
type HttpGateway interface {
	Get(HttpRequestEntity) RequestAdaptor
	Post(HttpRequestEntity, io.Reader) RequestAdaptor
	Put(HttpRequestEntity, io.Reader) RequestAdaptor
}

//NewHttpGateway -
var NewHttpGateway = func() HttpGateway {
	return &DefaultHttpGateway{}
}

//DefaultHttpGateway -
type DefaultHttpGateway struct{}

//Get -
func (gateway *DefaultHttpGateway) Get(entity HttpRequestEntity) RequestAdaptor {
	return func() (*http.Response, error) {
		return Request(entity, "GET", nil)
	}
}

//Post -
func (gateway *DefaultHttpGateway) Post(entity HttpRequestEntity, body io.Reader) RequestAdaptor {
	return func() (*http.Response, error) {
		return Request(entity, "POST", body)
	}
}

//Put -
func (gateway *DefaultHttpGateway) Put(entity HttpRequestEntity, body io.Reader) RequestAdaptor {
	return func() (*http.Response, error) {
		return Request(entity, "PUT", body)
	}
}

//NewRoundTripper -
var NewRoundTripper = func() http.RoundTripper {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}
