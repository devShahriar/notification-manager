package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RichardKnop/machinery/v2/tasks"
	nm "github.com/Traders-Connect/esb-contract/golang/notification_manager"
	"github.com/Traders-Connect/notification-manager/contract"
)

func (n *NotificationService) IntAddUserConfig(ctx context.Context, payload *nm.UserMetaReq) (*nm.IntAddUserConfigReply, error) {
	//TODO Insert user config
	//Insert in user table if user not exist
	meta := &nm.InstallIntegrationReq{
		UserId:       payload.UserId,
		DefaultEmail: payload.DefaultEmail,
		FirstName:    payload.FirstName,
	}
	err := n.Db.SetUserConfig(ctx, meta)
	if err != nil {
		return nil, err
	}
	return &nm.IntAddUserConfigReply{}, nil
}

func (n *NotificationService) IntAddAccountConfig(ctx context.Context, payload *nm.UserMetaReq) (*nm.IntAddAccountConfigReply, error) {
	if err := n.Db.SetAccountConfig(ctx, payload); err != nil {
		return nil, err
	}
	return &nm.IntAddAccountConfigReply{}, nil
}

func (n *NotificationService) IntDeleteAccountConfig(ctx context.Context, payload *nm.IntDeleteAccountConfigReq) (*nm.IntDeleteAccountConfigReply, error) {

	ntPayload := &nm.NotificationReq{
		UserId:    payload.UserId,
		AccountId: payload.AccountId,
		EventType: nm.EventType_ACCOUNT_DELETED.String(),
		Data:      payload.Data,
	}
	err := n.Db.DeleteAccountConfig(ctx, payload.AccountId)
	if err != nil {
		n.Log.Errorw("error:", err)
		return nil, err
	}

	_, ntErr := n.IntSendNotification(ctx, ntPayload)
	n.Log.Errorw("Error", ntErr)
	return &nm.IntDeleteAccountConfigReply{}, nil
}

func (n *NotificationService) IntSendNotification(ctx context.Context, payload *nm.NotificationReq) (*nm.IntSendNotificationReply, error) {

	if payload.AccountId == "" || payload.EventType == "" {
		return nil, fmt.Errorf("AccountId or eventType is empty")
	}

	n.Logger.Infof("SendNotification req body %+v", payload)

	dataBytes, err := json.Marshal(payload.Data)
	if err != nil {
		n.Logger.Errorw("Error: while converting payload.Data into bytes", err)
		return nil, err
	}
	n.Logger.Info(contract.GetWorkerArgs().WorkerConfig.AMQP.BindingKey)
	taskSignature := &tasks.Signature{
		Name:       "task_route_notification",
		RoutingKey: contract.GetWorkerArgs().WorkerConfig.AMQP.BindingKey,
		Args: []tasks.Arg{
			{
				Name:  "userId",
				Type:  "string",
				Value: payload.UserId,
			},
			{
				Name:  "accId",
				Type:  "string",
				Value: payload.AccountId,
			},
			{
				Name:  "eventType",
				Type:  "string",
				Value: payload.EventType,
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

	if n.MachinaryServer == nil {
		return nil, fmt.Errorf("MachinaryServer is null")
	}

	_, err = n.MachinaryServer.SendTask(taskSignature)
	if err != nil {
		n.Logger.Info(err)
	}
	return &nm.IntSendNotificationReply{}, nil
}
