package cli

import (
	"github.com/spf13/cobra"
	"os"
)

var completionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long: `To load completions:

Bash:

# Dependencies:
# bash 4.1+
# bash_completion from Kubectl

$ source <(pocket completion bash)

# To load completions for each session, execute once:
Linux:
  $ pocket completion bash > /etc/bash_completion.d/pocket
MacOS:
  $ pocket completion bash > /usr/local/etc/bash_completion.d/pocket

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ pocket completion zsh > "${fpath[1]}/pocket"

# You will need to start a new shell for this setup to take effect.

Fish:

$ pocket completion fish | source

# To load completions for each session, execute once:
$ pocket completion fish > ~/.config/fish/completions/pocket.fish
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}
