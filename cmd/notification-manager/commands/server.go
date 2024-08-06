package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Traders-Connect/utils"
	"github.com/devshahriar/notification-manager/contract"
	"github.com/devshahriar/notification-manager/db"
	"github.com/devshahriar/notification-manager/server"
	"github.com/devshahriar/notification-manager/worker"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "server",
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		arg := contract.GetWorkerArgs()

		//Initiating logger
		logger, err := utils.NewLogger("notification-server", "info")
		if err != nil {
			logger.Fatal(err)
		}

		//Initiating DB
		DBDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", arg.DbUser, arg.DbPass, arg.DbHost, arg.DbName)

		Db, err := db.NewMysql(DBDsn, logger)
		if err != nil {
			logger.Fatal(err)
		}

		Db.IngestDefaultConfigTable()

		w := &worker.Worker{
			WorkerConfig: arg.WorkerConfig,
		}
		w.InitMachineryWorker()

		mServer := server.GetMachineryServer()

		service := server.NewNotificationService(Db, mServer, logger, contract.GetServerAgrs())
		if service.MachinaryServer == nil {
			logger.Info("mserver is nil")
		} else {
			logger.Info("mserver is okay")
		}

		ctx, cancel := context.WithCancel(context.Background())

		// handle signals
		go func() {
			sig := <-sigs
			logger.Infow("received signals", "signal", sig.String())
			cancel()
		}()
		service.Run(ctx)

	},
}

func init() {
	registerServiceFlags(serverCmd)
	rootCmd.AddCommand(serverCmd)
}

func registerServiceFlags(c *cobra.Command) {
	registerFlags(c)
	workerAgs := contract.GetWorkerArgs()
	serviceArgs := contract.GetServerAgrs()
	serviceArgs.WorkerArgs = workerAgs
	c.Flags().StringVarP(&serviceArgs.Addr, "addr", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_ADDR", "0.0.0.0:9030"), "Grpc service address")
	c.Flags().StringVarP(&serviceArgs.AddrInt, "addr-int", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_ADDR_INT", "0.0.0.0:9031"), "Grpc service port")
	c.Flags().StringVarP(&serviceArgs.MetricsAddr, "metrics-addr", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_METRICS_ADDR", "0.0.0.0:9035"), "Grpc service prometheus metrics address")

	c.Flags().StringVarP(&serviceArgs.TelegramBotToken, "telegram-bot-token", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_TELEGRAM_BOT_TOKEN", ""), "Telegram default bot token")
}
