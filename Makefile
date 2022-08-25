swagger:
	swag init

local-run: swagger
	go run main.go