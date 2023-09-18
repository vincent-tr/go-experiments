package cmd

import (
	"mylife-home-core/pkg/plugins"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"mylife-home-common/bus" // tmp
	_ "mylife-home-common/bus"
	"mylife-home-common/config"
	"mylife-home-common/log"
)

var logger = log.CreateLogger("mylife:home:core:main")

var configFile string
var logConsole bool

var rootCmd = &cobra.Command{
	Use:   "mylife-home-core",
	Short: "mylife-home-core - Mylife Home Core",
	Run: func(_ *cobra.Command, _ []string) {
		log.Init(logConsole)
		config.Init(configFile)
		plugins.Build()

		logger.WithError(errors.Errorf("failed")).Error("bam")

		testComponent()
		testBus()
		testExit()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $(PWD)/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&logConsole, "log-console", false, "Log to console")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func testBus() {
	bus.NewClient("instance-name")
}

func testComponent() {

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

	logger.Infof("Terminate")
	comp.Termainte()
}

func testExit() {

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	s := <-exit
	logger.Infof("Got signal %s", s)

}
