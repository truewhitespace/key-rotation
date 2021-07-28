package awskeystore

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/truewhitespace/key-rotation/rotation"
)

func NewAWSUserKeyStore(username string, client *iam.IAM) *AWSUserKeyStore {
	return &AWSUserKeyStore{
		client:   client,
		username: username,
	}
}

type AWSUserKeyStore struct {
	client   *iam.IAM
	username string
}

func (a *AWSUserKeyStore) CreateKey(ctx context.Context) (rotation.Key, error) {
	response, err := a.client.CreateAccessKeyWithContext(ctx, &iam.CreateAccessKeyInput{UserName: &a.username})
	if err != nil {
		if awsErr, ok := err.(awserr.RequestFailure); ok {
			statusCode := awsErr.StatusCode()
			if statusCode == 404 {
				return nil, fmt.Errorf("no such user %q", a.username)
			}
		}
		return nil, err
	}
	return internalizeKeyFromKey(response.AccessKey), nil
}

func (a *AWSUserKeyStore) DeleteKey(ctx context.Context, key rotation.Key) error {
	actualKey := key.(*AWSAccessKey)
	_, err := a.client.DeleteAccessKeyWithContext(ctx, &iam.DeleteAccessKeyInput{
		AccessKeyId: &actualKey.ID,
		UserName:    &a.username,
	})
	return err
}

func (a *AWSUserKeyStore) ListKeys(ctx context.Context) (rotation.KeyList, error) {
	response, err := a.client.ListAccessKeysWithContext(ctx, &iam.ListAccessKeysInput{
		UserName: &a.username,
	})
	if err != nil {
		if awsErr, ok := err.(awserr.RequestFailure); ok {
			statusCode := awsErr.StatusCode()
			if statusCode == 404 {
				return make(rotation.KeyList, 0), nil
			}
		}
		return nil, err
	}

	out := make([]rotation.Key, len(response.AccessKeyMetadata))
	for i, m := range response.AccessKeyMetadata {
		key := internalizeKeyFromMetadata(m)
		out[i] = key
	}
	return out, nil
}

func (a *AWSUserKeyStore) MaximumKeys() int {
	return 2
}
