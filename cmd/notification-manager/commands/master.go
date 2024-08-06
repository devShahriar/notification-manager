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
	"github.com/devshahriar/notification-manager/worker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "master",
	Run: func(cmd *cobra.Command, args []string) {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		arg := contract.GetWorkerArgs()

		//Initiating logger
		logger, err := utils.NewLogger("notification-manager", "info")
		if err != nil {
			log.Fatal(err)
		}

		//Initiating DB
		DBDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", arg.DbUser, arg.DbPass, arg.DbHost, arg.DbName)

		Db, err := db.NewMysql(DBDsn, logger)
		if err != nil {
			logger.Fatal(err)
		}

		w := &worker.Worker{
			Name:         arg.Name,
			WorkerType:   arg.WorkerType,
			Concurrency:  arg.Concurrency,
			WorkerConfig: arg.WorkerConfig,
			Logger:       logger,
			Db:           Db,
		}
		w.InitTaskFactory()
		w.InitMachineryWorker()

		//Registers slave workers
		w.InitWorkerPool()

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			sig := <-sigs
			log.Info("received signals", "signal", sig.String())
			cancel()
		}()

		w.Run(ctx)
	},
}

func init() {
	registerFlags(masterCmd)
	workerCmd.AddCommand(masterCmd)
}
