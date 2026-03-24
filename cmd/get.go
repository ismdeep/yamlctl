package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ismdeep/yamlctl/pkg/yamlx"
)

var getCmd = &cobra.Command{
	Use:   "get [file] [key]",
	Short: "Get a YAML value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := yamlx.Load(args[0])
		if err != nil {
			return err
		}

		val, err := yamlx.Get(root, args[1])
		if err != nil {
			return err
		}

		fmt.Println(val)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
