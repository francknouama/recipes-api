swagger:
	swag init

local-run: swagger
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run main.go

apache-benchmark:
	ab -n 2000 -c 100 -g without-cache.data http://localhost:8080/recipes