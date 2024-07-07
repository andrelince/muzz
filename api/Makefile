.PHONY: swagg
swagg.fmt:
	@go run github.com/swaggo/swag/cmd/swag fmt -g rest/rest.go
swagg: swagg.fmt
	@go run github.com/swaggo/swag/cmd/swag init -g rest/rest.go