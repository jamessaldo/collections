run:
	swag init && go run .
build:
	go build
install:
	go mod tidy
access:
	go run . seed AccessSeed
gotest:
	go test ./test/... -cover -v