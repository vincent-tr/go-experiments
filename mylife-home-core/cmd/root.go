package cmd

import (
	"mylife-home-core/pkg/plugins"
	"os"
	"os/signal"
	"syscall"

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
		logger.WithError(errors.Errorf("failed")).Error("bam")

		testComponent()
		testBus()

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

	plugin := plugins.GetPlugin("logic-base.value-binary")

	logger.Infof("Metadata = '%s'", plugin.Metadata())

	comp, err := plugin.Instantiate("test", map[string]any{"initialValue": true})
	if err != nil {
		panic(err)
	}

	comp.SetOnStateChange(func(name string, value any) {
		logger.Infof("State '%s' changed to '%v'", name, value)
	})

	logger.Infof("State = '%v'", comp.GetState())

	logger.Infof("Execute")
	comp.Execute("setValue", false)

	logger.Infof("Execute no change")
	comp.Execute("setValue", false)

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	s := <-exit
	logger.Infof("Got signal %s", s)

	logger.Infof("Terminate")
	comp.Termainte()
}
