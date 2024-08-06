package test

import (
	"fmt"
	"testing"

	"github.com/devshahriar/notification-manager/template"
)

func TestTemplateEngine(t *testing.T) {
	temp := "Your trading account %ACCOUNT_NUMBER% has been successfully added to Traders Connect with the name %ACCOUNT_NAME%"
	data := map[string]string{
		"ACCOUNT_NUMBER": "123432",
		"ACCOUNT_NAME":   "shahriar",
	}
	body := template.IngestDataIntoMsgBody("shahirar", template.EVENT_ACCOUNT_ADDED, temp, data)
	fmt.Println(body)
}
