package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Traders-Connect/esb-contract/golang/notification_manager"
	"github.com/Traders-Connect/notification-manager/contract"
	"github.com/Traders-Connect/notification-manager/model"
	"github.com/Traders-Connect/notification-manager/template"
	"github.com/devShahriar/H"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sirupsen/logrus"
)

var ctx context.Context = context.Background()

type TaskSendEmail struct {
	*Worker
}

func (t *TaskSendEmail) SendEmail(userConfig, accId, eventType string, data []byte) error {

	var emailMeta contract.EmailMeta
	var err error

	if eventType == notification_manager.EventType_ACCOUNT_DELETED.String() {
		t.Worker.Logger.Info("Sending user specific notification")
		emailMeta, err = t.Worker.Db.GetEmailMetaForUserOnly(ctx, userConfig, eventType)
	} else {
		emailMeta, err = t.Worker.Db.GetEmailMeta(ctx, userConfig, accId, eventType)
	}

	t.Worker.Logger.Info(emailMeta)

	if err != nil {
		logrus.Info("Couldn't fetch email")
		return err
	}

	t.Worker.Logger.Info(emailMeta)
	t.Logger.Info()

	api := contract.GetWorkerArgs().EmailBaseUrl
	key := contract.GetWorkerArgs().EmailApiKey

	t.Logger.Info(api)
	t.Logger.Info(key)

	mg := mailgun.NewMailgun(api, key)
	mg.SetAPIBase("https://api.eu.mailgun.net/v3")

	EmailList := []string{emailMeta.DefaultEmail}

	if emailMeta.Email != nil && *emailMeta.Email != "" {
		fmt.Println("emailMeta.Email is not null")
		EmailList = append(EmailList, *emailMeta.Email)
	}

	t.Logger.Info(EmailList)

	var dataObj map[string]string
	err = json.Unmarshal(data, &dataObj)

	if err != nil {
		t.Logger.Errorw("Error while decoding data []byte to map[string]string", err)
		return err
	}

	body := template.IngestDataIntoMsgBody(emailMeta.FirstName, eventType, emailMeta.MessageTemplate, dataObj)

	for _, email := range EmailList {
		subject := emailMeta.Subject
		sender := "Traders Connect noreply@mg.tradersconnect.com"
		recipient := email

		template := GetHtmlTemplate(subject, body)
		message := mg.NewMessage(sender, subject, "")

		message.SetHtml(template)

		message.AddRecipient(recipient)
		_, _, err = mg.Send(context.Background(), message)
		if err != nil {
			fmt.Println("Error sending email:", err)
		}

		reqMeta := struct {
			Subject string
			Message string
			Email   string
		}{
			Subject: sender,
			Message: body,
			Email:   email,
		}

		reqMetaBytes, _ := json.Marshal(reqMeta)

		t.Db.DumpLog(model.Logs{
			UserConfig:       userConfig,
			AccountId:        accId,
			EventType:        eventType,
			NotificationType: contract.EMAIL,
			ReqMeta:          reqMetaBytes,
			Status:           H.If(err != nil, "FAILED", "SUCCESS"),
		})
	}

	t.Logger.Info("Email sent successfully!")
	return err
}

func GetHtmlTemplate(subject, body string) string {
	htmlBody := template.Template

	htmlBody = strings.ReplaceAll(htmlBody, "{{Subject}}", subject)
	htmlBody = strings.ReplaceAll(htmlBody, "{{Body}}", body)

	return htmlBody
}
