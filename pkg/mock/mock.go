package mock

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yesoreyeram/grafana-infinity-datasource/pkg/infinity"
	settingsSrv "github.com/yesoreyeram/grafana-infinity-datasource/pkg/settings"
)

type InfinityMocker struct {
	Body     string
	FileName string // filename (relative path of where it is being called)
}

func (rt *InfinityMocker) RoundTrip(req *http.Request) (*http.Response, error) {
	responseBody := "{}"
	if rt.Body != "" {
		responseBody = rt.Body
	}
	res := &http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Body:       ioutil.NopCloser(bytes.NewBufferString(responseBody)),
	}
	if rt.FileName != "" {
		b, err := ioutil.ReadFile(rt.FileName)
		if err != nil {
			return res, fmt.Errorf("error reading testdata file %s", rt.FileName)
		}
		reader := ioutil.NopCloser(bytes.NewReader(b))
		res.Body = reader
	}
	if res.Body != nil {
		return res, nil
	}
	return nil, errors.New("fake client not working as expected. If you got this error fix this method")
}

func New(body string) *infinity.Client {
	client, _ := infinity.NewClient(settingsSrv.InfinitySettings{})
	client.HttpClient.Transport = &InfinityMocker{Body: body}
	client.IsMock = true
	return client
}