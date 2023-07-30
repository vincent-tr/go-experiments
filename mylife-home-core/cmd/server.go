package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"mylife-home-core/pkg/"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "server",
		Short: "Starts the server",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Hello\n");
		},
	})
}
