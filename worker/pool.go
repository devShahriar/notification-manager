package worker

import (
	"fmt"

	machineryConf "github.com/RichardKnop/machinery/v2/config"
	"github.com/devshahriar/notification-manager/contract"
	"github.com/sirupsen/logrus"
)

var WorkerPool map[string]*Worker

func (w *Worker) InitWorkerPool() {
	//TODO: GET list or meta of slave workers from DB
	//workerMeta := []contract.WorkerMeta{}
	WorkerPool = make(map[string]*Worker)
	// workerMeta = append(workerMeta, contract.WorkerMeta{
	// 	Name:             "nt-email",
	// 	NotificationType: "email",
	// 	WorkerType:       "slave",
	// 	Event:            "send_email",
	// 	Exchange:         "nt-email",
	// 	Queue:            "nt-email",
	// 	ExchangeType:     "direct",
	// 	BindingKey:       "nt-email"})
	workerMeta, err := w.Db.GetWorkerMeta()
	if err != nil {
		w.Logger.Info("Failed to load slave workers meta")
		panic("Failed to load slave workers meta")
	}
	for i := 0; i < len(workerMeta); i++ {

		workerInstance := &Worker{
			Name:       workerMeta[i].Name,
			WorkerType: workerMeta[i].WorkerType,
			WorkerConfig: &machineryConf.Config{
				Broker:       contract.GetWorkerArgs().WorkerConfig.Broker,
				DefaultQueue: workerMeta[i].Queue,
				AMQP: &machineryConf.AMQPConfig{
					Exchange:      workerMeta[i].Exchange,
					ExchangeType:  workerMeta[i].ExchangeType,
					BindingKey:    workerMeta[i].BindingKey,
					PrefetchCount: 150,
				},
				ResultBackend: contract.GetWorkerArgs().WorkerConfig.ResultBackend,
				Redis: &machineryConf.RedisConfig{

					MaxIdle:                3,
					IdleTimeout:            240,
					ReadTimeout:            15,
					WriteTimeout:           15,
					ConnectTimeout:         15,
					NormalTasksPollPeriod:  1000,
					DelayedTasksPollPeriod: 500,
				},
			},
		}

		workerInstance.InitMachineryWorker()
		logrus.Infof("[ * ] Registering slave worker: %v for notifcationType:%v", workerInstance.Name, workerMeta[i].NotificationType)
		WorkerPool[workerMeta[i].NotificationType] = workerInstance
	}
}

func GetSlaveFromPool(NotificationType string) *Worker {
	return WorkerPool[NotificationType]
}

func GetTaskName(notificationType string) string {
	return fmt.Sprintf("task_send_%s", notificationType)
}
