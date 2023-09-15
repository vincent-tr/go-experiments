package cmd

import (
	"fmt"
	"mylife-home-core/pkg/plugins"

	"github.com/spf13/cobra"

	_ "mylife-home-core-plugins-logic-base"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Starts the server",
		Run: func(_ *cobra.Command, _ []string) {
			plugins.Build()

			for _, id := range plugins.Ids() {
				fmt.Printf("plugin: '%s'\n", id)
			}
		},
	})
}
