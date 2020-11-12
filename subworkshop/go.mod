module github.com/go-masonry/mortar-demo/subworkshop

go 1.14

require (
	github.com/alecthomas/kong v0.2.11
	github.com/go-masonry/bjaeger v0.1.1
	github.com/go-masonry/bprometheus v0.1.1
	github.com/go-masonry/bviper v0.1.1
	github.com/go-masonry/bzerolog v0.1.1
	github.com/go-masonry/mortar v0.1.1
	github.com/go-masonry/mortar-demo/workshop v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.0.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/fx v1.13.1
	google.golang.org/genproto v0.0.0-20201007142714-5c0e72c5e71e
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/go-masonry/mortar-demo/workshop => ../workshop
