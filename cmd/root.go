package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/techierishi/pal/cmd/alias"
	"github.com/techierishi/pal/cmd/clip"
	"github.com/techierishi/pal/cmd/cred"
	"github.com/techierishi/pal/cmd/hist"
	"github.com/techierishi/pal/cmd/snip"
	"github.com/techierishi/pal/cmd/svc"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/util"
)

var (
	configFile string
	version    = "dev"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           "pal",
	Short:         " Simple cli app which makes your cli interaction easier.",
	Long:          `pal -  Simple cli app which makes your cli interaction easier.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

// Execute adds all child commands to the root command sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Error executing pal cmd, %v", err)
		os.Exit(-1)
	}
}

func init() {

	dir := defaultDir()
	cobra.OnInitialize(initConfig(dir))
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(snip.RootCmdSnip)
	RootCmd.AddCommand(clip.RootCmdClip)
	RootCmd.AddCommand(cred.RootCmdCred)
	RootCmd.AddCommand(svc.RootCmdSvc)
	RootCmd.AddCommand(hist.RootCmdHist)
	RootCmd.AddCommand(alias.RootCmdAlias)

	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", fmt.Sprintf("config file (default is %s )", dir))
	RootCmd.PersistentFlags().BoolVarP(&config.Flag.Debug, "debug", "", false, "debug mode")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("pal version %s\n", version)
	},
}

func defaultDir() string {
	dir, err := config.GetDefaultConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	return dir
}

// initConfig reads in config file and ENV variables if set.
func initConfig(dir string) func() {
	return func() {
		if configFile == "" {
			configFile = filepath.Join(dir, "config.yaml")
		}

		if err := config.Conf.Load(configFile); err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(1)
		}

		err := util.CheckIfInitRan()
		if err != nil {
			os.Exit(1)
		}

	}

}
