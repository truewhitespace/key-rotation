package awskeystore

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/truewhitespace/key-rotation/rotation"
	"time"
)

type AWSAccessKey struct {
	ID      string
	Secret  *string
	created time.Time
}

func (a *AWSAccessKey) Created() time.Time {
	return a.created
}

func (a *AWSAccessKey) MaybeSecret() string {
	if a.Secret == nil {
		return "{unknown}"
	} else {
		return *a.Secret
	}
}

func internalizeKeyFromKey(k *iam.AccessKey) *AWSAccessKey {
	status := *k.Status
	var createdAt time.Time
	if status == "Active" {
		createdAt = *k.CreateDate
	} else {
		createdAt = rotation.InvalidTime()
	}

	return &AWSAccessKey{
		ID:      *k.AccessKeyId,
		Secret:  k.SecretAccessKey,
		created: createdAt,
	}
}

func internalizeKeyFromMetadata(k *iam.AccessKeyMetadata) *AWSAccessKey {
	status := *k.Status
	var createdAt time.Time
	if status == "Active" {
		createdAt = *k.CreateDate
	} else {
		createdAt = rotation.InvalidTime()
	}

	return &AWSAccessKey{
		ID:      *k.AccessKeyId,
		Secret:  nil,
		created: createdAt,
	}
}
