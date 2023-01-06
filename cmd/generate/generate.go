package generate

import (
	"github.com/spf13/cobra"

	"easyhelm/assets"
	"easyhelm/internal/templates"
)

var Cmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Helm Chart",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(cmd)
	},
}

func run(cmd *cobra.Command) error {
	fs := assets.Assets
	err := templates.NewClient(&fs).GenerateTemplates()
	if err != nil {
		return err
	}

	return nil
}
