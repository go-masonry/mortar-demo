package controllers_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	workshop "github.com/go-masonry/mortar-demo/workshop/api"
	"github.com/go-masonry/mortar-demo/workshop/app/controllers"
	"github.com/go-masonry/mortar-demo/workshop/app/data"
	"github.com/go-masonry/mortar-demo/workshop/app/mortar"
	"github.com/go-masonry/mortar/http/client"
	clientInt "github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type workshopSuite struct {
	suite.Suite
	pwd        string
	ctrl       *gomock.Controller
	app        *fxtest.App
	carDB      data.CarDB
	controller controllers.WorkshopController
}

func TestWorkshop(t *testing.T) {
	suite.Run(t, new(workshopSuite))
}

func (s *workshopSuite) TestAcceptCar() {
	_, err := s.controller.AcceptCar(context.Background(), &workshop.Car{
		Number:    "1234",
		Owner:     "test owner",
		BodyStyle: workshop.Car_PHAETON,
		Color:     "cyan",
	})
	s.NoError(err)
	car, err := s.carDB.GetCar(context.Background(), "1234")
	s.NoError(err)
	s.Equal("1234", car.CarNumber)
	s.Equal("cyan", car.CurrentColor)
	s.Equal("cyan", car.OriginalColor)
	s.Equal("PHAETON", car.BodyStyle)
	s.Equal("test owner", car.Owner)
	s.False(car.Painted)
}

// TestPaintCar is special test, if you look at business logic you can see that the client tries to connect to a remote
// service. We don't have it up and running, but we fake it with the special HTTP client interceptor.
// This will allow us to simply return the response without really going anywhere.
// Check s.specialHTTPClientBuilder function to understand
//
// You are encouraged to replace this special client with a default one, look at the commented line below
// If you do you should see an error similar to this one
//	Post "http://localhost:5381/v1/subworkshop/paint": dial tcp [::1]:5381: connect: connection refused
func (s *workshopSuite) TestPaintCar() {
	_, err := s.controller.AcceptCar(context.Background(), &workshop.Car{
		Number:    "1234567",
		Owner:     "test owner",
		BodyStyle: workshop.Car_PHAETON,
		Color:     "yellow",
	})
	s.NoError(err)
	// test no such car
	_, err = s.controller.PaintCar(context.Background(), &workshop.PaintCarRequest{CarNumber: "fake car number"})
	s.Error(err)
	// now paint
	_, err = s.controller.PaintCar(context.Background(), &workshop.PaintCarRequest{
		CarNumber:    "1234567",
		DesiredColor: "orange",
	})
	s.NoError(err)
}

func (s *workshopSuite) TestRetrieveCar() {
	// No car
	_, err := s.controller.RetrieveCar(context.Background(), &workshop.RetrieveCarRequest{CarNumber: "not yet there"})
	s.EqualError(err, "unknown car ID not yet there")
	car := &data.CarEntity{
		CarNumber:     "12345",
		Owner:         "test owner",
		BodyStyle:     "SEDAN",
		OriginalColor: "white",
		CurrentColor:  "white",
		Painted:       false,
	}
	// Insert car, but it's not yet painted
	err = s.carDB.InsertCar(context.Background(), car)
	s.NoError(err)
	_, err = s.controller.RetrieveCar(context.Background(), &workshop.RetrieveCarRequest{CarNumber: "12345"})
	s.EqualError(err, "car 12345 is not painted")
	// Now paint the car and get it
	err = s.carDB.PaintCar(context.Background(), "12345", "black")
	s.NoError(err)
	carProto, err := s.controller.RetrieveCar(context.Background(), &workshop.RetrieveCarRequest{CarNumber: "12345"})
	s.NoError(err)
	s.Equal("black", carProto.GetColor())
	s.Equal("12345", carProto.GetNumber())
	s.Equal("test owner", carProto.GetOwner())
	s.Equal(workshop.Car_SEDAN, carProto.GetBodyStyle())
}

func (s *workshopSuite) TestCarPainted() {
	_, err := s.controller.AcceptCar(context.Background(), &workshop.Car{
		Number:    "123456",
		Owner:     "test owner",
		BodyStyle: workshop.Car_PHAETON,
		Color:     "fuchsia",
	})
	s.NoError(err)
	_, err = s.controller.CarPainted(context.Background(), &workshop.PaintFinishedRequest{
		CarNumber:    "123456",
		DesiredColor: "indigo",
	})
	s.NoError(err)
	car, err := s.carDB.GetCar(context.Background(), "123456")
	s.NoError(err)
	s.Equal("indigo", car.CurrentColor)
	s.Equal("fuchsia", car.OriginalColor)
	s.True(car.Painted)
}

func (s *workshopSuite) SetupSuite() {
	var err error
	s.pwd, err = os.Getwd()
	s.Require().NoError(err)
}

func (s *workshopSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.app = fxtest.New(s.T(),
		s.fxOptions(),
	)
	s.app.RequireStart()
}

func (s *workshopSuite) TearDownTest() {
	s.app.RequireStop()
	s.ctrl.Finish()
}

func (s *workshopSuite) fxOptions() fx.Option {
	return fx.Options(
		fx.NopLogger, // remove fx debug prints
		mortar.ViperFxOption(s.pwd+"/../../config/config.yml", s.pwd+"/../../config/config_test.yml"),
		mortar.LoggerFxOption(),
		//providers.HTTPClientBuildersFxOption(), // uncomment this line to see that TestPaintCar fails
		fx.Provide(s.specialHTTPClientBuilder),
		fx.Provide(data.CreateCarDB),
		fx.Provide(controllers.CreateWorkshopController),
		fx.Populate(&s.carDB),
		fx.Populate(&s.controller),
	)
}

func (s *workshopSuite) specialHTTPClientBuilder() clientInt.NewHTTPClientBuilder {
	return func() clientInt.HTTPClientBuilder {
		return client.HTTPClientBuilder().AddInterceptors(func(*http.Request, clientInt.HTTPHandler) (*http.Response, error) {
			// special case, don't go anywhere just return the response
			return &http.Response{
				Status:        "200 OK",
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				ContentLength: 11,
				Body:          ioutil.NopCloser(strings.NewReader("car painted")),
			}, nil
		})
	}
}
