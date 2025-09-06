generate:
	swag init -g main.go --parseDependency --parseInternal --exclude vendor,static,internal/config,internal/domain,internal/repository,internal/service
run:
	go run main.go