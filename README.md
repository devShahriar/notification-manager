# notificationmanager-manager

A distributed notificationmanager system

# Scalable Notification Service for Traders Connect

A highly scalable, multi-platform notification service built with gRPC, GORM, MySQL, and RabbitMQ. This service was custom-developed for and is currently in use by [Traders Connect](https://tradersconnect.com/).

## Overview

This notification service is designed to handle large-scale, real-time notifications across multiple platforms including Discord, Slack, Telegram, and Email. It employs a master-slave architecture for efficient message routing and delivery, meeting the specific needs of Traders Connect's forex trading alert system.

## About Traders Connect

[Traders Connect](https://tradersconnect.com/) is a leading platform in the forex trading industry. They utilize this notification service to deliver critical forex trading alerts to their users across various communication channels. The service's ability to handle high volumes of time-sensitive notifications is crucial for providing Traders Connect's users with valuable, real-time data on trade outcomes and status updates.

Key benefits for Traders Connect:

- Timely delivery of trading alerts
- Multi-channel support to reach users on their preferred platforms
- Scalability to handle increasing user base and trading volumes
- Customizable templates for different types of trading notifications
- Reliable delivery ensuring no critical alerts are missed

## Features

- **Multi-Platform Support**: Send notifications to Discord, Slack, Telegram, and Email.
- **Customizable Templates**: Users can create and manage message templates.
- **User Preferences**: Notifications are sent based on individual user settings.
- **Scalable Architecture**: Master-slave setup for high throughput and reliability.
- **gRPC API**: Efficient, low-latency communication for pushing notifications.
- **Event-Driven**: Each type of event has its own gRPC API for pushing notifications.

## Architecture

### Master Node

- Decides message routing based on user settings and message templates.
- Pushes messages to appropriate RabbitMQ queues.

### Slave Nodes (Workers)

- Each worker is responsible for a specific platform (e.g., Discord, Slack).
- Listens to dedicated RabbitMQ queues.
- Processes and sends messages to the respective platforms.

## Tech Stack

- **gRPC**: For efficient API communication.
- **GORM**: ORM library for database operations.
- **MySQL**: Primary database for storing user preferences and message templates.
- **RabbitMQ**: Message broker for distributing notifications to workers.

## What It Does

Our Scalable Notification Service is a robust, flexible system designed to handle complex notification scenarios across multiple platforms. Here's a detailed look at its capabilities:

### Event-Based Notification Triggering

- Provides distinct gRPC APIs for different types of events that can trigger notifications.
- Allows services to easily integrate and push notifications for specific events in their systems.

### User-Centric Notification Management

- Maintains user preferences for notification delivery.
- Allows users to choose which platforms they want to receive notifications on (Discord, Slack, Telegram, Email).
- Supports time-based preferences, allowing users to specify when they want to receive notifications.

### Customizable Message Templates

- Offers a template management system where users or administrators can create and edit message templates.
- Supports dynamic content insertion into templates, allowing for personalized messages.
- Enables different templates for different types of notifications and platforms.

### Intelligent Message Routing

- The master node analyzes each notification request, considering:
  - The type of event
  - User preferences
  - Message template
  - Current time and user's preferred notification times
- Based on this analysis, it routes the message to the appropriate platform-specific queue.

### Multi-Platform Delivery

- Dedicated workers for each supported platform (Discord, Slack, Telegram, Email) ensure efficient, platform-optimized delivery.
- Each worker is specialized in the API and requirements of its specific platform.

### Scalable Message Processing

- Utilizes RabbitMQ for reliable message queuing, ensuring no notification is lost even during high load.
- Workers can be scaled independently based on the load for each platform.

### Real-Time and Batch Processing

- Supports both real-time notification delivery for urgent messages and batch processing for less time-sensitive notifications.
- Batch processing can be configured to optimize API usage for platforms with rate limits.

### Delivery Tracking and Retry Mechanism

- Tracks the status of each notification (sent, failed, pending).
- Implements a retry mechanism for failed notifications with configurable retry attempts and intervals.

### Analytics and Reporting

- Provides insights into notification patterns, delivery rates, and user engagement.
- Offers dashboards for monitoring system performance and notification trends.

### Extensibility

- Designed with pluggable architecture, allowing easy addition of new notification platforms.
- Supports custom logic injection for specialized notification handling.

### Security and Compliance

- Implements encryption for sensitive data in transit and at rest.
- Provides mechanisms to comply with data protection regulations (e.g., GDPR) including data retention policies and user data export.

## Setup and Installation [command]

### installation

``go build -o nt

```
### Server
```

./nt worker server -n nt-master -t server -b amqp://guest:guest@127.0.0.1:5672/ -q nt-master -e nt-master --exchange-type direct -k nt-master --redis-backend 127.0.0.1:6379 -c 1 \
--db-user root \
--db-pass asd \
--db-host 127.0.0.1:3306 \
--db-name nt \
--addr 0.0.0.0:9000 \
--port 0.0.0.0:9001 --metrics-addr 0.0.0.0:8000

```

### Master
```

./nt worker master -n nt-master -t master -b amqp://guest:guest@127.0.0.1:5672/ -q nt-master -e nt-master --exchange-type direct -k nt-master --redis-backend 127.0.0.1:6379 -c 1 \
--db-user root \
--db-pass asd \
--db-host 127.0.0.1:3306 \
--db-name nt

````

### Slave
./nt worker slave -n nt-email -t email -b amqp://guest:guest@127.0.0.1:5672/ -q nt-email -e nt-email --exchange-type direct -k nt-email --redis-backend 127.0.0.1:6379 -c 1 \
--db-user root \
--db-pass asd \
--db-host 127.0.0.1:3306 \
--db-name nt

## API Documentation

## GRPCurl

### Internal calls

```bash
cat examples/internal/send_notification_internal.json | grpcurl -plaintext -d @ localhost:9031 notificationmanager.Internal/IntSendNotification
````

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
