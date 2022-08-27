swagger:
	swag init

local-run: 
	AUTH0_DOMAIN="dev-6yej6d9j.auth0.com" AUTH0_API_IDENTIFIER="https://api.recipes.io" MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run cmd/server/main.go

init-users:
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" MONGO_DATABASE=demo go run cmd/makeuser/main.go

apache-benchmark:
	ab -n 2000 -c 100 -g without-cache.data http://localhost:8080/recipess