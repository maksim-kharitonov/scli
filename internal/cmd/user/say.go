package say

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SampleDB struct {
	URL     string
	Timeout time.Duration
}

// Env main app environment
type Env struct {
	SampleDB *SampleDB
}

// New returns a new command.
func New(env *Env) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "say",
		Short: "say smth",
		Long:  `say smth`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Root().PersistentPreRunE(cmd, args); err != nil {
				return err
			}

			environment := viper.GetString("environment")

			if err := viper.UnmarshalKey(environment, &env); err != nil {
				return err
			}

			if env.SampleDB == nil {
				return fmt.Errorf("section [sample db] not found in config")
			}

			return nil
		},
	}

	cmd.AddCommand(env.NewSayHelloCmd())

	return cmd
}
