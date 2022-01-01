DAEMON_VERSION:=$(shell grep -P -o -w '(?<=version[\ |]=[\ |]\").*?(?=\")' ./daemon/version.go)
CTL_VERSION:=$(shell grep -P -o -w '(?<=version[\ |]=[\ |]\").*?(?=\")' ./ctl/version.go)

export GOARCH=amd64
export GOOS=linux

all: ifman_ctl ifman_daemon

ifman_ctl:
	@echo "Build ifman ctl..."
	@mkdir ./bin -p
	@go build -o ./bin/ifman-ctl ./ctl
	@echo "Renaming..."
	@mv ./bin/ifman-ctl "./bin/ifman-ctl-${CTL_VERSION}-`sha256sum -b ./bin/ifman-ctl | cut -d ' ' -f 1`"

ifman_daemon:
	@echo "Build ifman daemon..."
	@mkdir ./bin -p
	@go build -o ./bin/ifman-daemon ./daemon
	@echo "Renaming..."
	@mv ./bin/ifman-daemon "./bin/ifman-daemon-${DAEMON_VERSION}-`sha256sum -b ./bin/ifman-daemon | cut -d ' ' -f 1`"

clean:
	@rm ./bin/ifman-* -rf

gen:
	@protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

prepare:
	@sudo yum install protobuf-compiler protobuf-devel protobuf -y
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
