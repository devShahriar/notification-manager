package worker

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/devShahriar/H"
	"github.com/devshahriar/notification-manager/contract"
	"github.com/devshahriar/notification-manager/model"
	"github.com/devshahriar/notification-manager/template"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"gorm.io/datatypes"
)

type TaskSendTelegramNotification struct {
	*Worker
}

func (t *TaskSendTelegramNotification) SendTelegramNotification(userConfig, accId, eventType string, data []byte) error {

	ntMeta, err := t.Db.GetBotNotificationMeta(ctx, userConfig, eventType, contract.TELEGRAM)

	if err != nil {
		t.Logger.Errorw("Error: while retrieving BotNotification Meta for userConfig", userConfig, "eventType", eventType)
		return err
	}

	for _, v := range ntMeta {

		var dataObj map[string]string
		err := json.Unmarshal(data, &dataObj)

		if err != nil {
			t.Logger.Errorw("Error while parsing notification data", err)
			continue
		}

		message := template.IngestDataIntoMsgBody(v.FirstName, v.EventType, v.MessageTemplate, dataObj)

		sendErr := t.Send(v.BotToken, v.ChannelId, message)
		if sendErr != nil {
			t.Logger.Errorw("Error while sending telegram notification", sendErr)
			continue
		}

		reqMeta := struct {
			BotToken  string
			ChannelId string
			Message   string
		}{
			BotToken:  v.BotToken,
			ChannelId: v.ChannelId,
			Message:   message,
		}

		reqMetaBytes, _ := json.Marshal(reqMeta)

		dumpLogErr := t.Db.DumpLog(model.Logs{
			UserConfig:       userConfig,
			AccountId:        accId,
			EventType:        eventType,
			NotificationType: contract.TELEGRAM,
			ReqMeta:          datatypes.JSON(reqMetaBytes),
			Status:           H.If(sendErr != nil, "FAILED", "SUCCESS"),
		})

		if dumpLogErr != nil {
			t.Logger.Errorw(dumpLogErr.Error())
		}

		log.Println("Message sent successfully!")
	}
	return nil
}

func (t *TaskSendTelegramNotification) Send(botToken, channelId, message string) error {

	t.Logger.Info(message)
	bot, err := tgbotapi.NewBotAPI(botToken)

	if err != nil {
		t.Logger.Errorw("Failed to created new BotAPi for telegram")
		return err
	}

	chatID, err := strconv.Atoi(channelId)
	if err != nil {
		t.Logger.Errorw("Error while converting channel id %v", channelId)
		return err
	}

	message = strings.Replace(message, "\\n", "\n", -1) //Handle \\n from the request payload
	msg := tgbotapi.NewMessage(int64(chatID), message)

	msg.ParseMode = "Markdown"
	_, err = bot.Send(msg)

	if err != nil {
		t.Logger.Errorw("Failed to send telegram notification")
		return err
	}

	return nil
}
