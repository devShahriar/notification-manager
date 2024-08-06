package db

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	nm "github.com/Traders-Connect/esb-contract/golang/notification_manager"
	"github.com/devShahriar/H"
	"github.com/devshahriar/notification-manager/contract"
	"github.com/devshahriar/notification-manager/model"
)

// This will return bot_token,message_template,channel_id etc
// Bot Worker will send notification from this meta data
func (m *Mysql) GetBotNotificationMeta(ctx context.Context, userConfigId, eventType, notificationType string) ([]model.BotNotificationMeta, error) {

	fName := "GetBotNotificationMeta"
	start := time.Now()

	var results []model.BotNotificationMeta
	err := m.DB.WithContext(ctx).Table("bot_configs bc").
		Select("uc.first_name, bc.bot_token, cc.channel_id, nc.event_type, nc.message_template, nc.subject, nc.notification_type").
		Joins("join bot_events_rules ber on bc.id = ber.bot_config_id").
		Joins("join channel_rules cr on ber.id = cr.bot_event_rules_id").
		Joins("join notification_configs nc on ber.notification_config_id = nc.id").
		Joins("join channel_configs cc on cc.id = cr.channel_config_id").
		Joins("join user_configs uc on uc.id = bc.user_config").
		Where("bc.enabled = true AND nc.enabled = true AND cc.enabled = true AND bc.user_config = ? AND nc.event_type = ? AND nc.notification_type = ?",
			userConfigId, eventType, notificationType).Scan(&results).Error

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: While getting BotNotificationMeta:%v err:%+v", userConfigId, err),
		fmt.Sprintf("Success: Deleted notification config for userConfig:%v", userConfigId),
		start)

	return results, err
}

// Bot
func (m *Mysql) AddBot(ctx context.Context, req *nm.AddBotReq) error {
	fName := "AddBot"
	start := time.Now()

	userConfig, err := m.GetUserConfigId(ctx, req.UserId)
	if err != nil {
		m.Log.Errorw("Error while feting userConfig id for userId:", req.UserId)
		return err
	}

	if req.BotName == contract.DefaultTelegramBot {
		return fmt.Errorf("bot name conflict with default bot name:%v", contract.DefaultTelegramBot)
	}

	botMeta := model.BotConfigs{
		UserConfig:       *userConfig,
		BotName:          req.BotName,
		BotToken:         req.BotToken,
		BotDescription:   req.BotDescription,
		NotificationType: req.NotificationType,
		Enabled:          true,
	}

	result := m.DB.WithContext(ctx).Where(botMeta).Attrs(botMeta).FirstOrCreate(&model.BotConfigs{})

	if result.Error == nil && result.RowsAffected > 0 {
		m.Log.Info("New Bot config added")
	} else {
		m.Log.Errorw("Error: Bot Config already exists")
	}
	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: While adding BotConfig for userId:%v err:%+v", req.UserId, err),
		fmt.Sprintf("Success: Added new bot config for userId:%v", req.UserId),
		start)
	return result.Error
}

func (m *Mysql) EditBot(ctx context.Context, meta *nm.EditBotReq) error {

	fName := "EditBot"
	start := time.Now()

	if meta.BotName == contract.DefaultTelegramBot {
		return fmt.Errorf("bot name conflict with default bot name:%v", contract.DefaultTelegramBot)
	}

	result := m.DB.WithContext(ctx).Model(&model.BotConfigs{}).
		Where("id = ?", meta.BotConfigId).
		Select("bot_name", "bot_description", "bot_token").
		Updates(
			&model.BotConfigs{
				BotName:        meta.BotName,
				BotDescription: meta.BotDescription,
				BotToken:       meta.BotToken,
			},
		)

	if result.RowsAffected > 0 || result.Error == nil {
		m.Log.Infof("BotConfig edited for botConfigId %v", meta.BotConfigId)

	} else {
		m.Log.Errorw("Error while updating botConfig for botConfigId %v", meta.BotConfigId)
		return fmt.Errorf("bot config id doesn't exist %v", meta.BotConfigId)
	}
	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error while updating botConfig for botConfigId:%v error:%v", meta.BotConfigId, result.Error),
		fmt.Sprintf("Success: Edited botConfig for botConfigId:%v", meta.BotConfigId),
		start)

	return result.Error
}

func (m *Mysql) GetBot(ctx context.Context, meta *nm.GetBotsReq) (*nm.GetBotsReply, error) {

	fName := "GetBot"
	start := time.Now()

	var result []model.BotConfigs

	userConfig, err := m.GetUserConfigId(ctx, meta.UserId)
	if err != nil {
		m.Log.Errorw("Error while getting userConfig for UserId:", meta.UserId)
		return nil, fmt.Errorf("error while getting userConfig for UserId:%v", meta.UserId)
	}

	err = m.DB.WithContext(ctx).Model(&model.BotConfigs{}).
		Select("id", "bot_name", "bot_description", "bot_token", "enabled").
		Where("user_config = ? AND notification_type = ?", userConfig, meta.NotificationType).
		Scan(&result).Error

	var reply *nm.GetBotsReply
	botConfigs := []*nm.BotMeta{}

	for _, v := range result {

		botConfig := nm.BotMeta{
			BotConfigId:    v.ID,
			BotName:        v.BotName,
			BotDescription: v.BotDescription,
			BotToken:       H.If(v.BotName == contract.DefaultTelegramBot, "", v.BotToken),
			Enabled:        v.Enabled,
		}

		botConfigs = append(botConfigs, &botConfig)
	}
	reply = &nm.GetBotsReply{
		UserConfigId: *userConfig,
		BotConfigs:   botConfigs,
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error while getting botConfig for userId:%v error:%v", meta.UserId, err),
		fmt.Sprintf("Success: Got botConfig for userId:%v", meta.UserId),
		start)

	return reply, err
}

func (m *Mysql) EditBotStatus(ctx context.Context, meta *nm.EditBotStatusReq) error {

	fName := "EditBotStatus"
	start := time.Now()

	result := m.DB.WithContext(ctx).Model(&model.BotConfigs{}).
		Where("id = ?", meta.BotConfigId).
		Select("enabled").
		Updates(
			&model.BotConfigs{
				Enabled: meta.Enable,
			},
		)

	if result.RowsAffected == 0 || result.Error != nil {
		m.Log.Errorw("Error while enabling bot . Bot Config Doesn't Exist")
		return fmt.Errorf("error while enabling bot . Bot Config Doesn't Exist %v", meta.BotConfigId)
	} else {
		m.Log.Infof("Updated bot status for botConfigId %v", meta.BotConfigId)
	}

	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error: While updating  bot config status for botConfigId:%v error:%v", meta.BotConfigId, result.Error),
		fmt.Sprintf("Success: Updated Bot config status for botConfigId:%v", meta.BotConfigId),
		start)

	return result.Error
}

func (m *Mysql) DeleteBot(ctx context.Context, meta *nm.DeleteBotReq) error {

	fName := "DeleteBot"
	start := time.Now()

	valid := m.CheckValidUser(ctx, meta.UserId, meta.BotConfigId, &model.BotConfigs{})
	if !valid {
		return fmt.Errorf("invalid request for userId:%v , configId: %v", meta.UserId, meta.BotConfigId)
	}

	result := m.DB.WithContext(ctx).Where("id = ?", meta.BotConfigId).Delete(&model.BotConfigs{})

	if result.RowsAffected == 0 || result.Error != nil {
		m.Log.Errorw("Error while enabling bot . Bot Config Doesn't Exist")
		return fmt.Errorf("error while deleting bot . Bot Config Doesn't Exist %v", meta.BotConfigId)
	} else {
		m.Log.Infof("Deleted bot for botConfigId %v", meta.BotConfigId)
	}

	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error: While deleting  bot config status for botConfigId:%v error:%v", meta.BotConfigId, result.Error),
		fmt.Sprintf("Success: Deleted Bot config status for botConfigId:%v", meta.BotConfigId),
		start)

	return result.Error

}

//Channel

func (m *Mysql) AddChannel(ctx context.Context, req *nm.AddChannelReq) error {

	fName := "AddChannel"
	start := time.Now()

	userConfig, err := m.GetUserConfigId(ctx, req.UserId)
	if err != nil {
		m.Log.Errorw("Error while feting userConfig id for userId:", req.UserId)
		return err
	}

	channelMeta := model.ChannelConfig{
		UserConfig:         *userConfig,
		ChannelName:        req.ChannelName,
		ChannelId:          req.ChannelId,
		ChannelDescription: req.ChannelDescription,
		Enabled:            true,
	}
	condition := model.ChannelConfig{
		UserConfig: *userConfig,
		ChannelId:  req.ChannelId,
	}

	result := m.DB.WithContext(ctx).Where(condition).Attrs(channelMeta).FirstOrCreate(&model.ChannelConfig{})

	if result.Error == nil && result.RowsAffected > 0 {
		m.Log.Info("New Channel config added")
	} else {
		m.Log.Errorw("Error: Channel Config already exists")
		return fmt.Errorf("this channel: %v has already been added for userId:%v", req.ChannelId, req.UserId)
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: While adding Channel Config for userId:%v err:%+v", req.UserId, err),
		fmt.Sprintf("Success: Added new Channel config for userId:%v", req.UserId),
		start)

	return result.Error
}

func (m *Mysql) EditChannel(ctx context.Context, meta *nm.EditChannelReq) error {

	fName := "EditChannel"
	start := time.Now()

	result := m.DB.WithContext(ctx).Model(&model.ChannelConfig{}).
		Where("id = ?", meta.ChannelConfigId).
		Select("channel_name", "channel_description", "channel_id").
		Updates(
			&model.ChannelConfig{
				ChannelName:        meta.ChannelName,
				ChannelDescription: meta.ChannelDescription,
				ChannelId:          meta.ChannelId,
			},
		)

	if result.RowsAffected > 0 || result.Error == nil {
		m.Log.Infof("ChannelConfig edited for ChannelConfigId %v", meta.ChannelConfigId)
	} else {
		m.Log.Errorw("Error while updating ChannelConfig for ChannelConfigId %v", meta.ChannelConfigId)
		return fmt.Errorf("channel config id doesn't exist %v", meta.ChannelConfigId)
	}
	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error while updating ChannelConfig for ChannelConfigId:%v error:%v", meta.ChannelConfigId, result.Error),
		fmt.Sprintf("Success: Edited ChannelConfig for ChannelConfigId:%v", meta.ChannelConfigId),
		start)

	return result.Error
}

func (m *Mysql) GetChannel(ctx context.Context, meta *nm.GetChannelReq) (*nm.GetChannelReply, error) {

	fName := "GetChannel"
	start := time.Now()

	var result []model.ChannelConfig

	userConfig, err := m.GetUserConfigId(ctx, meta.UserId)
	if err != nil {
		m.Log.Errorw("Error while getting userConfig for UserId:", meta.UserId)
		return nil, fmt.Errorf("error while getting userConfig for UserId:%v", meta.UserId)
	}

	err = m.DB.WithContext(ctx).Model(&model.ChannelConfig{}).
		Select("id", "channel_name", "channel_description", "channel_id", "enabled").
		Where("user_config = ?", userConfig).
		Scan(&result).Error

	var reply *nm.GetChannelReply
	channelConfigs := []*nm.ChannelConfig{}

	for _, v := range result {
		channelConfig := nm.ChannelConfig{
			ChannelConfigId:    v.ID,
			ChannelName:        v.ChannelName,
			ChannelDescription: v.ChannelDescription,
			ChannelId:          v.ChannelId,
			Enabled:            v.Enabled,
		}
		channelConfigs = append(channelConfigs, &channelConfig)
	}
	reply = &nm.GetChannelReply{
		UserConfigId:  *userConfig,
		ChannelConfig: channelConfigs,
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error while getting ChannelConfig for userId:%v error:%v", meta.UserId, err),
		fmt.Sprintf("Success: Got ChannelConfig for userId:%v", meta.UserId),
		start)

	return reply, err
}

func (m *Mysql) EditChannelStatus(ctx context.Context, meta *nm.EditChannelStatusReq) error {

	fName := "EditChannelStatus"
	start := time.Now()

	result := m.DB.WithContext(ctx).Model(&model.ChannelConfig{}).
		Where("id = ?", meta.ChannelConfigId).
		Select("enabled").
		Updates(
			&model.ChannelConfig{
				Enabled: meta.Enabled,
			},
		)

	if result.RowsAffected == 0 || result.Error != nil {
		m.Log.Errorw("Error while enabling channel . channle Config Doesn't Exist")
		return fmt.Errorf("error while enabling channel . channel Config Doesn't Exist %v", meta.ChannelConfigId)
	} else {
		m.Log.Infof("Updated channel status for ChannelConfigId %v", meta.ChannelConfigId)
	}

	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error: While updating  channel config status for ChannelConfigId:%v error:%v", meta.ChannelConfigId, result.Error),
		fmt.Sprintf("Success: Updated channel config status for ChannelConfigId:%v", meta.ChannelConfigId),
		start)

	return result.Error
}

func (m *Mysql) DeleteChannel(ctx context.Context, meta *nm.DeleteChannelReq) error {

	fName := "DeleteChannel"
	start := time.Now()

	valid := m.CheckValidUser(ctx, meta.UserId, meta.ChannelConfigId, &model.ChannelConfig{})
	if !valid {
		return fmt.Errorf("invalid request for userId:%v , configId: %v", meta.UserId, meta.ChannelConfigId)
	}
	result := m.DB.WithContext(ctx).Where("id = ?", meta.ChannelConfigId).Delete(&model.ChannelConfig{})

	if result.RowsAffected == 0 || result.Error != nil {
		m.Log.Errorw("Error while enabling Chanel . Channel Config Doesn't Exist")
		return fmt.Errorf("error while deleting Chanel . Channel Config Doesn't Exist %v", meta.ChannelConfigId)
	} else {
		m.Log.Infof("Deleted channel for ChannelConfigId %v", meta.ChannelConfigId)
	}

	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error: While deleting  bot config status for ChannelConfigId:%v error:%v", meta.ChannelConfigId, result.Error),
		fmt.Sprintf("Success: Deleted Bot config status for ChannelConfigId:%v", meta.ChannelConfigId),
		start)

	return result.Error

}

func (m *Mysql) CheckValidUser(ctx context.Context, reqUserId string, configId uint64, tableName interface{}) bool {

	column := "user_config"

	var userConfig string
	err := m.DB.WithContext(ctx).Model(tableName).Where("id = ? ", configId).Select(column).Scan(&userConfig).Error
	if err != nil {
		m.Log.Error(err)
		return false
	}

	var actualUserId string
	err1 := m.DB.WithContext(ctx).Model(&model.UserConfig{}).Where("id = ?", userConfig).Select("user_id").Scan(&actualUserId).Error
	if err1 != nil {
		m.Log.Errorw("-", "error", err1)
		return false
	}

	m.Log.Errorw("-", "acId", actualUserId)
	valid := H.If(actualUserId == reqUserId, true, false)
	return valid
}

func (m *Mysql) GetBotEventConfigs(ctx context.Context, req *nm.GetBotEventConfigsReq) (*nm.GetBotEventConfigsReply, error) {

	fName := "GetBotEventConfigs"
	start := time.Now()

	var botEventRules []contract.BotEventRule

	err := m.DB.WithContext(ctx).Table("bot_events_rules be").
		Select("be.bot_config_id, be.notification_config_id, nc.enabled, nc.event_type, nc.notification_type").
		Joins("JOIN notification_configs nc ON be.notification_config_id = nc.id").
		Joins("JOIN bot_configs bc ON be.bot_config_id = bc.id").
		Where("be.bot_config_id = ? AND bc.notification_type = ?", req.BotConfigId, req.NotificationType).
		Scan(&botEventRules).Error

	if err != nil {
		m.Log.Errorw("Error while getting BotEventConfigs for botConfigId %v", req.BotConfigId)
		return nil, fmt.Errorf("error while getting BotEventConfigs for botConfigId %v", req.BotConfigId)
	}

	ntConfigs := []*nm.NotificationConfig{}

	for _, v := range botEventRules {

		ntConfig := &nm.NotificationConfig{
			EventType:            v.EventType,
			NotificationConfigId: uint64(v.NotificationConfigID),
			NotificationType:     v.NotificationType,
			Enabled:              v.Enabled,
		}
		ntConfigs = append(ntConfigs, ntConfig)
	}

	reply := &nm.GetBotEventConfigsReply{
		BotConfigId:         req.BotConfigId,
		NotificationConfigs: ntConfigs,
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: While getting  bot GetBotEventConfigs for botConfigId:%v error:%v", req.BotConfigId, err),
		fmt.Sprintf("Success: Got GetBotEventConfigs for botConfigId:%v", req.BotConfigId),
		start)

	return reply, err
}

func (m *Mysql) GetBotEventDetails(ctx context.Context, req *nm.GetBotEventDetailsReq) (*nm.GetBotEventDetailsReply, error) {

	fName := "GetBotEventDetails"
	start := time.Now()

	userConfig, err := m.GetUserConfigId(ctx, req.UserId)
	if err != nil {
		m.Log.Errorw("Error while getting userConfig for UserId:", req.UserId)
		return nil, fmt.Errorf("error while getting userConfig for UserId:%v", req.UserId)
	}
	reply := &nm.GetBotEventDetailsReply{}

	//Populate Account Block/Unblock list
	accountList := &nm.ConfigDetailsReply{}
	err1 := m.PopulateEnabledAccountList(ctx, *userConfig, req.NotificationConfigId, accountList)
	err2 := m.PopulateBlockedAccountList(ctx, req.NotificationConfigId, accountList)

	if err1 != nil || err2 != nil {
		m.Log.Errorw("Error while getting account block/unblock list", err1, err2)
		return nil, fmt.Errorf("error while getting account block/unblock list for botConfigId %v notificationConfigId %v", req.BotConfigId, req.NotificationConfigId)
	}

	//Populate Channel Block/Unblock list
	err3 := m.PopulateBlockChannelList(ctx, *userConfig, req.BotConfigId, req.NotificationConfigId, reply)
	err4 := m.PopulateUnBlockChannelList(ctx, *userConfig, req.BotConfigId, req.NotificationConfigId, reply)

	if err3 != nil || err4 != nil {
		m.Log.Errorw("Error while getting account block/unblock list", err1, err2)
		return nil, fmt.Errorf("error while getting channel block/unblock list for botConfigId %v notificationConfigId %v", req.BotConfigId, req.NotificationConfigId)
	}

	//Get Notification Config
	var ntConf model.NotificationConfig
	err = m.DB.WithContext(ctx).Model(&model.NotificationConfig{}).
		Select("id", "event_type", "notification_type", "message_template", "subject", "enabled").
		Where("id = ?", req.NotificationConfigId).Scan(&ntConf).Error

	if err != nil {
		m.Log.Errorw("Error while getting notification config details", err)
		return nil, fmt.Errorf("error while getting notificationConfig for botConfigId %v notificationConfigId %v", req.BotConfigId, req.NotificationConfigId)
	}
	reply.NotificationConfig = &nm.NotificationConfig{
		NotificationConfigId: ntConf.ID,
		EventType:            ntConf.EventType,
		NotificationType:     ntConf.NotificationType,
		MessageTemplate:      ntConf.MessageTemplate,
		Subject:              ntConf.Subject,
		Enabled:              ntConf.Enabled,
	}

	reply.AccountBlockList = accountList.BlockList
	reply.AccountEnabledList = accountList.EnabledList

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: While getting  bot GetBotEventDetails for botConfigId:%v error:%v", req.BotConfigId, err),
		fmt.Sprintf("Success: Got GetBotEventDetails for botConfigId:%v", req.BotConfigId),
		start)

	return reply, nil
}

func (m *Mysql) PopulateBlockChannelList(ctx context.Context, userConfig, botConfigId, notificationConfigId uint64, data *nm.GetBotEventDetailsReply) error {

	fName := "PopulateBlockChannelList"
	start := time.Now()

	var excludedChannelIDs []uint64

	subQuery := m.DB.WithContext(ctx).Model(&model.BotEventsRules{}).
		Select("cr.channel_config_id as id").
		Joins("JOIN channel_rules cr ON bot_events_rules.id = cr.bot_event_rules_id").
		Joins("JOIN channel_configs cc ON cc.id = cr.channel_config_id").
		Where("bot_events_rules.bot_config_id = ? AND bot_events_rules.notification_config_id = ?", botConfigId, notificationConfigId).
		Pluck("id", &excludedChannelIDs)

	if subQuery.Error != nil {
		return subQuery.Error
	}

	var channels []contract.ChannelMeta
	if len(excludedChannelIDs) == 0 || excludedChannelIDs == nil {
		excludedChannelIDs = []uint64{0}
	}
	err := m.DB.WithContext(ctx).Model(&model.ChannelConfig{}).Select("id", "channel_name").
		Where("id NOT IN (?)", excludedChannelIDs).
		Where("user_config = ?", userConfig).Scan(&channels).Error

	if err != nil {
		return err
	}

	data.ChannelBlockList = []*nm.ChannelMeta{}

	for _, v := range channels {
		channelMeta := &nm.ChannelMeta{
			ChannelConfigId: v.Id,
			ChannelName:     v.ChannelName,
		}
		data.ChannelBlockList = append(data.ChannelBlockList, channelMeta)
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: error while populating block channel list for userConfig:%v,ntConfig:%v,botConf:%v", userConfig, notificationConfigId, botConfigId),
		fmt.Sprintf("Success: Got block channel list for botConfigId:%v", botConfigId),
		start)
	return err
}

func (m *Mysql) PopulateUnBlockChannelList(ctx context.Context, userConfig, botConfigId, notificationConfigId uint64, data *nm.GetBotEventDetailsReply) error {

	fName := "PopulateUnBlockChannelList"
	start := time.Now()
	var channels []contract.ChannelMeta

	err := m.DB.WithContext(ctx).Model(&model.BotEventsRules{}).
		Select("cc.id, cc.channel_name").
		Joins("JOIN channel_rules cr ON bot_events_rules.id = cr.bot_event_rules_id").
		Joins("JOIN channel_configs cc ON cc.id = cr.channel_config_id").
		Where("bot_events_rules.bot_config_id = ? AND bot_events_rules.notification_config_id = ?", botConfigId, notificationConfigId).
		Find(&channels).Error

	data.ChannelEnabledList = []*nm.ChannelMeta{}

	for _, v := range channels {
		channelMeta := &nm.ChannelMeta{
			ChannelConfigId: v.Id,
			ChannelName:     v.ChannelName,
		}
		data.ChannelEnabledList = append(data.ChannelEnabledList, channelMeta)
	}

	m.LogError(fName,
		err != nil,
		fmt.Sprintf("Error: error while populating block channel list for userConfig:%v,ntConfig:%v,botConf:%v", userConfig, notificationConfigId, botConfigId),
		fmt.Sprintf("Success: Got block channel list for botConfigId:%v", botConfigId),
		start)
	return err
}

func ValidateEditBotEventReq(data string) bool {

	return data != ""
}

func DoesContainHTML(data string) bool {
	htmlTagPattern := `<[^>]+>`

	// Create a regex pattern
	re := regexp.MustCompile(htmlTagPattern)

	// Use the regex pattern to check if the input string contains HTML tags
	htmlContent := re.MatchString(data)
	return htmlContent
}
func (m *Mysql) EditBotEventDetails(ctx context.Context, req *nm.EditBotEventDetailsReq) error {

	fName := "EditBotEventDetails"
	start := time.Now()

	notificationConfig := req.NotificationConfig

	if DoesContainHTML(req.NotificationConfig.MessageTemplate) {
		return fmt.Errorf("HTML template not allowed")
	}
	if ValidateEditBotEventReq(req.NotificationConfig.MessageTemplate) && !DoesContainHTML(req.NotificationConfig.MessageTemplate) {
		msgTemplate := strings.Replace(notificationConfig.MessageTemplate, "\n", "\\n", -1)

		if req.NotificationConfig.NotificationType == nm.NotificationType_DISCORD.String() {
			msgTemplate = strings.Replace(msgTemplate, "**", "*", -1)
		}

		result := m.DB.WithContext(ctx).Debug().Model(&model.NotificationConfig{}).
			Where("id = ? ", notificationConfig.NotificationConfigId).
			Select("enabled", "message_template", "subject").
			Updates(
				&model.NotificationConfig{
					Enabled:         notificationConfig.Enabled,
					MessageTemplate: msgTemplate,
					Subject:         notificationConfig.Subject,
				},
			)

		if result.Error != nil || result.RowsAffected == 0 {
			return fmt.Errorf("notification config was not found userConfig:%v | NotificationConfig:%v",
				notificationConfig.UserConfig,
				notificationConfig.NotificationConfigId)
		}
	}

	accountList := nm.EditConfigReq{
		BlockList:          req.AccountBlockList,
		UnblockList:        req.AccountEnabledList,
		NotificationConfig: req.NotificationConfig,
	}

	err := m.BlockAccountList(ctx, &accountList)
	err2 := m.UnBlockAccountList(ctx, &accountList)

	err3 := m.BlockChannel(ctx, req.BotConfigId, req.NotificationConfig.NotificationConfigId, req.ChannelBlockList)

	err4 := m.UnBlockChannel(ctx, req.BotConfigId, req.NotificationConfig.NotificationConfigId, req.ChannelEnabledList)

	if err != nil || err2 != nil || err3 != nil || err4 != nil {
		return fmt.Errorf("error while blocking/unblocking channel/account")
	}

	m.LogError(fName,
		err != nil || err2 != nil,
		fmt.Sprintf("Error: While editing notification config for userConfigId:%v", notificationConfig.UserConfig),
		fmt.Sprintf("Success: Edited notification config for userConfigId:%v", notificationConfig.UserConfig),
		start)

	return nil
}

func (m *Mysql) UnBlockChannel(ctx context.Context, botConfigId, ntConfigId uint64, channels []*nm.ChannelMeta) error {

	if len(channels) == 0 {
		m.Log.Info("No channel found to block")
		return nil
	}

	//Get BotEventRuleId
	var botEventRuleId uint64
	err := m.DB.WithContext(ctx).Model(&model.BotEventsRules{}).
		Select("id").
		Where("bot_config_id = ? AND notification_config_id = ?", botConfigId, ntConfigId).
		Scan(&botEventRuleId).Error

	if err != nil {
		m.Log.Errorw("Error: while getting botEventRuleId for botConfigId %v,ntConfigId %v", botConfigId, ntConfigId)
		return err
	}

	//Insert into channelRules
	for _, v := range channels {
		chanelRule := model.ChannelRules{
			BotEventRulesId: botEventRuleId,
			ChannelConfigId: v.ChannelConfigId,
		}

		err := m.DB.WithContext(ctx).Where(chanelRule).Attrs(chanelRule).FirstOrCreate(&model.ChannelRules{}).Error
		if err != nil {
			return err
		}

	}

	return nil
}

func (m *Mysql) BlockChannel(ctx context.Context, botConfigId, ntConfigId uint64, channels []*nm.ChannelMeta) error {

	if len(channels) == 0 {
		m.Log.Info("No channel found to block")
		return nil
	}

	//Get BotEventRuleId
	var botEventRuleId uint64
	err := m.DB.WithContext(ctx).Model(&model.BotEventsRules{}).
		Select("id").
		Where("bot_config_id = ? AND notification_config_id = ?", botConfigId, ntConfigId).
		Scan(&botEventRuleId).Error

	if err != nil {
		m.Log.Errorw("Error: while getting botEventRuleId for botConfigId %v,ntConfigId %v", botConfigId, ntConfigId)
		return err
	}

	//Insert into channelRules
	for _, v := range channels {
		chanelRule := model.ChannelRules{
			BotEventRulesId: botEventRuleId,
			ChannelConfigId: v.ChannelConfigId}

		err := m.DB.WithContext(ctx).Delete(&model.ChannelRules{}, "bot_event_rules_id = ? AND channel_config_id = ?", chanelRule.BotEventRulesId, chanelRule.ChannelConfigId).Error
		if err != nil {
			return err
		}

	}

	return nil
}

func (m *Mysql) EditBotEventStatus(ctx context.Context, meta *nm.EditBotEventStatusReq) error {

	fName := "EditBotEventStatus"
	start := time.Now()

	callbacks := func(ntConfId uint64) {

		userConfig, err := m.GetUserConfigId(ctx, meta.UserId)
		if err != nil {
			m.Log.Errorw("Error", err)
			return
		}
		m.AddBotEventRules(ctx, meta.BotConfigId, ntConfId)
		m.BlockAccountAndChannel(ctx, *userConfig, ntConfId)
	}

	result := m.DB.WithContext(ctx).Model(&model.NotificationConfig{}).
		Where("id = ?", meta.NotificationConfigId).
		Select("enabled").
		Updates(
			&model.NotificationConfig{
				Enabled: meta.Enabled,
			},
		)

	if (result.RowsAffected == 0 || result.Error != nil) && meta.Enabled {

		exist, err := m.IsBotNotificationConfigExist(ctx, meta.EventType, meta.BotConfigId)

		if exist || err != nil {
			m.Log.Errorw("Bot Event Config Already Exist for botConfig:", meta.BotConfigId, "eventType", meta.EventType)
			return fmt.Errorf("bot config for this event already exits botConfigId:%v. eventType:%v", meta.BotConfigId, meta.EventType)
		}

		m.Log.Errorw("Error while enabling notificationConfig . Notification Config Doesn't Exist")
		m.Log.Infof("Ingesting default configs for userId:%v", meta.UserId)
		userConfigId, err := m.GetUserConfigId(ctx, meta.UserId)
		if err != nil {
			m.Log.Errorw("Error:", err)
			return err
		}

		m.Log.Info(meta.NotificationType.String())
		err = m.AddDefaultNotificationConfig(ctx, uint(*userConfigId), meta.NotificationType.String(), []string{meta.EventType}, callbacks)
		if err != nil {
			m.Log.Errorw("Error:", err)
			return err
		}
	}
	m.LogError(fName,
		result.Error != nil || result.RowsAffected == 0,
		fmt.Sprintf("Error: While editing notification config status for notificationConfId:%v error:%v", meta.NotificationConfigId, result.Error),
		fmt.Sprintf("Success: Edited account config status for accountId:%v", meta.NotificationConfigId),
		start)

	return nil
}

func (m *Mysql) BlockAccountAndChannel(ctx context.Context, userConfId, ntConfId uint64) error {

	req := &nm.EditConfigReq{
		NotificationConfig: &nm.NotificationConfig{NotificationConfigId: ntConfId},
	}
	accountList, err := m.GetAccountListByUserConfig(ctx, userConfId)

	if err != nil {
		m.Log.Errorw("Error", err)
	}
	req.BlockList = accountList
	m.BlockAccountList(ctx, req)
	return nil
}

func (m *Mysql) GetAccountListByUserConfig(ctx context.Context, userConfId uint64) ([]*nm.AccountMeta, error) {

	var accountList []*nm.AccountMeta
	err := m.DB.WithContext(ctx).Model(&model.AccountConfig{}).Select("id as account_config_id", "account_name").
		Where("config_id = ?", userConfId).Scan(&accountList).Error

	return accountList, err

}

func (m *Mysql) AddBotEventRules(ctx context.Context, botConfigId, ntConfId uint64) error {

	data := &model.BotEventsRules{
		BotConfigId:          botConfigId,
		NotificationConfigId: ntConfId,
	}

	err := m.DB.WithContext(ctx).Where(data).Attrs(data).FirstOrCreate(&model.BotEventsRules{}).Error
	if err != nil {
		m.Log.Errorw("Error:", err)
		return err
	}
	return err
}

func (m *Mysql) UninstallIntegration(ctx context.Context, req *nm.UninstallIntegrationReq) error {

	// delete notification config
	userConfig, err := m.GetUserConfigId(ctx, req.UserId)
	if err != nil {
		m.Log.Errorw("Error while getting userConfig for UserId:", req.UserId)
		return fmt.Errorf("error while getting userConfig for UserId:%v", req.UserId)
	}

	err = m.DeleteNotificationByUser(*userConfig, req.NotificationType.String())
	if err != nil {
		return err
	}

	if req.NotificationType != nm.NotificationType_EMAIL {

		// Delete Bot
		err1 := m.DB.WithContext(ctx).Where("user_config = ? AND notification_type = ?", userConfig, req.NotificationType.String()).Delete(&model.BotConfigs{}).Error
		// Delete Channel
		err2 := m.DB.WithContext(ctx).Where("user_config = ?", userConfig).Delete(&model.ChannelConfig{}).Error

		if err1 != nil || err2 != nil {
			return fmt.Errorf("error while deleting bot/channel for userId %v, notification type %v", req.UserId, req.NotificationType)
		}

	}

	//Update user config status column
	userConf := model.UserConfig{}
	column := GetColumnName(req.NotificationType.String())
	PopulateUserConfigDisabledNotifications(&userConf, req.NotificationType.String())
	result := m.DB.WithContext(ctx).Model(&model.UserConfig{}).
		Where("user_id = ?", req.UserId).
		Select(column).
		Updates(
			userConf,
		)

	if result.Error != nil && result.RowsAffected == 0 {
		return fmt.Errorf("error while uninstalling integration userId %v ntType %v", req.UserId, req.NotificationType)
	}
	return nil
}

func (m *Mysql) DeleteNotificationByUser(userConfigId uint64, ntType string) error {

	err := m.DB.Where("user_config = ? AND notification_type = ?", userConfigId, ntType).Delete(&model.NotificationConfig{}).Error
	if err != nil {
		return err
	}

	return nil
}

func PopulateUserConfigDisabledNotifications(userConfig *model.UserConfig, ntType string) {
	switch ntType {
	case nm.NotificationType_EMAIL.String():
		userConfig.EmailEnabled = false
		return
	case nm.NotificationType_TELEGRAM.String():
		userConfig.TelegramEnabled = false
		return
	case nm.NotificationType_DISCORD.String():
		userConfig.DiscordEnabled = false
		return
	case nm.NotificationType_SLACK.String():
		userConfig.SlackEnabled = false
		return
	case nm.NotificationType_WHATSAPP.String():
		userConfig.WhatsAppEnabled = false
		return
	}
}
