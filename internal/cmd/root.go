package cmd

import (
	"github.com/kolesnikovm/tcp-pow/internal/cmd/client"
	"github.com/kolesnikovm/tcp-pow/internal/cmd/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "tcp-pow",
}

func init() {
	rootCmd.AddCommand(
		server.Cmd,
		client.Cmd,
	)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to execute root cmd")
	}
}
