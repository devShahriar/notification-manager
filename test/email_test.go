package test

import (
	"testing"

	"github.com/Traders-Connect/utils"
	"github.com/devshahriar/notification-manager/worker"
)

func TestEmail(t *testing.T) {
	log, _ := utils.NewLogger("notification-server", "info")

	taskEmail := worker.TaskSendEmail{&worker.Worker{Db: GetTestDbConn(), Logger: log}}
	taskEmail.SendEmail("1001", "898737", "trade_failed", []byte{})
}
