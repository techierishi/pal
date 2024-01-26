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
	"github.com/techierishi/pal/cmd/sync"
	"github.com/techierishi/pal/config"
	"github.com/techierishi/pal/ds"
	"github.com/techierishi/pal/util"
	"golang.design/x/clipboard"
)

var (
	configFile string
	version    = "dev"
)

type InitData struct {
	dir string
}

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
	RootCmd.AddCommand(sync.RootCmdSync)

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

func clipboardCheckHandler(io *InitData, next func(error)) {
	err := clipboard.Init()
	config.Flag.HasClipboard = true
	if err != nil {
		config.Flag.HasClipboard = false
	}
	next(nil)
}

func configLoadHandler(io *InitData, next func(error)) {
	if configFile == "" {
		configFile = filepath.Join(io.dir, "config.yaml")
	}

	if err := config.Conf.Load(configFile); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	next(nil)
}

func initRunCheckHandler(io *InitData, next func(error)) {
	err := util.CheckIfInitRan()
	if err != nil {
		os.Exit(1)
	}
	next(nil)
}

// initConfig reads in config file and ENV variables if set.
func initConfig(dir string) func() {
	return func() {
		initIO := InitData{dir: dir}
		chain := &ds.Chain[*InitData]{}
		chain.Use(clipboardCheckHandler)
		chain.Use(configLoadHandler)
		chain.Use(initRunCheckHandler)
		chain.Execute(&initIO)
	}

}
