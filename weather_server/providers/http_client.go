package providers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPClient struct{}

var (
	httpClient = HTTPClient{}
)

func (c HTTPClient) get(url string) ([]byte, error) {

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		respErr := fmt.Errorf("Unexpected response: %s", response.Status)
		return nil, respErr
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
