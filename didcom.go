package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DIDCommConnection struct {
	RemoteDid     string            `json:"remoteDid"`
	RoutingKey    string            `json:"routingKey"`
	Protocol      string            `json:"protocol"`
	Topic         string            `json:"topic"`
	EventType     string            `json:"eventType"`
	Properties    map[string]string `json:"properties"`
	RecipientDids []string          `json:"recipientDids"`
	Added         time.Time         `json:"added"`
	Group         string            `json:"group"`
}

const (
	connectionEndpoint = "/admin/connections"
	queryParamGroup    = "group"
	queryParamSearch   = "search"
)

func GetDidCommConnectionList(client *http.Client, accountID string, search string) ([]DIDCommConnection, error) {
	connectionURL := fmt.Sprintf("%s%s", libConfig.DIDComm.Url, connectionEndpoint)
	req, err := http.NewRequest(http.MethodGet, connectionURL, nil)
	if err != nil {
		return nil, err
	}
	req.URL.Query().Add(queryParamGroup, accountID)
	if search != "" {
		req.URL.Query().Add(queryParamSearch, search)
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var remoteControlList []DIDCommConnection
	err = json.Unmarshal(body, &remoteControlList)
	if err != nil {
		return nil, err
	}

	return remoteControlList, nil
}
