package test

import (
	"testing"

	"github.com/Traders-Connect/notification-manager/worker"
	"github.com/Traders-Connect/utils"
)

func TestEmail(t *testing.T) {
	log, _ := utils.NewLogger("notification-server", "info")

	taskEmail := worker.TaskSendEmail{&worker.Worker{Db: GetTestDbConn(), Logger: log}}
	taskEmail.SendEmail("1001", "898737", "trade_failed", []byte{})
}
