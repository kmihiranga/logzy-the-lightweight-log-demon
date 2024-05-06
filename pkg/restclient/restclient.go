package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"

	logger "logzy/pkg/log"
)

type RestClient struct {
	httpClient *http.Client
}

var (
	instance *RestClient
	once     sync.Once
	log      *zap.SugaredLogger = logger.GetLogger().Sugar()
)

// GetInstance return a singleton instance of RestClient with an optional timeout
func GetInstance(timeout *time.Duration) *RestClient {
	once.Do(func() {
		defaultTimeout := 30 * time.Second
		if timeout == nil {
			timeout = &defaultTimeout
		}
		instance = &RestClient{
			httpClient: &http.Client{
				Timeout: *timeout,
			},
		}
	})
	return instance
}

func (client *RestClient) Get(url string) ([]byte, error) {
	return client.doRequest(http.MethodGet, url, nil)
}

// POST request
func (client *RestClient) Post(url string, payload interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Errorf(fmt.Sprintf("Error marshalling post data. %v", err))
		return nil, err
	}
	return client.doRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
}

// doRequest handles the common logic for sending HTTP requests
func (client *RestClient) doRequest(method, uri string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		log.Errorf(fmt.Sprintf("Error requesting data to uri. %v", err))
		return nil, err
	}

	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		log.Errorf(fmt.Sprintf("Error processing request data. %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
