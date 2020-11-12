package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-masonry/mortar/interfaces/cfg"
	"net/http"

	workshop "github.com/go-masonry/mortar-demo/workshop/api"
	"github.com/go-masonry/mortar-demo/workshop/app/data"
	"github.com/go-masonry/mortar/interfaces/http/client"
	"github.com/go-masonry/mortar/interfaces/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/fx"
)

// WorkshopController responsible for the business logic of our Workshop
type WorkshopController interface {
	workshop.WorkshopServer
}

type workshopControllerDeps struct {
	fx.In

	DB                data.CarDB
	Config            cfg.Config
	Logger            log.Logger
	HTTPClientBuilder client.NewHTTPClientBuilder
}

type workshopController struct {
	deps                                 workshopControllerDeps
	client                               *http.Client
	workshop.UnimplementedWorkshopServer // if keep this one added even when you change your interface this code will compile
}

// CreateWorkshopController is a constructor for Fx
func CreateWorkshopController(deps workshopControllerDeps) WorkshopController {
	client := deps.HTTPClientBuilder().Build()
	return &workshopController{
		deps:    deps,
		client:  client,
	}
}

func (w *workshopController) AcceptCar(ctx context.Context, car *workshop.Car) (*empty.Empty, error) {
	err := w.deps.DB.InsertCar(ctx, FromProtoCarToModelCar(car))
	w.deps.Logger.WithError(err).Debug(ctx, "car accepted")
	return &empty.Empty{}, err
}

func (w *workshopController) PaintCar(ctx context.Context, request *workshop.PaintCarRequest) (*empty.Empty, error) {
	car, err := w.deps.DB.GetCar(ctx, request.GetCarNumber())
	if err != nil {
		return nil, err
	}
	httpReq, err := w.makePaintRestRequest(ctx, car, request)
	if err != nil {
		return nil, err
	}
	response, err := w.client.Do(httpReq)
	if err != nil {
		w.deps.Logger.WithError(err).Debug(ctx, "calling sub workshop failed")
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("painting failed with status %d", response.StatusCode)
	}
	return &empty.Empty{}, nil
}

func (w *workshopController) RetrieveCar(ctx context.Context, request *workshop.RetrieveCarRequest) (*workshop.Car, error) {
	car, err := w.deps.DB.GetCar(ctx, request.GetCarNumber())
	if err != nil {
		return nil, err
	}
	if car.Painted {
		car, err = w.deps.DB.RemoveCar(ctx, request.GetCarNumber())
		if err != nil {
			return nil, err
		}
		return FromModelCarToProtoCar(car), nil
	}
	return nil, fmt.Errorf("car %s is not painted", request.GetCarNumber())
}

func (w *workshopController) CarPainted(ctx context.Context, request *workshop.PaintFinishedRequest) (*empty.Empty, error) {
	err := w.deps.DB.PaintCar(ctx, request.GetCarNumber(), request.GetDesiredColor())
	return &empty.Empty{}, err
}

func (w *workshopController) makePaintRestRequest(ctx context.Context, car *data.CarEntity, request *workshop.PaintCarRequest) (httpReq *http.Request, err error) {
	subWorkshopHost := w.deps.Config.Get("workshop.services.subworkshop.host").String()
	subWorkshopRESTPort := w.deps.Config.Get("workshop.services.subworkshop.restport").String()
	ownGRPCPort := w.deps.Config.Get("mortar.server.grpc.port").String()
	ownHost := w.deps.Config.Get("workshop.services.subworkshop.callbackhost").String()
	requestBody := map[string]interface{}{
		"car":                      FromModelCarToSubWorkshopMap(car),
		"desired_color":            request.GetDesiredColor(),
		"callback_service_address": fmt.Sprintf("%s:%s", ownHost, ownGRPCPort),
	}

	body := new(bytes.Buffer)
	var marshaledBytes []byte
	if marshaledBytes, err = json.Marshal(requestBody); err == nil {
		if _, err = body.Write(marshaledBytes); err == nil {
			if httpReq, err = http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s:%s/v1/subworkshop/paint", subWorkshopHost, subWorkshopRESTPort), body); err == nil {
				httpReq = httpReq.WithContext(ctx)
			}
		}
	}

	return
}
