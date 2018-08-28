# Heroes

TBD

```
go run cmd/heroes/main.go
docker run -d -p 6379:6379 redis:latest

docker build -t heroes:0.0.1 .
docker run -d -p 6379:6379 --name heroredis --network=local redis:latest
docker run -e APP_PORT=3001 -e DB_HOST=heroredis -e DB_PORT=6379 --network=local heroes:0.0.1
```
