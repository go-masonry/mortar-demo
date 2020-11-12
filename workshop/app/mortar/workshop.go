package mortar

import (
	"context"

	workshop "github.com/go-masonry/mortar-demo/workshop/api"
	"github.com/go-masonry/mortar-demo/workshop/app/controllers"
	"github.com/go-masonry/mortar-demo/workshop/app/data"
	"github.com/go-masonry/mortar-demo/workshop/app/services"
	"github.com/go-masonry/mortar-demo/workshop/app/validations"
	serverInt "github.com/go-masonry/mortar/interfaces/http/server"
	"github.com/go-masonry/mortar/providers/groups"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type tutorialServiceDeps struct {
	fx.In

	// API Implementations
	Workshop    workshop.WorkshopServer
}

func TutorialAPIsAndOtherDependenciesFxOption() fx.Option {
	return fx.Options(
		// GRPC Service APIs registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCServerAPIs,
			Target: tutorialGRPCServiceAPIs,
		}),
		// GRPC Gateway Generated Handlers registration
		fx.Provide(fx.Annotated{
			Group:  groups.GRPCGatewayGeneratedHandlers + ",flatten", // "flatten" does this [][]serverInt.GRPCGatewayGeneratedHandlers -> []serverInt.GRPCGatewayGeneratedHandlers
			Target: tutorialGRPCGatewayHandlers,
		}),
		// All other tutorial dependencies
		tutorialDependencies(),
	)
}

func tutorialGRPCServiceAPIs(deps tutorialServiceDeps) serverInt.GRPCServerAPI {
	return func(srv *grpc.Server) {
		workshop.RegisterWorkshopServer(srv, deps.Workshop)
		// Any additional gRPC Implementations should be called here
	}
}

func tutorialGRPCGatewayHandlers() []serverInt.GRPCGatewayGeneratedHandlers {
	return []serverInt.GRPCGatewayGeneratedHandlers{
		// Register workshop REST API
		func(mux *runtime.ServeMux, endpoint string) error {
			return workshop.RegisterWorkshopHandlerFromEndpoint(context.Background(), mux, endpoint, []grpc.DialOption{grpc.WithInsecure()})
		},
		// Any additional gRPC gateway registrations should be called here
	}
}

func tutorialDependencies() fx.Option {
	return fx.Provide(
		services.CreateWorkshopService,
		controllers.CreateWorkshopController,
		data.CreateCarDB,
		validations.CreateWorkshopValidations,
	)
}
