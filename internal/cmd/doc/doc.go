package doc

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// completionCmd represents the completion command
func New() *cobra.Command {
	cmd := &cobra.Command{
		Hidden:                true,
		Use:                   "doc [man|md|rest]",
		Short:                 "Generate doc",
		Long:                  `Generate doc`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"man", "md", "rest"},
		Args:                  cobra.ExactValidArgs(1),
		// hack to disable persistent required 'environment' flag
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.InheritedFlags()
			flags.SetAnnotation("environment", cobra.BashCompOneRequiredFlag, []string{"unknown"})
			return nil
		},
		// end hack
	}

	pages := cmd.Flags().Bool("pages", false, "fix links for pages")
	cmd.Flags().MarkHidden("pages")

	path := cmd.Flags().String("path", "", "output docs to this folder")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		docsPath := "./docs"
		if *path != "" {
			docsPath = *path
		}

		if _, err := os.Stat(docsPath); os.IsNotExist(err) {
			err := os.Mkdir(docsPath, 0755)
			if err != nil {
				return err
			}
		}

		switch args[0] {
		case "man":
			header := &doc.GenManHeader{
				Title:   "SCLI",
				Section: "3",
			}
			err := doc.GenManTree(cmd.Root(), header, docsPath)
			if err != nil {
				return err
			}
		case "md":
			var identity func(s string) string
			if *pages {
				identity = func(s string) string { return s + ".html" }
			} else {
				identity = func(s string) string { return s }
			}
			emptyStr := func(s string) string { return "" }
			err := doc.GenMarkdownTreeCustom(cmd.Root(), docsPath, emptyStr, identity)
			if err != nil {
				return err
			}
		case "rest":
			err := doc.GenReSTTree(cmd.Root(), docsPath)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return cmd
}
