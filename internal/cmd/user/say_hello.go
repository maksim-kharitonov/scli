package say

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewSayHelloCmd represents say hello command
func (e *Env) NewSayHelloCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "hello",
		Short:   "say hello",
		Long:    `say hello`,
		Example: `Request: say hello`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Hello!")
			return nil
		},
	}

	return cmd
}
