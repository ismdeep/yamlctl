package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ismdeep/yamlctl/pkg/yamlx"
)

var setCmd = &cobra.Command{
	Use:   "set [file] [key] [value]",
	Short: "Set a YAML value",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		key := args[1]
		value := args[2]

		root, err := yamlx.Load(file)
		if err != nil {
			return err
		}

		target, err := yamlx.Set(root, key, value)
		if err != nil {
			return err
		}

		return yamlx.SaveScalar(file, target)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
}
