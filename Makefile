deps:	
	go mod init github.com/supunj/anticap	

deps-update:
	go mod tidy
	go mod vendor

build:
	go build -gcflags='-N -l' -o ./bin/anticap ./cmd/anticap/anticap.go

build-arm: deps
	env GOOS=linux GOARCH=arm go build -o ./bin/anticap_arm ./cmd/anticap/anticap.go

build-arm64: deps
	env GOOS=linux GOARCH=arm64 go build -o ./bin/anticap_arm64 ./cmd/anticap/anticap.go

docker:
	docker build -t anticap .

docker-run:
	docker run -p 8000:8000 -it anticap

clean:
	rm -rf ./bin

install:

run:
	go run cmd/anticap/anticap.go