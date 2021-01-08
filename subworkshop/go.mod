module github.com/go-masonry/mortar-demo/subworkshop

go 1.14

require (
	github.com/alecthomas/kong v0.2.12
	github.com/go-masonry/bjaeger v0.1.3
	github.com/go-masonry/bprometheus v0.1.3
	github.com/go-masonry/bviper v0.1.3
	github.com/go-masonry/bzerolog v0.1.3
	github.com/go-masonry/mortar v0.2.1
	github.com/go-masonry/mortar-demo/workshop v0.0.0-20201116100640-4a178a4540e1
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.1.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/stretchr/testify v1.6.1
	go.uber.org/fx v1.13.1
	google.golang.org/genproto v0.0.0-20210106152847-07624b53cd92
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)

replace github.com/go-masonry/mortar-demo/workshop => ../workshop
