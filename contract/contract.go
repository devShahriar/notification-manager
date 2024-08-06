package contract

import (
	"strings"

	nm "github.com/Traders-Connect/esb-contract/golang/notification_manager"
)

const (
	DEFAULT_EVENT             = "trade_failed"
	DEFAULT_NOTIFICATION_TYPE = "email"
	EMAIL                     = "email"
	TELEGRAM                  = "telegram"
	DISCORD                   = "discord"
	DefaultTelegramBot        = "Traders connect"
)

var NotificationType map[string]bool = map[string]bool{
	strings.ToLower(nm.NotificationType_EMAIL.String()):    true,
	strings.ToLower(nm.NotificationType_TELEGRAM.String()): true,
	strings.ToLower(nm.NotificationType_DISCORD.String()):  true,
	strings.ToLower(nm.NotificationType_SLACK.String()):    true,
	strings.ToLower(nm.NotificationType_WHATSAPP.String()): true,
}

type WorkerMeta struct {
	Name             string
	WorkerType       string
	Event            string
	Exchange         string
	Queue            string
	ExchangeType     string
	BindingKey       string // Routing key
	NotificationType string
}

type EnabledNotification struct {
	NotificationType string
}

// type EmailMeta struct {
// 	DefaultEmail    string
// 	Email           *string
// 	MessageTemplate string
// 	Subject         string
// 	EventType       string
// }

type EmailMeta struct {
	DefaultEmail    string  `gorm:"column:default_email"`
	FirstName       string  `gorm:"column:first_name"`
	Email           *string `gorm:"column:email"`
	MessageTemplate string  `gorm:"column:message_template"`
	Subject         string  `gorm:"column:subject"`
	EventType       string  // Assuming column name matches field name
}

func GetDefaultEventList() []string {
	return []string{
		nm.EventType_ACCOUNT_ADDED.String(),
		nm.EventType_ACCOUNT_CONNECTION_ERROR.String(),
		nm.EventType_ACCOUNT_CONNECTED.String(),
		nm.EventType_ACCOUNT_DELETED.String(),
		nm.EventType_TRADE_COPY_FAILURE.String(),
	}
}

func IsValidNotificationType(ntType string) bool {
	if _, ok := NotificationType[ntType]; ok {
		return ok
	}
	return false
}

type ConfigIds struct {
	UserConfigId  string `gorm:"column:config_id"`
	AccountConfId string `gorm:"column:id"`
	UserId        string `gorm:"column:user_id"`
}

type BotEventRule struct {
	BotConfigID          uint   `gorm:"column:bot_config_id"`
	NotificationConfigID uint   `gorm:"column:notification_config_id"`
	Enabled              bool   `gorm:"column:enabled"`
	EventType            string `gorm:"column:event_type"`
	NotificationType     string `gorm:"column:notification_type"`
}

type ChannelMeta struct {
	ChannelName string `gorm:"column:channel_name"`
	Id          uint64 `gorm:"column:id"`
}
