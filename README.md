# Run directly
```
export MAINCHANNEL=..
export TOKEN=..
go run main.go 
```


# Docker
```
docker build -t thirteenthirtyseven .
docker run thirteenthirtyseven
```



# Docker compose

## Save to docker-compose.yml
```
version: "3.9"

services:
  thirteenthirtyseven:
    build:
      context: .
    environment:
      TOKEN: "YOUR-DISCORD-TOKEN"
      MAINCHANNEL: "MAIN-CHANNEL-ID"
      DATABASEPATH: "/app/data/thirteenthirtyseven.db"
    volumes:
      - ./database:/app/data
    restart: unless-stopped
```