package mortar

import (
	"context"
	"github.com/go-masonry/mortar-demo/subworkshop/app/controllers"
	"github.com/go-masonry/mortar-demo/subworkshop/app/services"
	"github.com/go-masonry/mortar-demo/subworkshop/app/validations"

	"github.com/go-masonry/mortar-demo/subworkshop/api"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type subworkshopServiceDeps struct {
	fx.In

	// API Implementations
	SubWorkshop subworkshop.SubWorkshopServer
}

func SubWorkshopAPIsAndOtherDependenciesFxOption() fx.Option {
	return fx.Options(
		// GRPC Service APIs registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCServerAPIs,
			Target: subworkshopGRPCServiceAPIs,
		}),
		// GRPC Gateway Generated Handlers registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCGatewayGeneratedHandlers + ",flatten", // "flatten" does this [][]serverInt.GRPCGatewayGeneratedHandlers -> []serverInt.GRPCGatewayGeneratedHandlers
			Target: subworkshopGRPCGatewayHandlers,
		}),
		// All other tutorial dependencies
		subworkshopDependencies(),
	)
}

func subworkshopGRPCServiceAPIs(deps subworkshopServiceDeps) serverInt.GRPCServerAPI {
	return func(srv *grpc.Server) {
		subworkshop.RegisterSubWorkshopServer(srv, deps.SubWorkshop)
		// Any additional gRPC Implementations should be called here
	}
}

func subworkshopGRPCGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register sub subworkshop REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return subworkshop.RegisterSubWorkshopHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func subworkshopDependencies() fx.Option {
	return fx.Provide(
		services.CreateSubWorkshopService,
		controllers.CreateSubWorkshopController,
		validations.CreateSubWorkshopValidations,
	)
}
