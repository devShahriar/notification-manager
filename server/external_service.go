package server

import (
	"context"

	nm "github.com/Traders-Connect/esb-contract/golang/notification_manager"
)

func (n *NotificationService) GetIntegrationStatus(ctx context.Context, payload *nm.IntegrationStatusReq) (*nm.IntegrationStatusReply, error) {

	resp, err := n.Db.GetIntegrationStatus(ctx, payload)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (n *NotificationService) InstallIntegration(ctx context.Context, payload *nm.InstallIntegrationReq) (*nm.InstallIntegrationReply, error) {

	err := n.Db.InstallIntegration(ctx, payload)
	return &nm.InstallIntegrationReply{}, err
}

func (n *NotificationService) AddConfig(ctx context.Context, payload *nm.NotificationConfig) (*nm.AddConfigReply, error) {
	if err := n.Db.AddConfig(ctx, payload); err != nil {
		return nil, err
	}
	return &nm.AddConfigReply{}, nil
}

func (n *NotificationService) EditConfig(ctx context.Context, payload *nm.EditConfigReq) (*nm.EditConfigReply, error) {

	err := n.Db.EditConfig(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &nm.EditConfigReply{}, nil
}

func (n *NotificationService) GetConfig(ctx context.Context, payload *nm.ConfigReq) (*nm.UserConfigResp, error) {
	UserConfig, err := n.Db.GetConfig(ctx, payload)

	if err != nil {
		n.Log.Errorw("Error: while fetching UserConfig for userId: %v", payload.UserId)
		return nil, err
	}
	return UserConfig, nil
}

func (n *NotificationService) GetConfigDetails(ctx context.Context, payload *nm.ConfigDetailsReq) (*nm.ConfigDetailsReply, error) {

	reply, err := n.Db.GetConfigDetails(ctx, payload)
	if err != nil {
		return nil, err
	}
	return reply, nil

}

func (n *NotificationService) DeleteConfig(ctx context.Context, payload *nm.DeleteConfigReq) (*nm.DeleteConfigReply, error) {

	err := n.Db.DeleteConfig(ctx, payload)

	if err != nil {
		n.Log.Errorw("Error: while fetching UserConfig for userId: %v", payload.UserConfigId)
		return nil, err
	}
	return &nm.DeleteConfigReply{}, nil
}

func (n *NotificationService) EditAccountMeta(ctx context.Context, payload *nm.AccountMetaReq) (*nm.EditAccountMetaReply, error) {

	err := n.Db.EditAccountConfig(ctx, payload)

	if err != nil {
		n.Log.Errorw("Error: while fetching UserConfig for userId: %v", payload.AccountId)
		return nil, err
	}
	return &nm.EditAccountMetaReply{}, nil
}

func (n *NotificationService) EditConfigStatus(ctx context.Context, payload *nm.EditConfigStatusReq) (*nm.EditConfigStatusReply, error) {
	if err := n.Db.EditConfigStatus(ctx, payload); err != nil {
		return nil, err
	}
	return &nm.EditConfigStatusReply{}, nil
}

// Telegram-Bot
func (n *NotificationService) AddBot(ctx context.Context, payload *nm.AddBotReq) (*nm.AddBotReply, error) {

	if err := n.Db.AddBot(ctx, payload); err != nil {
		return nil, err
	}
	return &nm.AddBotReply{}, nil
}

func (n *NotificationService) EditBot(ctx context.Context, payload *nm.EditBotReq) (*nm.EditBotReply, error) {

	if err := n.Db.EditBot(ctx, payload); err != nil {
		return nil, err
	}
	return &nm.EditBotReply{}, nil
}

func (n *NotificationService) GetBots(ctx context.Context, payload *nm.GetBotsReq) (*nm.GetBotsReply, error) {

	reply, err := n.Db.GetBot(ctx, payload)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (n *NotificationService) EditBotStatus(ctx context.Context, payload *nm.EditBotStatusReq) (*nm.EditBotStatusReply, error) {

	err := n.Db.EditBotStatus(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &nm.EditBotStatusReply{}, nil
}

// Channel
func (n *NotificationService) AddChannel(ctx context.Context, payload *nm.AddChannelReq) (*nm.AddChannelReply, error) {

	err := n.Db.AddChannel(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &nm.AddChannelReply{}, nil
}

// Channel
func (n *NotificationService) EditChannel(ctx context.Context, payload *nm.EditChannelReq) (*nm.EditChannelReply, error) {

	err := n.Db.EditChannel(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &nm.EditChannelReply{}, nil
}

func (n *NotificationService) GetChannel(ctx context.Context, payload *nm.GetChannelReq) (*nm.GetChannelReply, error) {

	reply, err := n.Db.GetChannel(ctx, payload)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (n *NotificationService) EditChannelStatus(ctx context.Context, payload *nm.EditChannelStatusReq) (*nm.EditChannelStatusReply, error) {

	err := n.Db.EditChannelStatus(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &nm.EditChannelStatusReply{}, nil
}

func (n *NotificationService) DeleteChannel(ctx context.Context, payload *nm.DeleteChannelReq) (*nm.DeleteChannelReply, error) {

	err := n.Db.DeleteChannel(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &nm.DeleteChannelReply{}, nil
}

func (n *NotificationService) GetBotEventConfigs(ctx context.Context, payload *nm.GetBotEventConfigsReq) (*nm.GetBotEventConfigsReply, error) {

	reply, err := n.Db.GetBotEventConfigs(ctx, payload)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (n *NotificationService) GetBotEventDetails(ctx context.Context, payload *nm.GetBotEventDetailsReq) (*nm.GetBotEventDetailsReply, error) {
	reply, err := n.Db.GetBotEventDetails(ctx, payload)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (n *NotificationService) EditBotEventDetails(ctx context.Context, payload *nm.EditBotEventDetailsReq) (*nm.EditBotEventDetailsReply, error) {
	err := n.Db.EditBotEventDetails(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &nm.EditBotEventDetailsReply{}, nil
}

func (n *NotificationService) EditBotEventStatus(ctx context.Context, payload *nm.EditBotEventStatusReq) (*nm.EditBotEventStatusReply, error) {
	err := n.Db.EditBotEventStatus(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &nm.EditBotEventStatusReply{}, nil
}

func (n *NotificationService) DeleteBots(ctx context.Context, payload *nm.DeleteBotReq) (*nm.DeleteBotReply, error) {
	err := n.Db.DeleteBot(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &nm.DeleteBotReply{}, nil
}

func (n *NotificationService) UninstallIntegration(ctx context.Context, payload *nm.UninstallIntegrationReq) (*nm.UninstallIntegrationReply, error) {
	err := n.Db.UninstallIntegration(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &nm.UninstallIntegrationReply{}, nil
}
