package awskeystore

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func NewLocalstackProvider() (client *iam.IAM, err error) {
	providers := []credentials.Provider{
		&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     "test",
			SecretAccessKey: "test",
			SessionToken:    "",
		}},
		&credentials.EnvProvider{},
	}

	var sess *session.Session
	awsCfg := &aws.Config{
		Region: aws.String("us-east-1"),
	}

	awsCfg.Credentials = credentials.NewChainCredentials(providers)
	sess, err = session.NewSession(awsCfg)
	if err != nil {
		return
	}

	client = iam.New(sess, &aws.Config{
		Endpoint: aws.String("http://localhost:4566"),
	})
	return
}
