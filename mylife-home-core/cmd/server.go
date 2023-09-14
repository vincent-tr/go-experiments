package cmd

import (
	"fmt"
	"mylife-home-core-common/registry"

	"github.com/spf13/cobra"

	_ "mylife-home-core-plugins-logic-base"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Starts the server",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Hello\n")

			for index := 0; index < registry.NumPlugins(); index += 1 {
				plugin := registry.GetPlugin(index)
				fmt.Printf("plugin: '%s'\n", plugin.Metadata().Name())
			}
		},
	})
}
