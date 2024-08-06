package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	nm "github.com/Traders-Connect/esb-contract/golang/notification_manager"
	"github.com/Traders-Connect/utils/mysql"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gLog "gorm.io/gorm/logger"

	"github.com/devshahriar/notification-manager/contract"
	"github.com/devshahriar/notification-manager/model"
)

type Mysql struct {
	*mysql.InstrumentedMysql
}

func NewMysql(dsn string, logger *zap.SugaredLogger) (DB, error) {
	im, err := mysql.NewInstrumentedMysql(dsn, "", "", logger, mysql.SetLogLevel(gLog.Info))
	if err != nil {
		return nil, err
	}

	im.AutoMigrate(
		model.UserConfig{},
		model.AccountConfig{},
		model.NotificationConfig{},
		model.WorkerMeta{},
		model.Logs{},

		model.DefaultConfigs{},
		model.AccountNotificationRules{},

		model.BotConfigs{},
		model.BotEventsRules{},
		model.ChannelConfig{},
		model.ChannelRules{},
	)

	return &Mysql{
		InstrumentedMysql: im,
	}, nil
}

func (m *Mysql) CreateNotificationConfigIndex() {
	err := m.DB.Exec("CREATE UNIQUE INDEX idx_eventtype_notification_type_user_config ON notification_configs (event_type, notification_type, user_config)").Error
	if err != nil {
		m.Log.Info("notification_config unique index already exist")
		return
	}
}

func (m *Mysql) InstallIntegration(ctx context.Context, req *nm.InstallIntegrationReq) error {
	fName := "InstallIntegrationStatus"
	start := time.Now()

	userConfigId, err := m.GetUserConfigId(ctx, req.UserId)
	if err != nil {
		m.Log.Errorw("Error while getting user config user doesn't exist UserId", req.UserId)
		return fmt.Errorf("user config doesn't exist userId: %v", req.UserId)
	}

	var setUpdateErr error

	if userConfigId != nil && err == nil {
		m.Log.Info("User exist updating user UserConfig")
		setUpdateErr = m.UpdateUserConfig(ctx, req, userConfigId)
	}

	if setUpdateErr == nil && req.NotificationType != nm.NotificationType_EMAIL {
		err := m.AddDefaultBot(ctx, req)
		if err != nil {
			return err
		}
	}

	m.LogError(fName,
		setUpdateErr != nil,
		fmt.Sprintf("Error: while installing notification integration for UserId %v | Sent Default status", req.UserId),
		fmt.Sprintf("Success: Installed notification integration for UserId %v", req.UserId),
		start)
	return setUpdateErr
}

func (m *Mysql) AddDefaultBot(ctx context.Context, req *nm.InstallIntegrationReq) error {
	fName := "AddDefaultBot"
	start := time.Now()

	userConfig, err := m.GetUserConfigId(ctx, req.UserId)
	if err != nil {
		m.Log.Errorw("Error while feting userConfig id for userId:", req.UserId)
		return err
	}

	botConfig := model.BotConfigs{BotName: "Traders connect"}
	botConfig.BotToken = contract.GetServerAgrs().TelegramBotToken
	botConfig.UserConfig = *userConfig
	botConfig.Enabled = true
	botConfig.NotificationType = contract.TELEGRAM

	result := m.DB.Model(&model.BotConfigs{}).FirstOrCreate(&botConfig)

	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error: while adding default %v | Sent Default status", req.UserId),
		fmt.Sprintf("Success: Added default tradersConnect Bot %v", req.UserId),
		start)

	return result.Error
}

func (m *Mysql) SetUserConfig(ctx context.Context, UserMeta *nm.InstallIntegrationReq) error {

	fName := "SetUserConfig"
	start := time.Now()

	user := &model.UserConfig{
		UserId:       UserMeta.UserId,
		DefaultEmail: UserMeta.DefaultEmail,
		FirstName:    UserMeta.FirstName,
	}

	oldConfig := model.UserConfig{}
	result := m.DB.WithContext(ctx).WithContext(ctx).FirstOrCreate(&oldConfig, user)

	if result.Error == nil && result.RowsAffected > 0 {
		m.Log.Infof("New userConfig created userId: %v", UserMeta.UserId)
	} else if result.Error == nil && result.RowsAffected == 0 {
		m.Log.Infof("User Config already exist userId: %v | userConfigId:%v", UserMeta.UserId, oldConfig.ID)
		return fmt.Errorf("User config already exist")
	}

	m.LogError(fName,
		result.Error != nil,
		fmt.Sprintf("Error: While creating new userConfig:%v", result.Error),
		"Success: User Config created successfully",
		start)

	return nil

}

func (m *Mysql) UpdateUserConfig(ctx context.Context, UserMeta *nm.InstallIntegrationReq, userConfigId *uint64) error {
	fName := "UpdateUserConfig"
	start := time.Now()

	column := GetColumnName(UserMeta.NotificationType.String())
	userConfig := &model.UserConfig{}
	PopulateUserConfigEnabledNotifications(userConfig, UserMeta.NotificationType.String())
	result := m.DB.WithContext(ctx).Model(&model.UserConfig{}).Where("user_id = ?", UserMeta.UserId).Select(column).Updates(userConfig)

	if result.Error != nil {
		m.Log.Errorw("Error while updating user config for Install Integration call")
		return fmt.Errorf("UserConfig update failed")
	}
	defaultEvents := contract.GetDefaultEventList()

	if UserMeta.NotificationType == nm.NotificationType_EMAIL {
		err := m.AddDefaultNotificationConfig(ctx, uint(*userConfigId), UserMeta.NotificationType.String(), defaultEvents)
		if err != nil {
			m.Log.Errorw("Error while ingesting default config for Install Integration call", "error:", err)
			return fmt.Errorf("error while ingesting default config for Install integration call", err)
		}
	}

	m.LogError(fName,
		result.Error != nil,
		fmt.Sprintf("Error: While Updating userConfig:%v or ingesting default configs", result.Error),
		"Success: User Config updated successfully",
		start)
	return result.Error
}

func (m *Mysql) GetIntegrationStatus(ctx context.Context, req *nm.IntegrationStatusReq) (*nm.IntegrationStatusReply, error) {

	fName := "GetIntegrationStatus"
	start := time.Now()

	integrationStatus := &model.UserConfig{}
	err := m.DB.WithContext(ctx).Model(&model.UserConfig{}).Select("user_id", "default_email", "email_enabled", "telegram_enabled", "discord_enabled", "slack_enabled", "whats_app_enabled").
		Where("user_id = ?", req.UserId).First(integrationStatus).Error

	resp := &nm.IntegrationStatusReply{
		UserId:          req.UserId,
		DefaultEmail:    integrationStatus.DefaultEmail,
		EmailEnabled:    integrationStatus.EmailEnabled,
		TelegramEnabled: integrationStatus.TelegramEnabled,
		DiscordEnabled:  integrationStatus.DiscordEnabled,
		SlackEnabled:    integrationStatus.SlackEnabled,
		WhatsappEnabled: integrationStatus.WhatsAppEnabled,
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: while retrieving integration status for UserId %v | Sent Default status", req.UserId),
		fmt.Sprintf("Success: Retrieved integration status for UserId %v", req.UserId),
		start)

	return resp, err
}

// Master Retrieves UserConfig Id To route notification
func (m *Mysql) GetUserConfig(ctx context.Context, accId string) (contract.ConfigIds, error) {

	fName := "GetUserConfig"
	start := time.Now()
	var configId contract.ConfigIds

	err := m.DB.WithContext(ctx).Table("account_configs a").Select("a.id, a.config_id, u.user_id").
		Joins("JOIN user_configs u on a.config_id=u.id").
		Where("a.account_id = ?", accId).Scan(&configId).Error

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: while retrieving UserConfigId for AccountId %v", accId),
		fmt.Sprintf("Success: Retrieved UserConfig for AccountId %v", accId),
		start)

	return configId, err
}

func (m *Mysql) GetEnabledNotificationTypes(ctx context.Context, userConfigId string, eventType string) ([]model.NotificationConfig, error) {

	fName := "GetEnabledNotificationTypes"
	start := time.Now()

	var notificationType []model.NotificationConfig

	err := m.DB.WithContext(ctx).Model(&model.NotificationConfig{}).
		Select("id", "event_type", "notification_type").
		Where("user_config = ? AND event_type = ? AND enabled = ?", userConfigId, eventType, true).
		Scan(&notificationType).Error

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: Retrieving enabled notification type for userConfig:%v and EventType:%v err:%+v", userConfigId, eventType, err),
		fmt.Sprintf("Success: Retrieved enabled notification types for userConfig:%v", userConfigId),
		start)

	return notificationType, nil
}

func (m *Mysql) IsAccountNotificationDisabled(ctx context.Context, accountConfigId, notificationConfigId string) (bool, error) {
	fName := "IsAccountNotificationDisabled"
	start := time.Now()

	// Execute the query using GORM
	var exists bool
	err := m.DB.WithContext(ctx).Raw("select EXISTS(select 1 from account_notification_rules where account_config_id= ? and notification_config_id = ? and disabled = ?) as result", accountConfigId, notificationConfigId, true).
		Scan(&exists).Error

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: Retrieving notification rule for accountConfig:%v ntConfId %v err:%+v", accountConfigId, err),
		fmt.Sprintf("Success: Retrieved notification rule for accountConfig:%v ntConfId", accountConfigId, notificationConfigId),
		start)
	return exists, err
}

func (m *Mysql) LogError(fName string, condition bool, errorMsg string, successMsg string, startedAt time.Time) {
	result := mysql.DbOpSuccess
	if condition {
		m.Log.Errorw(errorMsg)
		result = mysql.DbOpError
	} else {
		m.Log.Infow(successMsg)
	}

	m.MetricCount.WithLabelValues(fName, result).Inc()
	m.MetricDuration.WithLabelValues(fName, result).Observe(float64(time.Since(startedAt)))
}

func (m *Mysql) GetEmailMeta(ctx context.Context, userConfig, accountId, eventType string) (contract.EmailMeta, error) {
	fName := "GetEmailMeta"
	start := time.Now()

	var emailMeta contract.EmailMeta

	err := m.DB.WithContext(ctx).Table("user_configs a").
		Select("a.first_name, a.default_email, c.email, b.message_template, b.subject").
		Joins("JOIN notification_configs b ON a.id = b.user_config").
		Joins("JOIN account_configs c ON b.user_config = c.config_id").
		Where("a.id = ? AND b.notification_type = ? AND b.event_type = ? AND c.account_id = ?", userConfig, contract.EMAIL, eventType, accountId).
		Scan(&emailMeta).Error

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: Getting email meta for userConfig:%v error:%+v", userConfig, err),
		fmt.Sprintf("Success: Email meta fetched for userConfig:%v", userConfig),
		start)

	return emailMeta, nil
}

func (m *Mysql) GetEmailMetaForUserOnly(ctx context.Context, userConfig, eventType string) (contract.EmailMeta, error) {
	fName := "GetEmailMetaForUserOnly"
	start := time.Now()

	var emailMeta contract.EmailMeta

	err := m.DB.WithContext(ctx).Table("user_configs a").
		Select("a.first_name, a.default_email, b.message_template, b.subject").
		Joins("JOIN notification_configs b ON a.id = b.user_config").
		Where("a.id = ? AND b.notification_type = ? AND b.event_type = ?", userConfig, contract.EMAIL, eventType).
		Scan(&emailMeta).Error

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: Getting email meta for userConfig:%v error:%+v", userConfig, err),
		fmt.Sprintf("Success: Email meta fetched for userConfig:%v", userConfig),
		start)

	return emailMeta, nil
}

func GetColumnName(ntType string) string {
	switch ntType {
	case nm.NotificationType_EMAIL.String():
		return "email_enabled"
	case nm.NotificationType_TELEGRAM.String():

		return "telegram_enabled"
	case nm.NotificationType_DISCORD.String():

		return "discord_enabled"
	case nm.NotificationType_SLACK.String():
		return "slack_enabled"
	case nm.NotificationType_WHATSAPP.String():
		return "whats_app_enabled"
	}
	return ""
}

func PopulateUserConfigEnabledNotifications(userConfig *model.UserConfig, ntType string) {
	switch ntType {
	case nm.NotificationType_EMAIL.String():
		userConfig.EmailEnabled = true
		return
	case nm.NotificationType_TELEGRAM.String():
		userConfig.TelegramEnabled = true
		return
	case nm.NotificationType_DISCORD.String():
		userConfig.DiscordEnabled = true
		return
	case nm.NotificationType_SLACK.String():
		userConfig.SlackEnabled = true
		return
	case nm.NotificationType_WHATSAPP.String():
		userConfig.WhatsAppEnabled = true
		return
	}
}

func (m *Mysql) DeleteUserConfig(ctx context.Context, userId string) {

	err := m.DB.WithContext(ctx).Where("user_id = ?", userId).Delete(&model.UserConfig{}).Error

	if err != nil {
		m.Log.Errorw("Error while deleting userConfig for UserId", userId)
	}
}

func (m *Mysql) DeleteAccountConfig(ctx context.Context, accId string) error {

	err := m.DB.WithContext(ctx).Where("account_id = ?", accId).Delete(&model.AccountConfig{}).Error

	if err != nil {
		m.Log.Errorw("Error while deleting userConfig for UserId", accId)
		return err
	}

	return nil
}

func (m *Mysql) SetAccountConfig(ctx context.Context, userMeta *nm.UserMetaReq) error {

	fName := "SetAccountConfig"
	start := time.Now()

	userConfigId, err := m.GetUserConfigId(ctx, userMeta.UserId)
	if err != nil {
		m.Log.Errorw("Error: occurred while fetching userConfigId for AccountId %v Or userId doesn't exist", userMeta.AccountMeta.AccountId)
		return err
	}
	fmt.Println("Creating account config for accountId", userMeta.AccountMeta.AccountId)
	accountConf := model.AccountConfig{
		ConfigId:    *userConfigId,
		AccountId:   userMeta.AccountMeta.AccountId,
		AccountName: userMeta.AccountMeta.Nickname,
	}
	if userMeta.AccountMeta.Email != "" {
		accountConf.Email = &userMeta.AccountMeta.Email
	}
	oldAccountConfig := model.AccountConfig{}

	result := m.DB.WithContext(ctx).FirstOrCreate(&oldAccountConfig, &accountConf)

	if result.Error == nil && result.RowsAffected > 0 {
		m.Log.Infof("New account Config created userId: %v", userMeta.UserId)
	} else {
		m.Log.Infof("Account Config already exist userId: %v | userConfigId:%v", userMeta.UserId, oldAccountConfig.ID)
	}

	m.LogError(fName,
		result.Error != nil,
		fmt.Sprintf("Error: While creating new accountConfig:%v", result.Error),
		fmt.Sprintf("Success: Account Config created successfully accountId %v", userMeta.AccountMeta.AccountId),
		start)

	return err
}

func (m *Mysql) GetUserConfigId(ctx context.Context, userId string) (*uint64, error) {

	var userConfigId uint64
	err := m.DB.WithContext(ctx).Model(&model.UserConfig{}).Where("user_id = ?", userId).Select("ID").First(&userConfigId).Error

	if err != nil {
		m.Log.Errorw("Error getting UserConfigId")
		return nil, err
	}
	return &userConfigId, nil
}

func (m *Mysql) GetWorkerMeta() ([]contract.WorkerMeta, error) {

	fName := "GetWorkerMeta"
	start := time.Now()
	var workerMeta []contract.WorkerMeta

	err := m.DB.Model(&model.WorkerMeta{}).
		Select("name", "worker_type", "exchange", "queue", "exchange_type", "binding_key", "notification_type").
		Scan(&workerMeta).Error

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: While getting worker meta error:%v", err),
		fmt.Sprintf("Success: Worker meta fetched"),
		start)

	return workerMeta, nil
}

func (m *Mysql) IngestWorkerMeta(args *contract.WorkerArgs) error {
	fName := "IngestWorkerMeta"
	start := time.Now()

	existingWorkerMeta := &model.WorkerMeta{}

	workerMeta := &model.WorkerMeta{
		Name:             args.Name,
		WorkerType:       args.WorkerType,
		Exchange:         args.WorkerConfig.AMQP.Exchange,
		Queue:            args.WorkerConfig.DefaultQueue,
		ExchangeType:     args.WorkerConfig.AMQP.ExchangeType,
		BindingKey:       args.WorkerConfig.AMQP.BindingKey,
		NotificationType: args.WorkerType,
	}

	result := m.DB.FirstOrCreate(existingWorkerMeta, workerMeta)

	if result.Error == nil && result.RowsAffected > 0 {
		m.Log.Infof("New worker meta created name: %v", workerMeta.Name)
	} else {
		m.Log.Infof("Worker meta already exist: %v", workerMeta.Name)
	}

	m.LogError(fName,
		result.Error != nil,
		fmt.Sprintf("Error: While Ingesting worker meta error:%v", result.Error),
		fmt.Sprintf("Success: Worker meta ingested"),
		start)
	return result.Error
}

func (m *Mysql) AddConfig(ctx context.Context, notificationConfig *nm.NotificationConfig) error {

	fName := "AddConfig"
	start := time.Now()

	userConfig, err := m.GetUserConfigId(ctx, notificationConfig.UserId)

	if err != nil {
		m.Log.Infof("Error: Failed to retrieved userConfigId for userId:%v", notificationConfig.UserId)
		return err
	}

	existingConf := model.NotificationConfig{}
	ntConfig := &model.NotificationConfig{
		EventType:        notificationConfig.EventType,
		NotificationType: notificationConfig.NotificationType,
		Enabled:          true,
		MessageTemplate:  notificationConfig.MessageTemplate,
		Subject:          notificationConfig.Subject,
		UserConfig:       *userConfig,
	}

	result := m.DB.WithContext(ctx).FirstOrCreate(&existingConf, ntConfig)

	if result.Error == nil && result.RowsAffected > 0 {
		m.Log.Infof("New notification config created for userId: %v | userConfigId: %v", notificationConfig.UserId, notificationConfig.UserConfig)
	} else {
		m.Log.Infof("Notification Config already exits for userId %v | eventType: %v | notification_type: %v",
			notificationConfig.UserId, notificationConfig.EventType, notificationConfig.NotificationType)
	}

	m.LogError(fName,
		result.Error != nil,
		fmt.Sprintf("Error: While Ingesting notification Config error:%v", result.Error),
		fmt.Sprintf("Success: Notification config ingested"),
		start)

	return result.Error
}

func (m *Mysql) EditConfig(ctx context.Context, req *nm.EditConfigReq) error {

	fName := "EditConfig"
	start := time.Now()
	notificationConfig := req.NotificationConfig
	result := m.DB.WithContext(ctx).Debug().Model(&model.NotificationConfig{}).
		Where("id = ? AND user_config = ?", notificationConfig.NotificationConfigId, notificationConfig.UserConfig).
		Select("enabled", "message_template", "subject").
		Updates(
			&model.NotificationConfig{
				Enabled:         notificationConfig.Enabled,
				MessageTemplate: notificationConfig.MessageTemplate,
				Subject:         notificationConfig.Subject,
			},
		)

	err := m.BlockAccountList(ctx, req)
	err2 := m.UnBlockAccountList(ctx, req)

	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0 || err != nil || err2 != nil,
		fmt.Sprintf("Error: While editing notification config for userConfigId:%v error:%v", notificationConfig.UserConfig, result.Error),
		fmt.Sprintf("Success: Edited notification config for userConfigId:%v", notificationConfig.UserConfig),
		start)

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("notification config was not found userConfig:%v | NotificationConfig:%v",
			notificationConfig.UserConfig,
			notificationConfig.NotificationConfigId)
	}
	return nil
}

func (m *Mysql) BlockAccountList(ctx context.Context, req *nm.EditConfigReq) error {
	if len(req.BlockList) == 0 {
		m.Log.Info("No Accounts to be blocked . Block list empty")
		return nil
	}
	ntConfigId := req.NotificationConfig.NotificationConfigId

	successBlockList := []*nm.AccountMeta{}
	failedBlockList := []*nm.AccountMeta{}
	var err error
	for _, v := range req.BlockList {

		existingRule := model.AccountNotificationRules{}
		newRule := model.AccountNotificationRules{AccountConfigId: v.AccountConfigId, NotificationConfigId: ntConfigId, Disabled: true}
		result := m.DB.WithContext(ctx).FirstOrCreate(&existingRule, &newRule)

		if result.RowsAffected > 0 && result.Error == nil {
			successBlockList = append(successBlockList, v)
			m.Log.Infof("Blocked account config %v", v.AccountConfigId)
		} else {
			failedBlockList = append(failedBlockList, v)
			m.Log.Errorw("Failed to insert in account rules or rule already exist")
			err = result.Error
		}
	}
	m.Log.Infof("Successfully blocked following account %v", successBlockList)
	m.Log.Infof("Failed to blocked following account %v", failedBlockList)
	return err
}

func (m *Mysql) UnBlockAccountList(ctx context.Context, req *nm.EditConfigReq) error {
	if len(req.UnblockList) == 0 {
		m.Log.Info("No Accounts to be Unblocked . UnBlock list empty")
		return nil
	}
	ntConfigId := req.NotificationConfig.NotificationConfigId

	successUnBlockList := []*nm.AccountMeta{}
	failedUnBlockList := []*nm.AccountMeta{}
	var err error
	for _, v := range req.UnblockList {

		result := m.DB.WithContext(ctx).Delete(&model.AccountNotificationRules{}, "account_config_id = ? AND notification_config_id = ?", v.AccountConfigId, ntConfigId)

		if result.RowsAffected > 0 && result.Error == nil {
			successUnBlockList = append(successUnBlockList, v)
			m.Log.Infof("Blocked account config %v", v.AccountConfigId)
		} else {
			failedUnBlockList = append(failedUnBlockList, v)
			m.Log.Errorw("Failed to insert in account rules or rule already exist")
			err = result.Error
		}
	}
	m.Log.Infof("Successfully blocked following account %v", successUnBlockList)
	m.Log.Infof("Failed to blocked following account %v", failedUnBlockList)
	return err
}

func (m *Mysql) GetConfig(ctx context.Context, confReq *nm.ConfigReq) (*nm.UserConfigResp, error) {

	fName := "GetConfig"
	start := time.Now()

	userConfigId, err := m.GetUserConfigId(ctx, confReq.UserId)

	if err != nil {
		m.Log.Errorw("Error: while fetching userConfigId for userId", confReq.UserId)
		return nil, err
	}

	if err != nil {
		m.Log.Errorw("Error: while fetching account meta for userId", confReq.UserId)
		return nil, err
	}

	notificationMeta, err := m.GetNotificationConfig(ctx, *userConfigId)

	if err != nil {
		m.Log.Errorw("Error: while fetching notification config for userId", confReq.UserId)
		return nil, err
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: While fetching notification config for userId:%v error:%v", confReq.UserId, err),
		fmt.Sprintf("Success: Edited notification config for userId:%v", confReq.UserId),
		start)

	return &nm.UserConfigResp{
		Id:                  confReq.UserId,
		DefaultEmail:        confReq.DefaultEmail,
		NotificationConfigs: notificationMeta}, nil
}

func (m *Mysql) GetAccountConfig(ctx context.Context, userConfigId uint64) ([]*nm.AccountMetaReq, error) {
	fName := "GetAccountConfig"
	start := time.Now()

	var accountMeta []*nm.AccountMetaReq
	err := m.DB.WithContext(ctx).Model(&model.AccountConfig{}).Where("config_id = ?", userConfigId).Select("account_id", "email").Scan(&accountMeta).Error

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: While fetching account config for userConfigId:%v error:%v", userConfigId, err),
		fmt.Sprintf("Success: Got account config for userConfigId:%v", userConfigId),
		start)

	return accountMeta, nil
}

func (m *Mysql) GetNotificationConfig(ctx context.Context, userConfigId uint64) ([]*nm.NotificationConfig, error) {

	fName := "GetNotificationConfig"
	start := time.Now()

	var notificationConfig []*nm.NotificationConfig
	var notificationModel []model.NotificationConfig

	err := m.DB.WithContext(ctx).Model(&model.NotificationConfig{}).
		Select("id", "event_type", "notification_type", "enabled", "user_config").
		Where("user_config = ?", userConfigId).
		Scan(&notificationModel).Error

	for _, v := range notificationModel {
		conf := &nm.NotificationConfig{
			NotificationConfigId: uint64(v.ID),
			EventType:            v.EventType,
			NotificationType:     v.NotificationType,
			Enabled:              v.Enabled,
			UserConfig:           uint64(v.UserConfig),
		}
		notificationConfig = append(notificationConfig, conf)
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: Retrieving enabled notification config for userConfig:%v err:%+v", userConfigId, err),
		fmt.Sprintf("Success: Retrieved enabled notification config for userConfig:%v", userConfigId),
		start)

	return notificationConfig, nil
}

func (m *Mysql) GetConfigDetails(ctx context.Context, req *nm.ConfigDetailsReq) (*nm.ConfigDetailsReply, error) {
	fName := "GetConfigDetails"
	start := time.Now()
	var confDetailsReply *nm.ConfigDetailsReply = &nm.ConfigDetailsReply{}
	var notificationConfig *nm.NotificationConfig
	var notificationModel model.NotificationConfig

	err := m.DB.WithContext(ctx).Model(&model.NotificationConfig{}).
		Select("id,event_type,notification_type,enabled,message_template,subject,user_config").
		Where("id = ?", req.NotificationId).
		First(&notificationModel).Error

	notificationConfig = &nm.NotificationConfig{
		NotificationConfigId: notificationModel.ID,
		EventType:            notificationModel.EventType,
		NotificationType:     notificationModel.NotificationType,
		Enabled:              notificationModel.Enabled,
		MessageTemplate:      notificationModel.MessageTemplate,
		Subject:              notificationModel.Subject,
		UserConfig:           notificationModel.UserConfig,
	}
	confDetailsReply.NotificationConfig = notificationConfig

	enableAccListErr := m.PopulateEnabledAccountList(ctx, notificationConfig.UserConfig, notificationModel.ID, confDetailsReply)
	if enableAccListErr != nil {
		m.Log.Errorw("Error while populating enabled account list error: ", enableAccListErr)
	}

	blockedAccListErr := m.PopulateBlockedAccountList(ctx, notificationModel.ID, confDetailsReply)
	if blockedAccListErr != nil {
		m.Log.Errorw("Error while populating blocked account list error: ", blockedAccListErr)
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: Retrieving notification config details notificationConfigId:%v err:%+v", req.NotificationId, err),
		fmt.Sprintf("Success: Retrieved notification config details for notificationConfigId:%v", req.NotificationId),
		start)

	return confDetailsReply, err
}

func (m *Mysql) PopulateBlockedAccountList(ctx context.Context, notificationConfId uint64, data *nm.ConfigDetailsReply) error {
	fName := "PopulateBlockedAccountList"
	start := time.Now()
	var accountList []model.AccountConfig
	err := m.DB.WithContext(ctx).Table("notification_configs").
		Select("account_configs.account_id, account_configs.id", "account_configs.account_name").
		Joins("JOIN account_notification_rules ON account_notification_rules.notification_config_id = notification_configs.id").
		Joins("JOIN account_configs ON account_configs.id = account_notification_rules.account_config_id").
		Where("notification_configs.id = ?", notificationConfId).
		Scan(&accountList).Error

	data.BlockList = make([]*nm.AccountMeta, 0)

	for _, v := range accountList {

		blockedAccount := &nm.AccountMeta{
			AccountId:       v.AccountId,
			AccountName:     v.AccountName,
			AccountConfigId: v.ID,
		}

		data.BlockList = append(data.BlockList, blockedAccount)
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: Generating notification config account block list notificationConfigId:%v err:%+v", notificationConfId, err),
		fmt.Sprintf("Success: Retrieved notification config details for notificationConfigId:%v", notificationConfId),
		start)
	return err
}

func (m *Mysql) PopulateEnabledAccountList(ctx context.Context, userConfig, notificationConfId uint64, data *nm.ConfigDetailsReply) error {

	fName := "PopulateEnabledAccountList"
	start := time.Now()

	var accountList []model.AccountConfig

	subQuery := m.DB.WithContext(ctx).Table("account_notification_rules").
		Select("account_config_id").
		Where("account_config_id = a.id AND notification_config_id = ?", notificationConfId)

	err := m.DB.WithContext(ctx).Table("account_configs a").
		Select("a.account_id, a.id as id, a.account_name").
		Where("config_id = ? AND id NOT IN (?)", userConfig, subQuery).
		Scan(&accountList).Error

	data.EnabledList = make([]*nm.AccountMeta, 0)

	for _, v := range accountList {

		enabledAccount := &nm.AccountMeta{
			AccountId:       v.AccountId,
			AccountName:     v.AccountName,
			AccountConfigId: v.ID,
		}

		data.EnabledList = append(data.EnabledList, enabledAccount)
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: Generating notification config account enabled list notificationConfigId:%v err:%+v", notificationConfId, err),
		fmt.Sprintf("Success: Retrieved notification config account enabled list notificationConfigId:%v", notificationConfId),
		start)

	return err
}

func (m *Mysql) DeleteConfig(ctx context.Context, req *nm.DeleteConfigReq) error {

	fName := "DeleteConfig"
	start := time.Now()

	result := m.DB.WithContext(ctx).Delete(&model.NotificationConfig{}, "id = ? AND user_config = ?", req.NotificationConfigId, req.UserConfigId)

	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error: While deleting notification config for userConfig:%v err:%+v", req.UserConfigId, result.Error),
		fmt.Sprintf("Success: Deleted notification config for userConfig:%v", req.UserConfigId),
		start)

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("user config was not found userConfig:%v notificationConfigId:%v", req.UserConfigId, req.NotificationConfigId)
	}

	return nil
}

func (m *Mysql) EditAccountConfig(ctx context.Context, meta *nm.AccountMetaReq) error {

	fName := "EditAccountConfig"
	start := time.Now()

	result := m.DB.WithContext(ctx).Debug().Model(&model.AccountConfig{}).
		Where("account_id = ?", meta.AccountId).
		Select("email").
		Updates(
			&model.AccountConfig{
				Email: &meta.Email,
			},
		)

	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error: While editing account config for accountId:%v error:%v", meta.AccountId, result.Error),
		fmt.Sprintf("Success: Edited account config for accountId:%v", meta.AccountId),
		start)

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("account id was not found account_id:%v", meta.AccountId)
	}

	return nil
}

func (m *Mysql) AddDefaultNotificationConfig(ctx context.Context, userConfigId uint, notificationType string, eventType []string, callback ...func(uint64)) error {

	fName := "AddDefaultNotificationConfig"
	start := time.Now()
	var defaultConfig []model.DefaultConfigs
	ntType := strings.ToLower(notificationType)

	if !contract.IsValidNotificationType(ntType) {
		m.Log.Errorw("Invalid notification type %v", ntType)
		return fmt.Errorf("invalid notification type")
	}

	result := m.DB.WithContext(ctx).Model(&model.DefaultConfigs{}).Select("event_type", "message_template", "subject").
		Where("event_type IN (?) and notification_type = ? ", eventType, notificationType).
		Scan(&defaultConfig)

	if result.Error != nil {
		m.Log.Errorw("Error while retrieving default config")
		return result.Error
	}
	m.Log.Infof("Default config : %v", defaultConfig)
	var res *gorm.DB
	for _, v := range defaultConfig {

		msgTemplate := strings.Replace(v.MessageTemplate, "\n", "\\n", -1)
		ntConf := model.NotificationConfig{
			EventType:        v.EventType,
			NotificationType: ntType,
			Enabled:          true,
			MessageTemplate:  msgTemplate,
			Subject:          v.Subject,
			UserConfig:       uint64(userConfigId),
		}

		//existingConf := model.NotificationConfig{}
		existingConfig := model.NotificationConfig{}
		m.Log.Info(ntType)
		if ntType == contract.EMAIL {
			condition := model.NotificationConfig{
				EventType:        v.EventType,
				NotificationType: ntType,
				UserConfig:       uint64(userConfigId),
			}
			res = m.DB.WithContext(ctx).Where(condition).Attrs(ntConf).FirstOrCreate(&existingConfig)
		} else {
			ntConf.UuId = uuid.New().String()
			res = m.DB.WithContext(ctx).Where(ntConf).Attrs(ntConf).FirstOrCreate(&existingConfig)
			m.Log.Infof("ntConfId %v", existingConfig.ID)
			callback[0](existingConfig.ID)

		}

		if res.Error == nil && res.RowsAffected > 0 {
			m.Log.Infof("New notification config created for userConfig: %v | userConfigId: %v", ntConf.UserConfig, ntConf.UserConfig)
		} else {
			m.Log.Infof("Notification Config already exits for userConfig %v | eventType: %v | notification_type: %v error %v",
				ntConf.UserConfig, ntConf.EventType, ntConf.NotificationType, res.Error)
		}

	}

	m.LogError(fName,
		res.Error != nil && res.RowsAffected == 0,
		fmt.Sprintf("Error: While adding default config error:%v", result.Error),
		fmt.Sprintf("Success: Added default config for userConfig:%v", userConfigId),
		start)

	return result.Error

}

func (m *Mysql) IsBotNotificationConfigExist(ctx context.Context, eventType string, botConfigId uint64) (bool, error) {
	var result bool

	err := m.DB.WithContext(ctx).Raw("select EXISTS(select 1 from bot_events_rules be join notification_configs nc on be.notification_config_id=nc.id where nc.event_type= ? and be.bot_config_id= ?) as result", eventType, botConfigId).Scan(&result).Error
	if err != nil {
		return false, err
	}
	return result, nil
}

func (m *Mysql) IngestDefaultConfigTable() {
	defaultConf := []model.DefaultConfigs{
		{
			EventType:        nm.EventType_ACCOUNT_ADDED.String(),
			MessageTemplate:  "<p>Hi %FIRST_NAME%,</p><p>Your trading account <strong>%ACCOUNT_NUMBER%</strong> has been successfully added to Traders Connect with the name <strong>%ACCOUNT_NAME%</strong>.</p>",
			Subject:          "Trading account added",
			NotificationType: contract.EMAIL,
		},
		{
			EventType:        nm.EventType_ACCOUNT_CONNECTION_ERROR.String(),
			MessageTemplate:  `<p>Hi %FIRST_NAME%,</p><p>Your trading account <strong>%ACCOUNT_NUMBER%</strong> has disconnected from Traders Connect. The current status is as below:</p><p><strong>Status</strong> - %CONNECTION_STATUS%</p><p><strong>Error</strong> - %CONNECTION_ERROR%</p><p>If this is unexpected please reach out to our support team.</p><p><br></p><p><br></p>`,
			Subject:          "Trading account disconnected",
			NotificationType: contract.EMAIL,
		},
		{
			EventType:        nm.EventType_ACCOUNT_CONNECTED.String(),
			MessageTemplate:  `<p>Hi %FIRST_NAME%</p><p>Your trading account <strong>%ACCOUNT_NUMBER%</strong> has successfully reconnected to Traders Connect.</p><p>The current connection status is - %CONNECTION_STATUS%</p>`,
			Subject:          "Trading account connected",
			NotificationType: contract.EMAIL,
		},
		{
			EventType:        nm.EventType_ACCOUNT_DELETED.String(),
			MessageTemplate:  "<p>Hi %FIRST_NAME%,</p><p>Your trading account <strong>%ACCOUNT_NUMBER%</strong> (%ACCOUNT_NAME%) has been successfully deleted.</p>",
			Subject:          "Trading account deleted",
			NotificationType: contract.EMAIL,
		},
		{
			EventType:        nm.EventType_TRADE_COPY_FAILURE.String(),
			MessageTemplate:  "<p>Hi %FIRST_NAME%,</p><p>Your trade copy action has failed.</p><p>We attempted to copy from account <strong>%COPIER_MASTER%</strong> to account <strong>%COPIER_SLAVE%</strong> but faced the error shown below:</p><p><strong>Error</strong> - %COPIER_ERROR%</p><p>If this is unexpected, or you are unsure what the error means, please reach out to our support team.</p>",
			Subject:          "Trade copy failure",
			NotificationType: contract.EMAIL,
		},
		{
			EventType:        nm.EventType_ACCOUNT_ENABLED.String(),
			MessageTemplate:  "<p>Hi %FIRST_NAME%,</p><p>Your trading account %ACCOUNT_NUMBER% has been successfully enabled.</p>",
			Subject:          "Trading account enabled",
			NotificationType: contract.EMAIL,
		},
		{
			EventType:        nm.EventType_ACCOUNT_DISABLED.String(),
			MessageTemplate:  "<p>Hi %FIRST_NAME%,</p><p>Your trading account %ACCOUNT_NUMBER% has been successfully disabled.</p>",
			Subject:          "Trading account disabled",
			NotificationType: contract.EMAIL,
		},

		//Telegram
		{
			EventType:        nm.EventType_ACCOUNT_ADDED.String(),
			MessageTemplate:  "Hi %FIRST_NAME%, \\nYour trading account *%ACCOUNT_NUMBER%* has been successfully added to Traders Connect with the name *%ACCOUNT_NAME%*.",
			Subject:          "Trading account added",
			NotificationType: contract.TELEGRAM,
		},
		{
			EventType:        nm.EventType_ACCOUNT_CONNECTION_ERROR.String(),
			MessageTemplate:  `Hi %FIRST_NAME%, \\nYour trading account *%ACCOUNT_NUMBER%* has disconnected from Traders Connect. The current status is as below: \\n*Status* - %CONNECTION_STATUS% \\n*Error* - %CONNECTION_ERROR% \\nIf this is unexpected please reach out to our support team.`,
			Subject:          "Trading account disconnected",
			NotificationType: contract.TELEGRAM,
		},
		{
			EventType:        nm.EventType_ACCOUNT_CONNECTED.String(),
			MessageTemplate:  `Hi %FIRST_NAME% \\nYour trading account *%ACCOUNT_NUMBER%* has successfully reconnected to Traders Connect. \\nThe current connection status is - %CONNECTION_STATUS%`,
			Subject:          "Trading account connected",
			NotificationType: contract.TELEGRAM,
		},
		{
			EventType:        nm.EventType_ACCOUNT_DELETED.String(),
			MessageTemplate:  "Hi %FIRST_NAME%, \\nYour trading account *%ACCOUNT_NUMBER%* (%ACCOUNT_NAME%) has been successfully deleted.",
			Subject:          "Trading account deleted",
			NotificationType: contract.TELEGRAM,
		},
		{
			EventType:        nm.EventType_TRADE_COPY_FAILURE.String(),
			MessageTemplate:  "Hi %FIRST_NAME%, \\nYour trade copy action has failed. \\nWe attempted to copy from account *%COPIER_MASTER%* to account *%COPIER_SLAVE%* but faced the error shown below: \\n*Error* - %COPIER_ERROR% \\nIf this is unexpected, or you are unsure what the error means, please reach out to our support team.",
			Subject:          "Trade copy failure",
			NotificationType: contract.TELEGRAM,
		},
		{
			EventType:        nm.EventType_ACCOUNT_ENABLED.String(),
			MessageTemplate:  "Hi %FIRST_NAME%, \\nYour trading account %ACCOUNT_NUMBER% has been successfully enabled.",
			Subject:          "Trading account enabled",
			NotificationType: contract.TELEGRAM,
		},
		{
			EventType:        nm.EventType_ACCOUNT_DISABLED.String(),
			MessageTemplate:  "Hi %FIRST_NAME%, \\nYour trading account %ACCOUNT_NUMBER% has been successfully disabled.",
			Subject:          "Trading account disabled",
			NotificationType: contract.TELEGRAM,
		},
	}
	created := 0
	for _, v := range defaultConf {
		var existingDefaultConf model.DefaultConfigs
		res := m.DB.FirstOrCreate(&existingDefaultConf, &v)
		if res.Error == nil && res.RowsAffected > 0 {
			created += 1
			m.Log.Infof("Default Config created for eventType %v", v.EventType)
		} else {
			m.Log.Infof("Default Config already exist for eventType %v", v.EventType)
		}
	}
	m.Log.Infof("Total default config added %v", created)
}

func (m *Mysql) EditConfigStatus(ctx context.Context, meta *nm.EditConfigStatusReq) error {

	fName := "EditConfigStatus"
	start := time.Now()

	result := m.DB.WithContext(ctx).Model(&model.NotificationConfig{}).
		Where("id = ?", meta.NotificationConfigId).
		Select("enabled").
		Updates(
			&model.NotificationConfig{
				Enabled: meta.Enabled,
			},
		)

	if (result.RowsAffected == 0 || result.Error != nil) && meta.Enabled {
		m.Log.Errorw("Error while enabling notificationConfig . Notification Config Doesn't Exist")
		m.Log.Infof("Ingesting default configs for eventType %v userId:%v", meta.EventType, meta.UserId)
		userConfigId, err := m.GetUserConfigId(ctx, meta.UserId)
		if err != nil {
			m.Log.Errorw("Error:", err)
			return err
		}

		err = m.AddDefaultNotificationConfig(ctx, uint(*userConfigId), contract.EMAIL, []string{meta.EventType})
		if err != nil {
			m.Log.Errorw("Error:", err)
			return err
		}
	}
	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error: While editing notification config status for notificationConfId:%v error:%v", meta.NotificationConfigId, result.Error),
		fmt.Sprintf("Success: Edited account config for accountId:%v", meta.NotificationConfigId),
		start)

	return nil
}

func (m *Mysql) DumpLog(log model.Logs) error {

	fName := "DumpLog"
	start := time.Now()

	err := m.DB.Create(&log).Error

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: While dumping notification Log"),
		fmt.Sprintf("Success: Dumped notification log"),
		start)

	return err
}
