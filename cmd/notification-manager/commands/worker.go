package commands

import (
	"log"

	"github.com/Traders-Connect/notification-manager/contract"
	"github.com/Traders-Connect/utils"
	"github.com/spf13/cobra"
)

func init() {
	registerFlags(workerCmd)
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "worker",
}

func registerFlags(c *cobra.Command) {
	args := contract.GetWorkerArgs()
	c.Flags().StringVarP(&args.Name, "worker-name", "n", utils.LookupEnvOrString("NOTIFICATION_MANAGER_WORKER_NAME", ""), "Worker name ex: nt_master|nt-email")
	c.Flags().StringVarP(&args.WorkerType, "type", "t", utils.LookupEnvOrString("NOTIFICATION_MANAGER_WORKER_TYPE", ""), "Worker type ex: master or slave")

	// machinery
	c.Flags().StringVarP(&args.WorkerConfig.Broker, "broker-uri", "b", utils.LookupEnvOrString("NOTIFICATION_MANAGER_BROKER_URI", "amqp://myuser:mypass@rabbitmq:5672/"), "Broker URL ")
	c.Flags().StringVarP(&args.WorkerConfig.ResultBackend, "redis-backend", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_DEFAULT_BACKEND", "yourpassword@redis:6379"), "Redis backend URL")
	c.Flags().StringVarP(&args.WorkerConfig.DefaultQueue, "queue", "q", utils.LookupEnvOrString("NOTIFICATION_MANAGER_DEFAULT_QUEUE", "nt-master"), "Queue name")
	c.Flags().StringVarP(&args.WorkerConfig.AMQP.Exchange, "exchange", "e", utils.LookupEnvOrString("NOTIFICATION_MANAGER_EXCHANGE", "notification_exchange"), "Exchange name depending on worker types")
	c.Flags().StringVarP(&args.WorkerConfig.AMQP.ExchangeType, "exchange-type", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_EXCHANGE_TYPE", "direct"), "Exchange type")
	c.Flags().StringVarP(&args.WorkerConfig.AMQP.BindingKey, "bind-key", "k", utils.LookupEnvOrString("NOTIFICATION_MANAGER_BINDING_KEY", ""), "RabbtiMq binding key")

	// worker
	wc, err := utils.LookupEnvOrInt64("NOTIFICATION_MANAGER_WORKER_COUNT", 10)
	if err != nil {
		log.Fatal(err)
	}
	c.Flags().IntVarP(&args.Concurrency, "concurrency", "c", int(wc), "Concurrency value for a worker")

	// db
	c.Flags().StringVarP(&args.DbHost, "db-host", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_DB_HOST", "db:3306"), "DB host name")
	c.Flags().StringVarP(&args.DbUser, "db-user", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_DB_USER", "notificationmanager"), "DB user name")
	c.Flags().StringVarP(&args.DbPass, "db-pass", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_DB_PASSWORD", "password"), "DB password")
	c.Flags().StringVarP(&args.DbName, "db-name", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_DB_NAME", "notificationmanager"), "DB name")

	//mailgun creds
	c.Flags().StringVarP(&args.EmailApiKey, "email-api-key", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_EMAIL_API_KEY", ""), "Email api key")
	c.Flags().StringVarP(&args.EmailBaseUrl, "email-baseurl", "", utils.LookupEnvOrString("NOTIFICATION_MANAGER_EMAIL_BASE_URL", ""), "Email base key")
}
