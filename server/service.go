package server

import (
	"context"
	"sync"

	"github.com/RichardKnop/machinery/v2"
	nm "github.com/Traders-Connect/esb-contract/golang/notification_manager"
	"github.com/Traders-Connect/utils"
	utilGrpc "github.com/Traders-Connect/utils/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/devshahriar/notification-manager/contract"
	"github.com/devshahriar/notification-manager/db"
	"google.golang.org/grpc/reflection"
)

type NotificationService struct {
	*utils.InstrumentedServer

	nm.UnimplementedInternalServer
	nm.UnimplementedNotificationManagerServer

	Db              db.DB
	MachinaryServer *machinery.Server
	Logger          *zap.SugaredLogger
	Args            *contract.ServiceArgs
}

func GetEndpointsRules() utilGrpc.RPCRules {
	servicePath := "/notificationmanager.NotificationManager/"
	return utilGrpc.RPCRules{
		Rules: map[string]utilGrpc.RPCRule{
			//NoAuthRequired: true
			servicePath + "GetIntegrationStatus": {AllowedPermissions: []string{"nt-config:getIntegrationStatus"}, NoAuthRequired: true},
			servicePath + "InstallIntegration":   {AllowedPermissions: []string{"nt-config:installIntegration"}, NoAuthRequired: true},
			servicePath + "AddConfig":            {AllowedPermissions: []string{"nt-config:addConfig"}, NoAuthRequired: true},
			servicePath + "GetConfig":            {AllowedPermissions: []string{"nt-config:getConfig"}, NoAuthRequired: true},
			servicePath + "EditConfig":           {AllowedPermissions: []string{"nt-config:editConfig"}, NoAuthRequired: true},
			servicePath + "DeleteConfig":         {AllowedPermissions: []string{"nt-config:deleteConfig"}, NoAuthRequired: true},
			servicePath + "EditAccountMeta":      {AllowedPermissions: []string{"nt-config:editAccountMeta"}, NoAuthRequired: true},
			servicePath + "GetConfigDetails":     {AllowedPermissions: []string{"nt-config:getConfigDetails"}, NoAuthRequired: true},
			servicePath + "EditConfigStatus":     {AllowedPermissions: []string{"nt-config:editConfigStatus"}, NoAuthRequired: true},

			servicePath + "AddBot":        {AllowedPermissions: []string{"nt-config:addBot"}, NoAuthRequired: true},
			servicePath + "EditBot":       {AllowedPermissions: []string{"nt-config:editBot"}, NoAuthRequired: true},
			servicePath + "GetBots":       {AllowedPermissions: []string{"nt-config:getBots"}, NoAuthRequired: true},
			servicePath + "EditBotStatus": {AllowedPermissions: []string{"nt-config:editBotStatus"}, NoAuthRequired: true},
			servicePath + "DeleteBots":    {AllowedPermissions: []string{"nt-config:deleteBots"}, NoAuthRequired: true},

			servicePath + "AddChannel":        {AllowedPermissions: []string{"nt-config:addChannel"}, NoAuthRequired: true},
			servicePath + "EditChannel":       {AllowedPermissions: []string{"nt-config:editChannel"}, NoAuthRequired: true},
			servicePath + "GetChannel":        {AllowedPermissions: []string{"nt-config:getChannel"}, NoAuthRequired: true},
			servicePath + "EditChannelStatus": {AllowedPermissions: []string{"nt-config:editChannelStatus"}, NoAuthRequired: true},
			servicePath + "DeleteChannel":     {AllowedPermissions: []string{"nt-config:deleteChannel"}, NoAuthRequired: true},

			servicePath + "GetBotEventConfigs":  {AllowedPermissions: []string{"nt-config:getBotEventConfigs"}, NoAuthRequired: true},
			servicePath + "GetBotEventDetails":  {AllowedPermissions: []string{"nt-config:getBotEventDetails"}, NoAuthRequired: true},
			servicePath + "EditBotEventDetails": {AllowedPermissions: []string{"nt-config:editBotEventDetails"}, NoAuthRequired: true},
			servicePath + "EditBotEventStatus":  {AllowedPermissions: []string{"nt-config:editBotEventStatus"}, NoAuthRequired: true},

			servicePath + "UninstallIntegration": {AllowedPermissions: []string{"nt-config:uninstallIntegration"}, NoAuthRequired: true},
		},
	}

}

func NewNotificationService(db db.DB, machinaryServer *machinery.Server, logger *zap.SugaredLogger, args *contract.ServiceArgs) *NotificationService {
	conf := utils.Config{
		ServiceName:      "notification_manager",
		APIAddr:          args.Addr,
		APIAddrInt:       args.AddrInt,
		Logger:           logger,
		MetricsAddr:      args.MetricsAddr,
		AuthDomain:       "auth.tradersconnect.com",
		RpcEnpointRules:  GetEndpointsRules(),
		TokenAuthEnabled: true,
	}

	is, err := utils.NewInstrumentedServer(conf)
	if err != nil {
		logger.Errorw("error while creating instrumented server", "error", err)
		return nil
	}
	return &NotificationService{
		InstrumentedServer: is,
		Db:                 db,
		MachinaryServer:    machinaryServer,
		Logger:             logger,
		Args:               args,
	}
}

func (n *NotificationService) Run(ctx context.Context) {
	wg := new(sync.WaitGroup)

	nm.RegisterInternalServer(n.InstrumentedServer.GRPCServerInternal, n)
	reflection.Register(n.InstrumentedServer.GRPCServerInternal)
	grpc_health_v1.RegisterHealthServer(n.InstrumentedServer.GRPCServerInternal, n)

	nm.RegisterNotificationManagerServer(n.InstrumentedServer.GRPCServer, n)
	reflection.Register(n.InstrumentedServer.GRPCServer)
	grpc_health_v1.RegisterHealthServer(n.InstrumentedServer.GRPCServer, n)

	n.InstrumentedServer.Run(ctx, wg)
	n.Logger.Info("Notification Service running")
	n.Logger.Info(n.Args.AddrInt)
	wg.Wait()
}
