package test

import (
	"encoding/json"
	"testing"

	nt "github.com/Traders-Connect/esb-contract/golang/notification_manager"
	"github.com/Traders-Connect/notification-manager/worker"
	"github.com/Traders-Connect/utils"
)

func TestSendTelegramNotification(t *testing.T) {
	log, _ := utils.NewLogger("notification-server", "info")
	task_send_telegram := worker.TaskSendTelegramNotification{Worker: &worker.Worker{Db: GetTestDB(), Logger: log}}
	/*
		<p>Hi %FIRST_NAME%,</p><p>Your trade copy action has failed.</p><p>We attempted to copy from account <strong>%COPIER_MASTER%</strong> to account <strong>%COPIER_SLAVE%</strong> but faced the error shown below:</p><p><strong>Error</strong> - %COPIER_ERROR%</p><p>If this is unexpected, or you are unsure what the error means, please reach out to our support team.</p>
	*/
	data := map[string]string{
		nt.NotificationDataKeys_COPIER_MASTER.String(): "master123",
		nt.NotificationDataKeys_COPIER_SLAVE.String():  "slave123",
		nt.NotificationDataKeys_COPIER_ERROR.String():  "copy error",
	}
	dataBytes, _ := json.Marshal(&data)
	task_send_telegram.SendTelegramNotification("3", "", "TRADE_COPY_FAILURE", dataBytes)
}
