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

			plugin := plugins.GetPlugin("logic-base.value-binary")

			fmt.Printf("Metadata = '%s'\n", plugin.Metadata())

			comp, err := plugin.Instantiate("test", map[string]any{"initialValue": true})
			if err != nil {
				panic(err)
			}

			comp.SetOnStateChange(func(name string, value any) {
				fmt.Printf("State '%s' changed to '%v'\n", name, value)
			})

			fmt.Printf("State = '%v'\n", comp.GetState())

			fmt.Printf("Execute\n")
			comp.Execute("setValue", false)

			fmt.Printf("Execute no change\n")
			comp.Execute("setValue", false)

			fmt.Printf("Terminate\n")
			comp.Termainte()
		},
	})
}
