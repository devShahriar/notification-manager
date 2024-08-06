package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Traders-Connect/esb-contract/golang/notification_manager"
	"github.com/Traders-Connect/notification-manager/contract"
	"github.com/Traders-Connect/notification-manager/db"
	"github.com/Traders-Connect/notification-manager/model"
	"github.com/Traders-Connect/utils"
	"github.com/davecgh/go-spew/spew"
)

func GetTestDB() db.DB {

	logger, err := utils.NewLogger("notification-server", "info")
	if err != nil {
		logger.Fatal(err)
	}
	arg := contract.GetWorkerArgs()
	arg.DbUser = "notificationmanager"
	arg.DbPass = "password"
	arg.DbHost = "127.0.0.1:3306"
	arg.DbName = "notificationmanager"
	DBDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", arg.DbUser, arg.DbPass, arg.DbHost, arg.DbName)

	DB, err := db.NewMysql(DBDsn, logger)
	if err != nil {
		logger.Info(err)
	}
	return DB
}

func TestSetUserConfig(t *testing.T) {
	logger, err := utils.NewLogger("notification-server", "info")
	if err != nil {
		logger.Fatal(err)
	}
	arg := contract.GetWorkerArgs()
	arg.DbUser = "notificationmanager"
	arg.DbPass = "password"
	arg.DbHost = "127.0.0.1:3306"
	arg.DbName = "notificationmanager"
	DBDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", arg.DbUser, arg.DbPass, arg.DbHost, arg.DbName)

	DB, err := db.NewMysql(DBDsn, logger)
	if err != nil {
		logger.Info(err)
	}
	//result, err := DB.IsAccountNotificationDisabled(11, 18)
	//Db.SetUserConfig("1234")
	//Db.SetAccountConfig("1234", "898737", "")
	// Db.SetAccountConfig("1234", "898717", "kokod")
	// Db.SetAccountConfig("1234", "898767", "kokod")
	//Db.DeleteConfig("1234")
	// res, _ := DB.GetUserConfig("00004")
	// fmt.Println(res.AccountConfId)
	// fmt.Println(res.UserConfigId)
	// fmt.Println(res.UserId)
	// fmt.Println(err)

	// reply := notification_manager.ConfigDetailsReply{}

	// DB.PopulateBlockedAccountList(24, &reply)
	//DB.PopulateEnabledAccountList(2, 24, &reply)
	cId, err := DB.GetUserConfigId(ctx, "234")
	fmt.Println(cId)
	fmt.Println(fmt.Sprintf("%d", *cId))
}

var ctx context.Context = context.Background()

func TestGetEmailMeta(t *testing.T) {

	Db := GetTestDbConn()
	res, err := Db.GetEmailMeta(ctx, "1001", "898737", "trade_failed")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func GetTestDbConn() db.DB {
	logger, err := utils.NewLogger("notification-server", "info")
	if err != nil {
		logger.Fatal(err)
	}
	arg := contract.GetWorkerArgs()
	arg.DbUser = "root"
	arg.DbPass = "asd"
	arg.DbHost = "127.0.0.1:3306"
	arg.DbName = "nt"
	DBDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", arg.DbUser, arg.DbPass, arg.DbHost, arg.DbName)

	Db, err := db.NewMysql(DBDsn, logger)
	if err != nil {
		logger.Info(err)
	}
	return Db
}

func TestAddDefaultConfig(t *testing.T) {
	logger, err := utils.NewLogger("notification-server", "info")
	if err != nil {
		logger.Fatal(err)
	}
	arg := contract.GetWorkerArgs()
	arg.DbUser = "notificationmanager"
	arg.DbPass = "password"
	arg.DbHost = "127.0.0.1:3306"
	arg.DbName = "notificationmanager"
	DBDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", arg.DbUser, arg.DbPass, arg.DbHost, arg.DbName)

	DB, err := db.NewMysql(DBDsn, logger)
	if err != nil {
		logger.Info(err)
	}

	err = DB.AddDefaultNotificationConfig(ctx, 3, "email", []string{"ACCOUNT_ENABLED"})
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetEnabledNotificationTypes(t *testing.T) {
	logger, err := utils.NewLogger("notification-server", "info")
	if err != nil {
		logger.Fatal(err)
	}
	arg := contract.GetWorkerArgs()
	arg.DbUser = "notificationmanager"
	arg.DbPass = "password"
	arg.DbHost = "127.0.0.1:3306"
	arg.DbName = "notificationmanager"
	DBDsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", arg.DbUser, arg.DbPass, arg.DbHost, arg.DbName)

	DB, err := db.NewMysql(DBDsn, logger)
	if err != nil {
		logger.Info(err)
	}

	res, err := DB.GetEnabledNotificationTypes(ctx, "3", "ACCOUNT_ENABLED")
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(res)
}

func TestGetBotNotificationMeta(t *testing.T) {

	db := GetTestDB()

	meta, err := db.GetBotNotificationMeta(ctx, "3", "TRADE_COPY_FAILURE", "telegram")
	if err != nil {
		fmt.Println(err)
	}

	spew.Dump(meta)

}

func TestPopulateBlockChannelList(t *testing.T) {

	db := GetTestDB()
	data := notification_manager.GetBotEventDetailsReply{}
	db.PopulateUnBlockChannelList(ctx, 3, 1, 18, &data)
	spew.Dump(data.ChannelEnabledList)
	fmt.Print("s")
}

func TestAddDefaultNotificationConfig(t *testing.T) {

	db := GetTestDB()
	fmt.Println("c")
	db.AddDefaultNotificationConfig(ctx, 3, contract.TELEGRAM, []string{"ACCOUNT_ENABLED"})
}

func TestGetAccountListByUserConfig(t *testing.T) {

	db := GetTestDB()
	fmt.Println("c")
	res, err := db.GetAccountListByUserConfig(ctx, 3)
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump(res)
}

func TestBlockAccountAndChannel(t *testing.T) {

	db := GetTestDB()
	fmt.Println("c")
	err := db.BlockAccountAndChannel(ctx, 3, 22)
	if err != nil {
		fmt.Println(err)
	}
}

func TestCheckValidUser(t *testing.T) {

	db := GetTestDB()
	fmt.Println("c")
	result := db.CheckValidUser(ctx, "123", 11, &model.ChannelConfig{})

	fmt.Println(result)
	fmt.Println(result)

}

//CheckValidUser(reqUserId string, configId uint64, tableName interface{})
