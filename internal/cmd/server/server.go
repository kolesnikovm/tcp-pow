package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kolesnikovm/tcp-pow/internal/configs"
	"github.com/kolesnikovm/tcp-pow/internal/controller/tcp"
	"github.com/kolesnikovm/tcp-pow/internal/pow"
	wisdomrepo "github.com/kolesnikovm/tcp-pow/internal/repository/wisdom"
	wisdomservice "github.com/kolesnikovm/tcp-pow/internal/service/wisdom"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "server",
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

		_, cancel := context.WithCancel(context.Background())

		metricsMux := http.NewServeMux()
		metricsMux.Handle("/metrics", promhttp.Handler())
		metricsSrv := http.Server{Addr: config.MetricsAddress, Handler: metricsMux}

		go func() {
			log.Info().Msgf("metrics server listening on %s", metricsSrv.Addr)
			if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("metrics server failure")
			}
		}()

		wisdomRepo, err := wisdomrepo.NewWisdomRepo(config.QuotesFile)
		if err != nil {
			return err
		}

		wisdomService := wisdomservice.NewWisdomService(wisdomRepo)

		powFactory := pow.NewPowShieldFactory()

		tcpServer := tcp.NewTCPServer(config, wisdomService, powFactory)

		go func() {
			if err := tcpServer.Run(); err != nil {
				log.Fatal().Err(err).Msg("tcp server failure")
			}
		}()

		sig := <-sigc
		log.Info().Stringer("signal", sig).Msg("received os signal")

		cancel()

		cleanup := func(ctx context.Context) {
			tcpServer.Stop(ctx)
			metricsSrv.Shutdown(ctx)
		}

		shutdown(cleanup)

		return nil
	},
}

func shutdown(cleaup func(context.Context)) {
	const shutdownTimeout = 10
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout*time.Second)

	go func() {
		cleaup(ctx)
		cancel()
	}()

	for {
		select {
		case <-time.Tick(1 * time.Second):
			deadline, _ := ctx.Deadline()
			log.Info().Msgf("cleaning up... force shutdown in %s", time.Until(deadline).Round(time.Second))
		case <-ctx.Done():
			if errors.Is(context.DeadlineExceeded, ctx.Err()) {
				log.Fatal().Msgf("failed to shutdown gracefully")
			}

			log.Info().Msg("server stopped")

			return
		}
	}
}
