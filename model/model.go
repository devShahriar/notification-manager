package model

import (
	"time"

	"gorm.io/datatypes"
	_ "gorm.io/gorm"
)

type WorkerMeta struct {
	Id               uint64 `gorm:"primarykey"`
	Name             string
	WorkerType       string
	Exchange         string
	Queue            string
	ExchangeType     string
	BindingKey       string // Routing key
	NotificationType string
}

type EmailNotificationMeta struct {
	Id           uint64 `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	BodyTemplate string
	Subject      string
	UserConfig   string `gorm:"type:varchar(64)"`
	DefaultEmail string
}

type UserConfig struct {
	ID                  uint64 `gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	UserId              string `gorm:"type:varchar(100);uniqueIndex:idx_user_id"`
	FirstName           string
	DefaultEmail        string
	EmailEnabled        bool
	TelegramEnabled     bool
	DiscordEnabled      bool
	SlackEnabled        bool
	WhatsAppEnabled     bool
	AccountConfigs      []AccountConfig      `gorm:"foreignKey:ConfigId;references:ID;constraint:OnDelete:CASCADE"`
	NotificationConfigs []NotificationConfig `gorm:"foreignKey:UserConfig;references:ID;constraint:OnDelete:CASCADE"`
	BotConfigs          []BotConfigs         `gorm:"foreignKey:UserConfig;references:ID;constraint:OnDelete:CASCADE"`
	ChannelConfig       []ChannelConfig      `gorm:"foreignKey:UserConfig;references:ID;constraint:OnDelete:CASCADE"`
}

type AccountConfig struct {
	ID                       uint64 `gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	CreatedAt                time.Time
	UpdatedAt                time.Time
	ConfigId                 uint64  `gorm:"type:bigint(20)"`
	AccountId                string  `gorm:"type:varchar(100);uniqueIndex:idx_account_id"`
	Email                    *string `gorm:"null"`
	AccountName              string
	AccountNotificationRules []AccountNotificationRules `gorm:"foreignKey:AccountConfigId;references:ID;constraint:OnDelete:CASCADE"`
}

type NotificationConfig struct {
	ID                       uint64 `gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	CreatedAt                time.Time
	UpdatedAt                time.Time
	EventType                string
	NotificationType         string
	Enabled                  bool
	MessageTemplate          string
	Subject                  string
	UserConfig               uint64 `gorm:"type:bigint(20)"`
	UuId                     string
	AccountNotificationRules []AccountNotificationRules `gorm:"foreignKey:NotificationConfigId;references:ID;constraint:OnDelete:CASCADE"`
	BotEventsRules           []BotEventsRules           `gorm:"foreignKey:NotificationConfigId;references:ID;constraint:OnDelete:CASCADE"`
}

type Logs struct {
	Id               uint64 `gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UserConfig       string
	AccountId        string
	EventType        string
	NotificationType string
	ReqMeta          datatypes.JSON `gorm:"type:json"`
	Status           string
}

type BotNotificationMeta struct {
	FirstName        string `gorm:"column:first_name"`
	BotToken         string `gorm:"column:bot_token"`
	ChannelId        string `gorm:"column:channel_id"`
	EventType        string `gorm:"column:event_type"`
	MessageTemplate  string `gorm:"column:message_template"`
	Subject          string `gorm:"column:subject"`
	NotificationType string `gorm:"column:notification_type"`
}
type DefaultConfigs struct {
	Id               uint64 `gorm:"primaryKey"`
	EventType        string
	MessageTemplate  string
	Subject          string
	NotificationType string
}

type AccountNotificationRules struct {
	Id                   uint64 `gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Disabled             bool
	AccountConfigId      uint64 `gorm:"type:bigint(20)"`
	NotificationConfigId uint64 `gorm:"type:bigint(20)"`
}

// Telegram
type BotConfigs struct {
	ID               uint64 `gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	UserConfig       uint64 `gorm:"type:bigint(20)"`
	BotName          string
	BotDescription   string
	BotToken         string
	NotificationType string
	Enabled          bool
	BotEventsRules   []BotEventsRules `gorm:"foreignKey:BotConfigId;references:ID;constraint:OnDelete:CASCADE"`
}

type BotEventsRules struct {
	ID                   uint64 `gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	BotConfigId          uint64         `gorm:"type:bigint(20)"`
	NotificationConfigId uint64         `gorm:"type:bigint(20)"`
	ChannelRules         []ChannelRules `gorm:"foreignKey:BotEventRulesId;references:ID;constraint:OnDelete:CASCADE"`
}

type ChannelRules struct {
	ID              uint64 `gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	BotEventRulesId uint64
	ChannelConfigId uint64
}

type ChannelConfig struct {
	ID                 uint64 `gorm:"primaryKey;autoIncrement;type:bigint(20)"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	ChannelName        string
	ChannelDescription string
	ChannelId          string
	Enabled            bool
	UserConfig         uint64         `gorm:"type:bigint(20)"`
	ChannelRules       []ChannelRules `gorm:"foreignKey:ChannelConfigId;references:ID;constraint:OnDelete:CASCADE"`
}
