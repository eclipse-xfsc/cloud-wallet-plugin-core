package core

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/cloudevents/sdk-go/v2/event"
	cloudeventprovider "github.com/eclipse-xfsc/cloud-event-provider"
	messaging "github.com/eclipse-xfsc/nats-message-library"
	msgCommon "github.com/eclipse-xfsc/nats-message-library/common"
	"github.com/google/uuid"
)

const messageErr = "message error"

func PublishHistoryEvent(topic string, eventType messaging.RecordEventType, record messaging.HistoryRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return SdkError{general: "failed to publish history event", specific: err}
	}
	e, err := cloudeventprovider.NewEvent(libConfig.Name, string(eventType), data)
	bus := NewEventBus()
	return bus.Publish(context.Background(), topic, e)
}

func PublishDidCommNotification(topic string, eventType messaging.RecordEventType, userId string, remoteDid string) error {
	var notification = DIDCommNotification{
		Account:   userId,
		RemoteDID: remoteDid,
	}
	data, err := json.Marshal(notification)
	if err != nil {
		return SdkError{general: "failed to publish didcomm notification", specific: err}
	}
	e, err := cloudeventprovider.NewEvent(libConfig.Name, string(eventType), data)
	bus := NewEventBus()
	return bus.Publish(context.Background(), topic, e)
}

func NewMessage() (Message, error) {
	bus := NewEventBus()
	m := message{eventBus: bus}
	return &m, nil
}

func NewEventBus() EventBus {
	return &eventBus{}
}

type eventBus struct {
}

func (b *eventBus) getConnection(topic string, connectionType cloudeventprovider.ConnectionType) (*cloudeventprovider.CloudEventProviderClient, error) {
	client, err := cloudEventsConnection(topic, connectionType)
	if err != nil {
		return nil, SdkError{general: "error establishing cloudEvents connection", specific: err}
	}
	return client, err
}
func (b *eventBus) Request(ctx context.Context, topic string, e event.Event) (*event.Event, error) {
	conn, err := b.getConnection(topic, cloudeventprovider.ConnectionTypeReq)
	if err != nil {
		return nil, err
	}
	return conn.RequestCtx(ctx, e)
}

func (b *eventBus) Reply(ctx context.Context, topic string, handler func(context.Context, event.Event) (*event.Event, error)) error {
	conn, err := b.getConnection(topic, cloudeventprovider.ConnectionTypeRep)
	if err != nil {
		return err
	}
	return conn.ReplyCtx(ctx, handler)
}

func (b *eventBus) Publish(ctx context.Context, topic string, e event.Event) error {
	conn, err := b.getConnection(topic, cloudeventprovider.ConnectionTypePub)
	if err != nil {
		return err
	}
	return conn.PubCtx(ctx, e)
}

func (b *eventBus) Subscribe(ctx context.Context, topic string, handler func(event.Event)) error {
	conn, err := b.getConnection(topic, cloudeventprovider.ConnectionTypeSub)
	if err != nil {
		return err
	}
	return conn.SubCtx(ctx, handler)
}

type message struct {
	eventBus EventBus
}

func (m *message) CreateKey(keyId string, accountId string, keyType string) error {
	eventData := createKeyRequest{
		Request: msgCommon.Request{
			TenantId:  libConfig.Tenant,
			RequestId: requestId(),
		},
		Namespace: libConfig.Crypto.Namespace,
		Group:     accountId,
		Key:       keyId,
		Type:      keyType,
	}
	data, err := json.Marshal(eventData)
	if err != nil {
		return SdkError{general: messageErr, specific: err}
	}
	ev, err := cloudeventprovider.NewEvent(libConfig.Name, string(CreateKeyEventType), data)
	if err != nil {
		return SdkError{general: messageErr, specific: err}
	}
	_, err = m.eventBus.Request(context.Background(), SignerTopic, ev)
	if err != nil {
		return SdkError{general: messageErr, specific: err}
	}
	return nil
}

func (m *message) CreateToken(keyId string, accountId string, data []byte) ([]byte, error) {
	eventData := createTokenRequest{
		Request: msgCommon.Request{
			TenantId:  libConfig.Tenant,
			RequestId: requestId(),
		},
		Namespace: libConfig.Tenant,
		Group:     accountId,
		Key:       keyId,
		Payload:   data,
		Header:    make([]byte, 0),
	}
	data, err := json.Marshal(eventData)
	if err != nil {
		return nil, SdkError{general: messageErr, specific: err}
	}
	ev, err := cloudeventprovider.NewEvent(libConfig.Name, string(SignTokenEventType), data)
	if err != nil {
		return nil, SdkError{general: messageErr, specific: err}
	}
	res, err := m.eventBus.Request(context.Background(), SignerTopic, ev)
	var output createTokenReply
	err = res.DataAs(&output)
	if err != nil {
		return nil, SdkError{general: messageErr, specific: err}
	}
	return output.Token, nil
}

func requestId() string {
	return strings.Join([]string{libConfig.Name, uuid.New().String()}, "_")
}
