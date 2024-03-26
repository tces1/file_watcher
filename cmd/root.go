package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tces1/file_watcher/local"
	"github.com/tces1/file_watcher/pkg"
	"github.com/tces1/file_watcher/webdav"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "file_watcher",
	Short: "fileWatcher is directory content watcher.",
	Long:  "fileWatcher is a tool to watch file server and notify user when new content is added.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Read cfgFile:", cfgFile)
		config, err := pkg.ReadConfig(cfgFile)
		if err != nil {
			fmt.Println("Error reading config file:", err)
			return
		}
		for _, watch := range config.Watches {
			if watch.Type == "local" {
				go local.Watch(watch, config)
			} else {
				go webdav.Watch(watch, config)
			}
		}
		select {}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVarP(&cfgFile, "config", "f", "./file_watcher.yaml", "file wathcer config file")
}
