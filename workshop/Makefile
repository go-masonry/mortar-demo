VER:='github.com/go-masonry/mortar/mortar.version=v1.2.3'
GIT:='github.com/go-masonry/mortar/mortar.gitCommit=$(shell git rev-parse --short HEAD)'
BUILD_TAG:='github.com/go-masonry/mortar/mortar.buildTag=42'
BUILD_TS:='github.com/go-masonry/mortar/mortar.buildTimestamp=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")'

run:
	@go run -ldflags="-X ${VER} -X ${GIT} -X ${BUILD_TAG} -X ${BUILD_TS}" main.go config config/config.yml

gen-api:
	@protoc -I.\
        -I$$GOPATH/src \
        -I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
        --go_out=:api \
		--go-grpc_out=:api \
        --grpc-gateway_out=:api \
        --openapiv2_out=:. \
        api/garage.proto\

test:
	@echo "Testing ..."
	@go test -failfast ./...

.PHONY: gen-api test run
