package worker

import (
	"context"

	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/backends/redis"
	amqpBroker "github.com/RichardKnop/machinery/v2/brokers/amqp"
	"github.com/RichardKnop/machinery/v2/config"
	machineryConf "github.com/RichardKnop/machinery/v2/config"
	lock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/devshahriar/notification-manager/db"
	log "github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

var (
	PrefetchCount = 150
)

type Worker struct {
	Name            string
	WorkerType      string
	Event           string
	WorkerConfig    *machineryConf.Config
	MachineryWorker *machinery.Worker
	MachineryServer *machinery.Server
	Concurrency     int
	Db              db.DB
	Logger          *zap.SugaredLogger
}

func (w *Worker) InitMachineryWorker() {

	conf := w.WorkerConfig

	mc := &config.Config{
		Broker:       conf.Broker,
		DefaultQueue: conf.DefaultQueue,
		AMQP: &config.AMQPConfig{
			Exchange:      conf.AMQP.Exchange,
			ExchangeType:  conf.AMQP.ExchangeType,
			BindingKey:    conf.AMQP.BindingKey,
			PrefetchCount: PrefetchCount,
		},
		Redis: &config.RedisConfig{

			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	log.Info(conf.ResultBackend)
	resultBackend := conf.ResultBackend
	broker := amqpBroker.New(mc)
	backend := redis.NewGR(mc, []string{resultBackend}, 5)
	lock := lock.New()

	machinaryServer := machinery.NewServer(mc, broker, backend, lock)
	machinaryWorker := machinaryServer.NewWorker("notification_worker", w.Concurrency)

	w.MachineryServer = machinaryServer
	w.MachineryWorker = machinaryWorker

}

func (w *Worker) ResisterTask() {
	task := GetTaskByWorkerName(w.Name)
	log.Info(task)
	err := w.MachineryServer.RegisterTasks(task)
	if err != nil {
		log.Info("Error while registering task")
		log.Info(err)
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.ResisterTask()
	if err := w.MachineryWorker.Launch(); err != nil {
		log.Info("[*] Error while launching worker")
		log.Info(err)
	}
}
