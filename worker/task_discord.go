package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/devShahriar/H"
	"github.com/devshahriar/notification-manager/contract"
	"github.com/devshahriar/notification-manager/model"
	"github.com/devshahriar/notification-manager/template"
	"gorm.io/datatypes"
)

type TaskSendDiscordNotification struct {
	*Worker
}

func (t *TaskSendDiscordNotification) SendDiscordNotification(userConfig, accId, eventType string, data []byte) error {

	ntMeta, err := t.Db.GetBotNotificationMeta(ctx, userConfig, eventType, contract.DISCORD)

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
			t.Logger.Errorw("Error while sending discord notification", sendErr)
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
			NotificationType: contract.DISCORD,
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

func (t *TaskSendDiscordNotification) Send(botToken, channelId, message string) error {

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return err
	}

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord connection:", err)
		return err
	}

	// Send a message to the specified channel
	t.Worker.Logger.Infof("message %+v", message)
	message = strings.ReplaceAll(message, "\\n", "\n")
	_, err = dg.ChannelMessageSend(channelId, message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}

	fmt.Println("Message sent!")

	// Wait for a termination signal (e.g., Ctrl+C)

	// Close the Discord session cleanly before exiting
	dg.Close()
	return nil
}
