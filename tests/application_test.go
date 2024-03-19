package tests

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/goccy/go-json"
	gotils "github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"

	"graphql-project/core"
)

func createApiRequest(token string) (*fasthttp.Request, error) {
	url := fasthttp.AcquireURI()
	err := url.Parse(nil, gotils.S2B(fmt.Sprintf("http://localhost:%d/graphql", Cfg.Port())))
	if err != nil {
		return nil, err
	}

	request := fasthttp.AcquireRequest()
	request.SetURI(url)

	fasthttp.ReleaseURI(url)

	request.Header.SetMethod(fasthttp.MethodPost)
	request.Header.Add("Authorization", "Bearer "+token)
	return request, err
}

func checkResponseError(response *fasthttp.Response) (m map[string]any, err error) {
	body := response.Body()
	if len(body) > 0 {
		if json.Unmarshal(body, &m) == nil {
			if errs, ok := m["errors"]; ok {
				s := "unknown error"
				if errs != nil {
					if b, e := json.MarshalIndent(errs, "", "  "); e == nil {
						s = gotils.B2S(b)
					}
				}
				err = errors.New(s)
			}
		}
	}
	if err == nil {
		if !(response.StatusCode() >= 200 && response.StatusCode() < 300) {
			if len(body) > 0 {
				err = errors.New(gotils.B2S(body))
			} else {
				err = errors.New(fasthttp.StatusMessage(response.StatusCode()))
			}
		}
	}
	return
}

func TestCreateUser(t *testing.T) {
	token, err := core.Anon(time.Hour, Cfg.JwtSecret())
	if err != nil {
		t.Errorf("create anon token %v", err)
		return
	}

	reqBody, err := loadRequestData("testdata/" + t.Name())
	if err != nil {
		t.Errorf("request %v", err)
		return
	}

	request, err := createApiRequest(token)
	if err != nil {
		t.Errorf("create request %v", err)
		return
	}
	defer fasthttp.ReleaseRequest(request)

	request.SetBody(reqBody)
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	client := &fasthttp.HostClient{
		Addr: fmt.Sprintf("localhost:%d", Cfg.Port()),
	}

	err = client.Do(request, response)

	if err != nil {
		t.Errorf("connection error %v", err)
		return
	}

	entity, err := checkResponseError(response)
	if err != nil {
		t.Errorf("request failed: %v", err)
		return
	}

	if err := Compare("testdata/"+t.Name(), entity); err != nil {
		t.Errorf("unexpected response:\n%v", err)
		return
	}
}
