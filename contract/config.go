package contract

import machineryConf "github.com/RichardKnop/machinery/v2/config"

var Args *WorkerArgs
var ServerArgs *ServiceArgs

type WorkerArgs struct {
	Name         string
	WorkerType   string
	WorkerConfig *machineryConf.Config
	Concurrency  int
	DbUser       string
	DbPass       string
	DbHost       string
	DbName       string
	EmailApiKey  string
	EmailBaseUrl string
}

type ServiceArgs struct {
	*WorkerArgs
	Addr             string
	AddrInt          string
	MetricsAddr      string
	TelegramBotToken string
}

func GetWorkerArgs() *WorkerArgs {
	if Args == nil {
		Args = &WorkerArgs{
			WorkerConfig: &machineryConf.Config{AMQP: &machineryConf.AMQPConfig{}},
		}
		return Args
	}
	return Args
}

func GetServerAgrs() *ServiceArgs {
	if ServerArgs == nil {
		ServerArgs = &ServiceArgs{}
	}
	return ServerArgs
}
