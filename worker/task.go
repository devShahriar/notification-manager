package worker

var TaskFactory map[string]map[string]interface{}

// WORKER NAMES
const (
	EMAIL_WORKER    = "nt-email"
	MASTER          = "nt-master"
	TELEGRAM_WORKER = "nt-telegram"
	DISCORD_WORKER  = "nt-discord"
)

func (w *Worker) InitTaskFactory() {

	email := &TaskSendEmail{Worker: w}
	emailTask := map[string]interface{}{
		"task_send_email": email.SendEmail,
	}

	telegram := &TaskSendTelegramNotification{Worker: w}
	telegramTask := map[string]interface{}{
		"task_send_telegram": telegram.SendTelegramNotification,
	}

	discord := &TaskSendDiscordNotification{Worker: w}
	discordTask := map[string]interface{}{
		"task_send_discord": discord.SendDiscordNotification,
	}

	ntRouter := NotificationRouter{Worker: w}
	RouteNotificationTask := map[string]interface{}{
		"task_route_notification": ntRouter.RouteNotification,
	}

	TaskFactory = map[string]map[string]interface{}{
		EMAIL_WORKER:    emailTask,
		MASTER:          RouteNotificationTask,
		TELEGRAM_WORKER: telegramTask,
		DISCORD_WORKER:  discordTask,
	}
}

func GetTaskByWorkerName(workerName string) map[string]interface{} {
	return TaskFactory[workerName]
}
