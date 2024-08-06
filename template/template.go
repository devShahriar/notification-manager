package template

import (
	"fmt"
	"reflect"
	"strings"

	nm "github.com/Traders-Connect/esb-contract/golang/notification_manager"
	"github.com/devShahriar/H"
)

const (
	EVENT_ACCOUNT_ADDED = "ACCOUNT_ADDED"
)

func IngestDataIntoMsgBody(firstName, eventName string, template string, data map[string]string) string {

	msgStruct := GetMsgStruct(eventName)
	H.PopulateStructFromMap(msgStruct, data)
	fmt.Println(msgStruct)
	body := IngestData(msgStruct, firstName, template)
	return body
}

func IngestData(obj interface{}, firstName, template string) string {
	objValue := reflect.ValueOf(obj).Elem()
	objType := objValue.Type()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		key := field.Tag.Get("key")

		template = strings.ReplaceAll(template, "%"+key+"%", objValue.Field(i).String())
	}

	// Add FirstName
	template = strings.ReplaceAll(template, "%FIRST_NAME%", firstName)

	//Add new line
	template = strings.ReplaceAll(template, "\n", "<br>")

	return template
}

func GetMsgStruct(eventName string) interface{} {
	switch eventName {
	case nm.EventType_ACCOUNT_ADDED.String(),
		nm.EventType_ACCOUNT_ENABLED.String(),
		nm.EventType_ACCOUNT_DISABLED.String(),
		nm.EventType_ACCOUNT_DELETED.String():
		return &AccountCommon{}

	case nm.EventType_ACCOUNT_CONNECTED.String():
		return &AccountConnected{}

	case nm.EventType_ACCOUNT_CONNECTION_ERROR.String():
		return &AccountConnectionError{}

	case nm.EventType_COPIER_CREATED.String():
		return &CopierCreated{}

	case nm.EventType_COPIER_ENABLED.String(),
		nm.EventType_COPIER_DISABLED.String(),
		nm.EventType_COPIER_MODIFIED.String(),
		nm.EventType_COPIER_DELETED.String():
		return &CopierCommon{}

	case nm.EventType_TRADE_COPIED_SUCCESSFULLY.String():
		return &TradeCopiedSuccessfully{}

	case nm.EventType_TRADE_MODIFIED_SUCCESSFULLY.String():
		return &TradeModifiedSuccessfully{}

	case nm.EventType_TRADE_COPY_FAILURE.String():
		return &TradeCopyFailure{}
	}

	return nil
}
