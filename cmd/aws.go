package cmd

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/spf13/cobra"
	"github.com/truewhitespace/key-rotation/awskeystore"
	"github.com/truewhitespace/key-rotation/rotation"
)

func updateAWSUser(cmd *cobra.Command, args []string, flags *awsFlags, rotationConfig *rotationFlags) (err error) {
	ctx := cmd.Context()
	username := args[0]

	keystore, err := flags.buildKeyStore(username)
	if err != nil {
		return err
	}

	var rotator *rotation.KeyRotation
	if rotator, err = rotationConfig.build(); err != nil {
		return err
	}

	var plan *rotation.KeyRotationPlan
	if plan, err = rotator.Plan(ctx, keystore); err != nil {
		return err
	}
	var keys rotation.KeyList
	if keys, err = plan.Apply(ctx, keystore); err != nil {
		return err
	}

	out := cmd.OutOrStdout()
	if _, err := fmt.Fprintf(out, "Keys for %s\n", username); err != nil {
		return err
	}
	for i, k := range keys {
		awsKey := k.(*awskeystore.AWSAccessKey)
		if _, err := fmt.Fprintf(out, "%d: %s -- %#v\n", i, awsKey.ID, awsKey.MaybeSecret()); err != nil {
			return err
		}
	}
	return nil
}

type awsFlags struct {
	providerType string
}

func (flags *awsFlags) buildKeyStore(forUser string) (result rotation.KeyStore, err error) {
	var awsClient *iam.IAM
	if flags.providerType == "default" {
		sess := session.Must(session.NewSession())
		awsClient = iam.New(sess)
	} else if flags.providerType == "localstack" {
		awsClient, err = awskeystore.NewLocalstackProvider()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("bad aws provider type " + flags.providerType)
	}
	result = awskeystore.NewAWSUserKeyStore(forUser, awsClient)
	return
}

func awsCmd() *cobra.Command {
	flags := &awsFlags{}
	config := &rotationFlags{}
	cmd := &cobra.Command{
		Use:     "aws [user]",
		Short:   "Rotates the specified AWS user",
		PreRunE: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return updateAWSUser(cmd, args, flags, config)
		},
	}
	cmd.Flags().StringVarP(&flags.providerType, "aws-provider", "a", "default", "Must be either {default,moto}")
	config.attach(cmd.Flags())
	return cmd
}
