package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/truewhitespace/key-rotation/rotation"
	"time"
)

type rotationFlags struct {
	validFor     time.Duration
	expiresAfter time.Duration
}

func (r *rotationFlags) attach(f *pflag.FlagSet) {
	f.DurationVar(&r.validFor, "valid-for", 20*24*time.Hour, "how long a key should be considered valid and usable")
	f.DurationVar(&r.expiresAfter, "expires-after", 10*24*time.Hour, "grace period before deletion after validity")
}

func (r *rotationFlags) build() (*rotation.GracefulExpiration, error) {
	grace := r.validFor
	expiry := r.expiresAfter + r.validFor
	return rotation.NewGracefulExpiration(expiry, grace)
}

func NewRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-rotation",
		Short: "Plans or updates key changes for various systems",
	}
	cmd.AddCommand(awsCmd())
	return cmd
}
