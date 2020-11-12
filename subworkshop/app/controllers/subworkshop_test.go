package controllers_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/go-masonry/mortar/interfaces/http/client"
	mock_client "github.com/go-masonry/mortar/interfaces/http/client/mock"
	workshop "github.com/go-masonry/mortar-demo/subworkshop/api"
	"github.com/go-masonry/mortar-demo/subworkshop/app/controllers"
	"github.com/go-masonry/mortar-demo/subworkshop/app/mortar"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"google.golang.org/grpc"
)

type subWorkshopSuite struct {
	suite.Suite
	pwd                 string
	ctrl                *gomock.Controller
	app                 *fxtest.App
	grpcConnBuilderMock *mock_client.MockGRPCClientConnectionBuilder
	gRPCWrapperMock     *mock_client.MockGRPCClientConnectionWrapper
	subController       controllers.SubWorkshopController
}

func TestSubWorkshop(t *testing.T) {
	suite.Run(t, new(subWorkshopSuite))
}

func (s *subWorkshopSuite) TestPaintCar() {
	// Prepare mocks
	fakeConnection := new(fakeGRPCConnection)
	s.gRPCWrapperMock.EXPECT().Dial(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeConnection, nil)
	s.grpcConnBuilderMock.EXPECT().Build().Return(s.gRPCWrapperMock)

	_, err := s.subController.PaintCar(context.Background(), &workshop.SubPaintCarRequest{
		Car:                    &workshop.Car{Number: "1234"},
		DesiredColor:           "black",
		CallbackServiceAddress: "/dev/null",
	})
	s.NoError(err)
	// It's not really necessary, but for the sake of the example you get see it's called
	s.Equal(1, fakeConnection.callCounter)
}

func (s *subWorkshopSuite) TestPaintCarWithFailingDialer() {
	// Prepare mocks
	s.gRPCWrapperMock.EXPECT().Dial(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("just void"))
	s.grpcConnBuilderMock.EXPECT().Build().Return(s.gRPCWrapperMock)

	_, err := s.subController.PaintCar(context.Background(), &workshop.SubPaintCarRequest{
		Car:                    &workshop.Car{Number: "1234"},
		DesiredColor:           "black",
		CallbackServiceAddress: "/dev/null",
	})
	s.EqualError(err, "car painted but we can't callback to /dev/null, just void")
}

func (s *subWorkshopSuite) SetupSuite() {
	var err error
	s.pwd, err = os.Getwd()
	s.Require().NoError(err)
}

func (s *subWorkshopSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.grpcConnBuilderMock = mock_client.NewMockGRPCClientConnectionBuilder(s.ctrl)
	s.gRPCWrapperMock = mock_client.NewMockGRPCClientConnectionWrapper(s.ctrl)
	s.app = fxtest.New(s.T(),
		fx.NopLogger, // remove fx debug prints
		mortar.ViperFxOption(s.pwd+"/../../config/config.yml", s.pwd+"/../../config/config_test.yml"),
		mortar.LoggerFxOption(),
		fx.Provide(func() client.GRPCClientConnectionBuilder {
			return s.grpcConnBuilderMock
		}),
		fx.Provide(controllers.CreateSubWorkshopController),
		fx.Populate(&s.subController),
	)
	s.app.RequireStart()
}

func (s *subWorkshopSuite) TearDownTest() {
	s.app.RequireStop()
	s.ctrl.Finish()
}

type fakeGRPCConnection struct {
	callCounter int
}

func (f *fakeGRPCConnection) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	f.callCounter++
	return nil // everything is great
}

func (f *fakeGRPCConnection) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	panic("implement me")
}
