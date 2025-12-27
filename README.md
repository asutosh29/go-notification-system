# Gin Notification Server

Basic Notification server using channels as message queue with GORM as Database ORM

# Docker
If docker is installed
```
> docker compose up --build
```

# Setup Instructions
## env
```
cp .env.sample .env
```

## Database
Make sure you have a postgres database running with credentials mentioned in ```.env```

Otherwise use with docker (make sure docker daemon is running...)

```
> docker compose up -d
```

## Server
For developement:
```
air
```

For prod:
```
go run ./cmd/main.go
```

# Usage

## Web UI (Easy Testing)

The server provides a web interface to easily test notifications without using curl commands.

### Access the UI
Open your browser and navigate to:
```
http://localhost:8080/test
```

The UI includes:
- **Notification Form** - Create new notifications with title and description
- **Notifications List** - View all notifications in real-time
- **Live Stream Display** - See notifications as they are created
- **Connection Status** - Visual indicator showing real-time connection status
- **Activity Logs** - Detailed logs of all operations

## API Endpoints

### Health Check
- **GET** `/health` - Check server health status
  ```bash
  curl http://localhost:8080/health
  ```

### Notifications

#### Get All Notifications
- **GET** `/notification` - Retrieve all notifications
  ```bash
  curl http://localhost:8080/notification
  ```

#### Get Specific Notification
- **GET** `/notification/:id` - Get a notification by ID
  ```bash
  curl http://localhost:8080/notification/{notification_id}
  ```

#### Create Notification
- **POST** `/notification` - Create a new notification
  ```bash
  curl -X POST http://localhost:8080/notification \
    -H "Content-Type: application/json" \
    -d '{
      "title": "Notification Title",
      "description": "Notification Description"
    }'
  ```

#### Delete Notification
- **DELETE** `/notification/:id` - Delete a notification
  ```bash
  curl -X DELETE http://localhost:8080/notification/{notification_id}
  ```

#### Stream Notifications (Server-Sent Events)
- **GET** `/notification/stream` - Real-time notification stream
  ```bash
  curl http://localhost:8080/notification/stream
  ```
  This endpoint uses an internal hub/channel system to broadcast new notifications to all connected clients in real-time.

### Welcome Page
- **GET** `/` - Welcome message

## Example Workflow

### Using the Web UI (Recommended for Testing)
1. Start the server (see Setup Instructions above)
2. Open your browser to `http://localhost:8080/test`
3. Fill in the notification form with a title and description
4. Click "Send Notification"
5. Watch the notification appear in the list and live stream display in real-time
6. Open the UI in another browser tab to see real-time updates across multiple clients

### Using curl Commands
1. Start the server (see Setup Instructions above)

2. Check server health:
   ```bash
   curl http://localhost:8080/health
   ```

3. Create a notification:
   ```bash
   curl -X POST http://localhost:8080/notification \
     -H "Content-Type: application/json" \
     -d '{
      "title": "Hello World",
      "description": "This is a test notification"
    }'
   ```

4. Get all notifications:
   ```bash
   curl http://localhost:8080/notification
   ```

5. Stream notifications in real-time (in one terminal):
   ```bash
   curl http://localhost:8080/notification/stream
   ```
   
   While streaming, create notifications in another terminal to see them broadcast in real-time.