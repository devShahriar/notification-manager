package worker

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RichardKnop/machinery/v2/tasks"
	nm "github.com/Traders-Connect/esb-contract/golang/notification_manager"
)

type NotificationRouter struct {
	*Worker
}

func (n *NotificationRouter) RouteNotification(userId, accId, eventType string, dataBytes []byte) error {

	n.Logger.Info("Receviced notification task")
	var data map[string]string
	_ = json.Unmarshal(dataBytes, &data)

	if userId != "" && eventType == nm.EventType_ACCOUNT_DELETED.String() {
		n.Logger.Info("Routing notification based on userId")
		err := n.SendUserSpecificNotification(userId, eventType, dataBytes)
		n.Logger.Info(err)
		return err
	}

	userConfigId, err := n.Worker.Db.GetUserConfig(ctx, accId)
	if err != nil {
		n.Worker.Logger.Infof("Failed to fetch UserConfig %v", userConfigId)
	}

	n.Worker.Logger.Infof("Retrieving UserConfig %v", userConfigId)
	//Check if notification for this event is disable dont send it to slave worker

	//GET enabled notificationTypes
	notificationsTypes, err := n.Worker.Db.GetEnabledNotificationTypes(ctx, userConfigId.UserConfigId, eventType)
	if err != nil {
		return err
	}

	sentNotificationType := map[string]bool{}
	for i := 0; i < len(notificationsTypes); i++ {

		if !n.ShouldSendNotification(userConfigId.UserId,
			userConfigId.AccountConfId,
			fmt.Sprintf("%d", notificationsTypes[i].ID),
			notificationsTypes[i].EventType,
			notificationsTypes[i].NotificationType) {
			n.Logger.Errorw("Not routing notification for the config below as the account or userConfig is disabled")
			continue
		}
		if _, ok := sentNotificationType[notificationsTypes[i].NotificationType]; ok {
			continue //NotificationConfig table might have duplicate notificationType . Ignore routing duplicate notification type
		} else {
			sentNotificationType[notificationsTypes[i].NotificationType] = true
		}

		worker := GetSlaveFromPool(notificationsTypes[i].NotificationType)

		taskSignature := GetRouteNotificationTask(
			GetTaskName(notificationsTypes[i].NotificationType),
			worker.WorkerConfig.AMQP.BindingKey,
			eventType,
			userConfigId.UserConfigId,
			accId,
			dataBytes)

		n.Logger.Info(notificationsTypes[i].NotificationType)
		n.Logger.Info(taskSignature)
		_, err = worker.MachineryServer.SendTask(taskSignature)
		if err != nil {
			n.Logger.Errorw("Error", err)
		} else {
			n.Logger.Infof("Send notification to %s", notificationsTypes[i].NotificationType)
		}

	}

	return nil
}

func (n *NotificationRouter) SendUserSpecificNotification(userId, eventType string, dataBytes []byte) error {

	status, err := n.Worker.Db.GetIntegrationStatus(ctx, &nm.IntegrationStatusReq{UserId: userId})
	if err != nil {
		n.Logger.Info("Failed to retrieved integration status")
		return err
	}

	userConfigId, _ := n.Db.GetUserConfigId(ctx, userId)
	n.Logger.Infof("Sending user specific notification for userConfigId:%v", userConfigId)

	notificationsTypes, err := n.Worker.Db.GetEnabledNotificationTypes(ctx, fmt.Sprintf("%d", *userConfigId), eventType)
	if err != nil {
		n.Logger.Errorw("Error:", err)
		return err
	}

	for i := 0; i < len(notificationsTypes); i++ {
		enabled := IsUserConfigEnabled(status, notificationsTypes[i].NotificationType)
		if !enabled {
			n.Logger.Infof("Notification integration not enabled for userId %v, notification type :%v", userId, notificationsTypes[i].NotificationType)
			continue
		}

		worker := GetSlaveFromPool(notificationsTypes[i].NotificationType)

		taskSignature := GetRouteNotificationTask(
			GetTaskName(notificationsTypes[i].NotificationType),
			worker.WorkerConfig.AMQP.BindingKey,
			eventType,
			fmt.Sprintf("%d", *userConfigId),
			"",
			dataBytes)

		n.Logger.Info(notificationsTypes[i].NotificationType)
		n.Logger.Info(taskSignature)
		_, err = worker.MachineryServer.SendTask(taskSignature)
		if err != nil {
			n.Logger.Errorw("Error", err)
		} else {
			n.Logger.Infof("Send notification to %s", notificationsTypes[i].NotificationType)
		}

	}
	return nil
}

func GetRouteNotificationTask(taskName, BindingKey, eventType, userConfig, accountId string, dataBytes []byte) *tasks.Signature {

	return &tasks.Signature{
		Name:       taskName,
		RoutingKey: BindingKey,
		Args: []tasks.Arg{
			{
				Name:  "userConfig",
				Type:  "string",
				Value: userConfig,
			},
			{
				Name:  "accountId",
				Type:  "string",
				Value: accountId,
			},
			{
				Name:  "eventType",
				Type:  "string",
				Value: eventType,
			},
			{
				Name:  "dataBytes",
				Type:  "[]byte",
				Value: dataBytes,
			},
		},
		RetryCount:   1,
		RetryTimeout: 100,
	}

}

func (n *NotificationRouter) ShouldSendNotification(userId string, accountConfigId string, notificationId string, eventType string, notificationType string) bool {

	disabled, err1 := n.Worker.Db.IsAccountNotificationDisabled(ctx, accountConfigId, notificationId)

	if err1 != nil {

		n.Logger.Errorw("Notification disabled for accountConfigId", accountConfigId, "eventType", eventType, "notificationType", notificationType)
		return false
	}

	status, err := n.Worker.Db.GetIntegrationStatus(ctx, &nm.IntegrationStatusReq{UserId: userId})
	if err != nil {
		return false
	}
	enabled := IsUserConfigEnabled(status, notificationType)

	return enabled && !disabled
}

func IsUserConfigEnabled(payload *nm.IntegrationStatusReply, notificationType string) bool {
	ntType := strings.ToUpper(notificationType)
	switch ntType {
	case nm.NotificationType_EMAIL.String():
		return payload.EmailEnabled
	case nm.NotificationType_DISCORD.String():
		return payload.DiscordEnabled
	case nm.NotificationType_SLACK.String():
		return payload.SlackEnabled
	case nm.NotificationType_TELEGRAM.String():
		return payload.TelegramEnabled
	case nm.NotificationType_WHATSAPP.String():
		return payload.WhatsappEnabled
	}
	return false
}
