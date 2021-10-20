package completion

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
func New() *cobra.Command {
	cmd := &cobra.Command{
		Hidden: true,
		Use:    "completion [bash|zsh|fish|powershell]",
		Short:  "Generate completion script",
		Long: `To load completions:
	
	Bash:
	
	  $ source <(scli completion bash)
	
	  # To load completions for each session, execute once:
	  # Linux:
	  $ scli completion bash > /etc/bash_completion.d/scli
	  # macOS:
	  $ scli completion bash > /usr/local/etc/bash_completion.d/scli
	
	Zsh:
	
	  # If shell completion is not already enabled in your environment,
	  # you will need to enable it.  You can execute the following once:
	
	  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
	
	  # To load completions for each session, execute once:
	  $ scli completion zsh > "${fpath[1]}/_scli"
	
	  # You will need to start a new shell for this setup to take effect.
	
	fish:
	
	  $ scli completion fish | source
	
	  # To load completions for each session, execute once:
	  $ scli completion fish > ~/.config/fish/completions/scli.fish
	
	PowerShell:
	
	  PS> scli completion powershell | Out-String | Invoke-Expression
	
	  # To load completions for every new session, run:
	  PS> scli completion powershell > scli.ps1
	  # and source this file from your PowerShell profile.
	`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		// hack to disable persistent required 'environment' flag
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.InheritedFlags()
			flags.SetAnnotation("environment", cobra.BashCompOneRequiredFlag, []string{"unknown"})
			return nil
		},
		// end hack
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
		},
	}

	return cmd
}
