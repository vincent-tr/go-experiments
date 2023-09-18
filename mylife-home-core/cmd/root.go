package cmd

import (
	"fmt"
	"mylife-home-core/pkg/plugins"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	_ "mylife-home-common/bus"    // tmp
	_ "mylife-home-common/config" // tmp
	"mylife-home-common/log"
	_ "mylife-home-core-plugins-logic-base"
)

var logger = log.CreateLogger("mylife:home:core:main")

var rootCmd = &cobra.Command{
	Use:   "mylife-home-core",
	Short: "mylife-home-core - Mylife Home Core",
	Run: func(_ *cobra.Command, _ []string) {
		log.Configure()

		testComponent()
		testBus()

		logger.WithError(errors.Errorf("failed")).Error("bam")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func testBus() {

}

func testComponent() {
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
}
