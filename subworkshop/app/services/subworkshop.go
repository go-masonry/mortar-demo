package services

import (
	"context"

	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/go-masonry/mortar-demo/subworkshop/api"
	"github.com/go-masonry/mortar-demo/subworkshop/app/controllers"
	"github.com/go-masonry/mortar-demo/subworkshop/app/validations"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
)

type subWorkshopServiceDeps struct {
	fx.In

	Logger      log.Logger
	Controller  controllers.SubWorkshopController
	Validations validations.SubWorkshopValidations
}

type subWorkshopImpl struct {
	deps subWorkshopServiceDeps
	subworkshop.UnimplementedSubWorkshopServer
}

func CreateSubWorkshopService(deps subWorkshopServiceDeps) subworkshop.SubWorkshopServer {
	return &subWorkshopImpl{
		deps: deps,
	}
}

func (s *subWorkshopImpl) PaintCar(ctx context.Context, request *subworkshop.SubPaintCarRequest) (*empty.Empty, error) {
	if err := s.deps.Validations.PaintCar(ctx, request); err != nil {
		return nil, err
	}
	s.deps.Logger.Debug(ctx, "sub workshop - actually painting the car")
	return s.deps.Controller.PaintCar(ctx, request)
}
