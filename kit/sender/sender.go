package sender

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

func NewMessage(phoneNumber, content string) *sns.PublishInput {
	return &sns.PublishInput{
		Message:     aws.String(content),
		PhoneNumber: aws.String(phoneNumber),
	}
}
