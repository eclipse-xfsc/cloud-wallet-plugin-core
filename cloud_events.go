package core

import (
	"fmt"

	cloudeventprovider "github.com/eclipse-xfsc/cloud-event-provider"
)

func cloudEventsConnection(topic string, connectionType cloudeventprovider.ConnectionType) (*cloudeventprovider.CloudEventProviderClient, error) {
	client, err := cloudeventprovider.New(cloudeventprovider.Config{
		Protocol: cloudeventprovider.ProtocolTypeNats,
		Settings: cloudeventprovider.NatsConfig{
			Url:        libConfig.Nats.Url,
			QueueGroup: libConfig.Nats.QueueGroup,
		},
	}, connectionType, topic)

	if err != nil {
		logger.Error(err, "error during processing message")
		return nil, err
	} else {
		logger.Info(fmt.Sprintf("cloudEvents can be received over topic: %s", topic))
	}
	return client, nil
}
