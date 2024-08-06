package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Traders-Connect/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Traders-Connect/notification-manager/contract"
	"github.com/Traders-Connect/notification-manager/db"
	"github.com/Traders-Connect/notification-manager/worker"
)

var slaveCmd = &cobra.Command{
	Use:   "slave",
	Short: "slave",
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

		db, err := db.NewMysql(DBDsn, logger)
		if err != nil {
			logger.Fatal(err)
		}

		w := &worker.Worker{
			Name:         arg.Name,
			WorkerType:   arg.WorkerType,
			Concurrency:  arg.Concurrency,
			WorkerConfig: arg.WorkerConfig,
			Logger:       logger,
			Db:           db,
		}

		w.InitTaskFactory()
		w.InitMachineryWorker()

		//Ingesting slave meta in worker meta table so master worker will be able to discover slave
		w.Db.IngestWorkerMeta(arg)

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
	registerFlags(slaveCmd)
	workerCmd.AddCommand(slaveCmd)
}
