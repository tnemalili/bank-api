package messaging

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/sftx/bank-api/core/ports"
	log "github.com/sirupsen/logrus"
)

type messagingClient struct {
	snsClient *sns.Client
}

func (c *messagingClient) Publish(topic string, event any) error {
	// Serialize the event to JSON
	message, err := json.Marshal(event)
	if err != nil {
		return err
	}
	// Publish the message to the specified SNS topic
	pubOutput, err := c.snsClient.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String(string(message)),
		TopicArn: aws.String(topic),
	})
	if err != nil {
		log.Printf("failed to publish message to SNS topic %s: %v", topic, err)
		return err
	}
	log.Printf("message published to SNS topic %s with message ID %s", topic, *pubOutput.MessageId)
	return nil
}


func NewMessagingClient() *messagingClient {
	
	ctx := context.Background()
	
	region := os.Getenv("AWS_REGION")
	// Default to us-east-1 if AWS_REGION is not set
	if region == "" {
		region = "us-east-1"
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Errorf("Failed to load SDK configuration, %v", err)
	}

	// Create an Amazon SNS client
	snsClient := sns.NewFromConfig(cfg)

	return &messagingClient{snsClient: snsClient}
}


// Get feedback in case we are not implementing the interface correctly
var _ ports.IMessagingService = (*messagingClient)(nil)