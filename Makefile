SERVER_DIR = server
CLIENT_DIR = client

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	cd $(SERVER_DIR) && swag init --parseDependency

.PHONY: server
server: swagger
	cd $(SERVER_DIR) && go run .

.PHONY: client
client:
	cd $(CLIENT_DIR) && npm install && npm run format && npm run dev
