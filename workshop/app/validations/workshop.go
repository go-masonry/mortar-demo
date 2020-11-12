package validations

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"

	workshop "github.com/go-masonry/mortar-demo/workshop/api"
	"google.golang.org/grpc/status"
)

type WorkshopValidations interface {
	AcceptCar(ctx context.Context, car *workshop.Car) error
	PaintCar(ctx context.Context, request *workshop.PaintCarRequest) error
	RetrieveCar(ctx context.Context, request *workshop.RetrieveCarRequest) error
	CarPainted(ctx context.Context, request *workshop.PaintFinishedRequest) error
}

type workshopValidations struct {
}

func CreateWorkshopValidations() WorkshopValidations {
	return new(workshopValidations)
}

func (w *workshopValidations) AcceptCar(ctx context.Context, car *workshop.Car) error {
	return carIdValidation(car.GetNumber())
}

func (w *workshopValidations) PaintCar(ctx context.Context, request *workshop.PaintCarRequest) error {
	supportedColors := map[string]struct{}{"red": {}, "green": {}, "blue": {}}
	if _, supported := supportedColors[strings.ToLower(request.GetDesiredColor())]; supported {
		return nil
	}
	return status.Errorf(codes.InvalidArgument, "out of ink for %s", request.GetDesiredColor())
}

func (w *workshopValidations) RetrieveCar(ctx context.Context, request *workshop.RetrieveCarRequest) error {
	return carIdValidation(request.GetCarNumber())
}

func (w *workshopValidations) CarPainted(ctx context.Context, request *workshop.PaintFinishedRequest) error {
	return carIdValidation(request.GetCarNumber())
}

func carIdValidation(carID string) error {
	if len(carID) != 8 {
		return status.Errorf(codes.InvalidArgument, "%s should be 8 chars long", carID)
	}
	return nil
}
