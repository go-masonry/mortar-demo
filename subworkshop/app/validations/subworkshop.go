package validations

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-masonry/mortar-demo/subworkshop/api"
)

type SubWorkshopValidations interface {
	PaintCar(ctx context.Context, request *subworkshop.SubPaintCarRequest) error
}

type subWorkshopValidations struct{}

func CreateSubWorkshopValidations() SubWorkshopValidations {
	return new(subWorkshopValidations)
}
func (s subWorkshopValidations) PaintCar(ctx context.Context, request *subworkshop.SubPaintCarRequest) error {
	if len(request.GetCallbackServiceAddress()) == 0 {
		return status.Errorf(codes.InvalidArgument, "callback service address cannot be empty")
	}
	if request.GetCar() == nil {
		return status.Errorf(codes.InvalidArgument, "car can't be empty")
	}
	return nil
}
