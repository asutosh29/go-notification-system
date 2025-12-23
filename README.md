# Gin Server Template

Basic Http Server with GORM as Database ORM

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