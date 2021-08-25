package pubsub

import (
	"context"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type cloudPubsubClientInterface interface {
	Close() error
	CreateTopic(ctx context.Context, topicId string) (*pubsub.Topic, error)
	GetTopic(ctx context.Context, topicId string) (*pubsub.Topic, error)
}

type cloudPubsubClient struct {
	cloudPubsubClient *pubsub.Client
}

func (c cloudPubsubClient) Close() error {
	return c.cloudPubsubClient.Close()
}

func (c cloudPubsubClient) CreateTopic(ctx context.Context, topicId string) (*pubsub.Topic, error) {
	return c.cloudPubsubClient.CreateTopic(ctx, topicId)
}

func (c cloudPubsubClient) GetTopic(ctx context.Context, topicId string) (*pubsub.Topic, error) {
	topic := c.cloudPubsubClient.Topic(topicId)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, grpc.Errorf(codes.NotFound, "Topic does not exist")
	}
	return topic, nil
}
