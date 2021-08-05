package awskeystore

import (
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/truewhitespace/key-rotation/rotation"
	"time"
)

//AWSAccessKey is a paired down interpretation of the AWS API suitable to be used as rotation.Key in a rotation.KeyStore
type AWSAccessKey struct {
	//ID is the AWS access key ID
	ID string
	//Secret is the secret key for the given ID.  This is only valid when the key has been created an will be null at
	//all other times.
	Secret *string
	//Internalized time the AWS API reports the key has been created.
	created time.Time
}

func (a *AWSAccessKey) Created() time.Time {
	return a.created
}

//MaybeSecret converts teh possible secret value into a humanized form.
func (a *AWSAccessKey) MaybeSecret() string {
	if a.Secret == nil {
		return "{unknown}"
	} else {
		return *a.Secret
	}
}

//internalizeKeyFromKey takes an AWS IAM key to create an AWSAccessKey.  Key which are not `Active` will be noted with
//an invalid creation time.
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

//internalizeKeyFromMetadata takes an AWS IAM key to create an AWSAccessKey
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
