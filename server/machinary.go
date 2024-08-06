package server

import (
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/backends/redis"
	amqpBroker "github.com/RichardKnop/machinery/v2/brokers/amqp"
	"github.com/RichardKnop/machinery/v2/config"
	lock "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/Traders-Connect/notification-manager/contract"
	"github.com/sirupsen/logrus"
)

func GetMachineryServer() *machinery.Server {
	conf := contract.GetWorkerArgs()
	logrus.Info(conf)
	mc := &config.Config{
		Broker:       conf.WorkerConfig.Broker,
		DefaultQueue: conf.WorkerConfig.DefaultQueue,
		AMQP: &config.AMQPConfig{
			Exchange:      conf.WorkerConfig.AMQP.Exchange,
			ExchangeType:  conf.WorkerConfig.AMQP.ExchangeType,
			BindingKey:    conf.WorkerConfig.AMQP.BindingKey,
			PrefetchCount: 1,
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
	resultBackend := conf.WorkerConfig.ResultBackend
	broker := amqpBroker.New(mc)
	backend := redis.NewGR(mc, []string{resultBackend}, 5)
	lock := lock.New()
	server := machinery.NewServer(mc, broker, backend, lock)

	return server

}
