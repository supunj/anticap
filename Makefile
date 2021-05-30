deps:
	go get -u github.com/gorilla/mux
	go get -u github.com/go-redis/redis
	go get -u github.com/spf13/viper
	go get -u github.com/gorilla/websocket

	#API Documentation
	go get -u github.com/swaggo/swag/cmd/swag

	#Test
	go get github.com/erggo/datafiller

build:
	go build -gcflags='-N -l' -o ./bin/anticap ./cmd/anticap/anticap.go

build-arm: deps
	env GOOS=linux GOARCH=arm go build -o ./bin/anticap_arm ./cmd/anticap/anticap.go

build-arm64: deps
	env GOOS=linux GOARCH=arm64 go build -o ./bin/anticap_arm64 ./cmd/anticap/anticap.go

clean:
	rm -rf ./bin

install:

run:
	go run server.go