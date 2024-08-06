package db

import (
	"context"

	nm "github.com/Traders-Connect/esb-contract/golang/notification_manager"
	"github.com/devshahriar/notification-manager/contract"
	"github.com/devshahriar/notification-manager/model"
)

type DB interface {
	GetUserConfig(context.Context, string) (contract.ConfigIds, error)
	GetEnabledNotificationTypes(context.Context, string, string) ([]model.NotificationConfig, error)
	GetEmailMeta(context.Context, string, string, string) (contract.EmailMeta, error)
	GetWorkerMeta() ([]contract.WorkerMeta, error)

	SetUserConfig(context.Context, *nm.InstallIntegrationReq) error
	SetAccountConfig(context.Context, *nm.UserMetaReq) error
	DeleteUserConfig(context.Context, string)
	IngestWorkerMeta(args *contract.WorkerArgs) error

	AddConfig(context.Context, *nm.NotificationConfig) error
	EditConfig(context.Context, *nm.EditConfigReq) error
	GetConfig(context.Context, *nm.ConfigReq) (*nm.UserConfigResp, error)
	DeleteConfig(context.Context, *nm.DeleteConfigReq) error
	EditAccountConfig(context.Context, *nm.AccountMetaReq) error
	AddDefaultNotificationConfig(context.Context, uint, string, []string, ...func(uint64)) error
	IngestDefaultConfigTable()
	GetIntegrationStatus(context.Context, *nm.IntegrationStatusReq) (*nm.IntegrationStatusReply, error)
	InstallIntegration(context.Context, *nm.InstallIntegrationReq) error
	IsAccountNotificationDisabled(context.Context, string, string) (bool, error)

	GetConfigDetails(context.Context, *nm.ConfigDetailsReq) (*nm.ConfigDetailsReply, error)
	PopulateBlockedAccountList(context.Context, uint64, *nm.ConfigDetailsReply) error
	PopulateEnabledAccountList(context.Context, uint64, uint64, *nm.ConfigDetailsReply) error
	CreateNotificationConfigIndex()
	EditConfigStatus(context.Context, *nm.EditConfigStatusReq) error

	DeleteAccountConfig(context.Context, string) error
	GetUserConfigId(context.Context, string) (*uint64, error)
	GetEmailMetaForUserOnly(context.Context, string, string) (contract.EmailMeta, error)

	//Telegram
	GetBotNotificationMeta(ctx context.Context, userConfigId, eventType, notificationType string) ([]model.BotNotificationMeta, error)

	//Bot
	AddBot(context.Context, *nm.AddBotReq) error
	EditBot(context.Context, *nm.EditBotReq) error
	GetBot(ctx context.Context, meta *nm.GetBotsReq) (*nm.GetBotsReply, error)
	EditBotStatus(ctx context.Context, meta *nm.EditBotStatusReq) error
	DeleteBot(ctx context.Context, meta *nm.DeleteBotReq) error

	AddChannel(context.Context, *nm.AddChannelReq) error
	EditChannel(context.Context, *nm.EditChannelReq) error
	GetChannel(ctx context.Context, meta *nm.GetChannelReq) (*nm.GetChannelReply, error)
	EditChannelStatus(ctx context.Context, meta *nm.EditChannelStatusReq) error
	DeleteChannel(ctx context.Context, meta *nm.DeleteChannelReq) error

	GetBotEventConfigs(ctx context.Context, req *nm.GetBotEventConfigsReq) (*nm.GetBotEventConfigsReply, error)
	GetBotEventDetails(ctx context.Context, req *nm.GetBotEventDetailsReq) (*nm.GetBotEventDetailsReply, error)
	EditBotEventDetails(ctx context.Context, req *nm.EditBotEventDetailsReq) error

	GetAccountListByUserConfig(ctx context.Context, userConfId uint64) ([]*nm.AccountMeta, error)
	BlockAccountAndChannel(ctx context.Context, userConfId, ntConfId uint64) error

	EditBotEventStatus(ctx context.Context, meta *nm.EditBotEventStatusReq) error

	PopulateBlockChannelList(ctx context.Context, userConfig, botConfigId, notificationConfigId uint64, data *nm.GetBotEventDetailsReply) error
	PopulateUnBlockChannelList(ctx context.Context, userConfig, botConfigId, notificationConfigId uint64, data *nm.GetBotEventDetailsReply) error
	UninstallIntegration(ctx context.Context, req *nm.UninstallIntegrationReq) error

	DumpLog(log model.Logs) error
	CheckValidUser(ctx context.Context, reqUserId string, configId uint64, tableName interface{}) bool
}
