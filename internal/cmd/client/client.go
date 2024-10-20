package client

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kolesnikovm/tcp-pow/internal/configs"
	wisdomclient "github.com/kolesnikovm/tcp-pow/pkg/client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "client",
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, err := cmd.InheritedFlags().GetString("config")
		if err != nil {
			return err
		}

		config, err := configs.NewConfig(configFile)
		if err != nil {
			return err
		}

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		ctx, cancel := context.WithCancel(context.Background())

		for i := 0; i < config.Concurrency; i++ {
			go func() {
				runClient(ctx, config.ServerAddress)
			}()
		}

		sig := <-sigc
		log.Info().Stringer("signal", sig).Msg("received os signal")

		cancel()

		return nil
	},
}

func runClient(ctx context.Context, addr string) {
	client, err := wisdomclient.NewWisdomClient(addr)
	if err != nil {
		log.Fatal().Err(err).Str("addr", addr).Msg("failed to create wisdom client")
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := client.GetQuote(ctx)
			if err != nil {
				log.Error().Err(err).Str("addr", addr).Msg("failed to get quote")
			}
		}
	}
}
