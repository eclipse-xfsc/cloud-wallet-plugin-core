package core

import (
	"context"
	"net/http"

	"github.com/cloudevents/sdk-go/v2/event"

	"github.com/Nerzal/gocloak/v13"
	cloudeventprovider "github.com/eclipse-xfsc/cloud-event-provider"
	"github.com/eclipse-xfsc/nats-message-library/common"
)

const UserKey = "user"

type EventType string

const (
	CreateKeyEventType EventType = "signer.createKey"
	SignTokenEventType EventType = "signer.signToken"
)

const SignerTopic = "signer-topic"

const (
	UserNotFound = "user not found"
)

type dataFetcher interface {
	GetUserInfo(ctx context.Context, accessToken string, realm string) (*gocloak.UserInfo, error)
}

type UserInfo struct {
	*gocloak.UserInfo
}

func (u *UserInfo) ID() string {
	sub := u.UserInfo.Sub
	return *sub
}

type Connector interface {
	getConnection(string, cloudeventprovider.ConnectionType) (*cloudeventprovider.CloudEventProviderClient, error)
}

type EventBus interface {
	Request(context.Context, string, event.Event) (*event.Event, error)
	Reply(context.Context, string, func(context.Context, event.Event) (*event.Event, error)) error
	Publish(context.Context, string, event.Event) error
	Subscribe(context.Context, string, func(event.Event)) error
}

type Message interface {
	CreateKey(string, string, string) error
	CreateToken(string, string, []byte) ([]byte, error)
}

type createKeyRequest struct {
	common.Request
	Namespace string `json:"namespace"`
	Group     string `json:"group"`
	Key       string `json:"key"`
	Type      string `json:"type"`
}

type createTokenRequest struct {
	common.Request
	Namespace string `json:"namespace"`
	Group     string `json:"group"`
	Key       string `json:"key"`
	Payload   []byte `json:"payload"`
	Header    []byte `json:"header"`
}

type createTokenReply struct {
	common.Reply
	Token []byte `json:"token"`
}

type Metadata struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Description string `json:"description"`
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type StandardQueryParams struct {
	Type        string `json:"type"`
	UserHint    string `json:"userHint"`
	Version     string `json:"version"`
	State       string `json:"state"`
	RedirectUri string `json:"redirect_uri"`
	Payload     any    `json:"payload"`
}

type DIDCommNotification struct {
	Account   string `json:"account"`
	RemoteDID string `json:"remoteDid"`
}
