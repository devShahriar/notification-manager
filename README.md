# notificationmanager-manager

A distributed notificationmanager system 



## GRPCurl

### Internal calls 

```bash
cat examples/internal/send_notification_internal.json | grpcurl -plaintext -d @ localhost:9031 notificationmanager.Internal/IntSendNotification
```

```bash
cat examples/internal/add_user_config_internal.json | grpcurl -plaintext -d @ localhost:9031 notificationmanager.Internal/IntAddUserConfig
```

```bash
cat examples/internal/add_account_config_internal.json | grpcurl -plaintext -d @ localhost:9031 notificationmanager.Internal/IntAddAccountConfig
```

```bash
cat examples/internal/delete_account_config_internal.json | grpcurl -plaintext -d @ localhost:9031 notificationmanager.Internal/IntDeleteAccountConfig
```

### External calls

```bash
cat examples/external/add_config_external.json | grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/AddConfig
```

```bash
cat examples/external/get_config_external.json | grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/GetConfig
```

```bash
cat examples/external/edit_config_external.json | grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/EditConfig
```

```bash
cat examples/external/delete_config_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/DeleteConfig
```

```bash
cat examples/external/edit_account_config_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/EditAccountConfig
```

```bash
cat examples/external/get_integration_status.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/GetIntegrationStatus
```

```bash
cat examples/external/install_integration.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/InstallIntegration
```

```bash
cat examples/external/get_config_details.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/GetConfigDetails
```

```bash
cat examples/external/edit_config_status.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/EditConfigStatus
```


## Telegram bot Rpc examples 

### Bot
```bash
cat examples/external/add_bot_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/AddBot
```

```bash
cat examples/external/edit_bot_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/EditBot
```

```bash
cat examples/external/get_bot_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/GetBots
```

```bash
cat examples/external/edit_bot_status_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/EditBotStatus
```

```bash
cat examples/external/delete_bot_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/DeleteBots
```

### Channel

```bash
cat examples/external/add_channel_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/AddChannel
```

```bash
cat examples/external/edit_channel_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/EditChannel
```

```bash
cat examples/external/get_channel_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/GetChannel
```

```bash
cat examples/external/edit_channel_status_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/EditChannelStatus
```

```bash
cat examples/external/delete_channel_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/DeleteChannel
```





```bash
cat examples/external/get_bot_event_configs_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/GetBotEventConfigs
```


```bash
cat examples/external/get_bot_event_details_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/GetBotEventDetails
```

```bash
cat examples/external/edit_bot_event_details_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/EditBotEventDetails
```

```bash
cat examples/external/edit_bot_event_status_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/EditBotEventStatus
```

```bash
cat examples/external/uninstall_integration_external.json |  grpcurl -H "authorization: Bearer $(cat examples/token.txt)" -plaintext -d @ localhost:9030 notificationmanager.NotificationManager/UninstallIntegration
```

