swagger:
	swag init

local-run: 
	JWT_SECRET=eUbP9shywUygMx7u MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run cmd/server/main.go

init-users:
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run cmd/makeuser/main.go

apache-benchmark:
	ab -n 2000 -c 100 -g without-cache.data http://localhost:8080/recipess