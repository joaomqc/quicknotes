SERVER_DIR = server

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	cd $(SERVER_DIR) && swag init --parseDependency

server: swagger
	cd $(SERVER_DIR) && go run .