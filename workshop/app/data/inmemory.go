package data

import (
	"context"
	"fmt"

	"go.uber.org/fx"
)

// This interface will represent our car db
type CarDB interface {
	InsertCar(ctx context.Context, car *CarEntity) error
	PaintCar(ctx context.Context, carNumber string, newColor string) error
	GetCar(ctx context.Context, carNumber string) (*CarEntity, error)
	RemoveCar(ctx context.Context, carNumber string) (*CarEntity, error)
}

type carDBDeps struct {
	fx.In
}

func CreateCarDB(deps carDBDeps) CarDB {
	return &carDB{
		deps: deps,
		cars: make(map[string]*CarEntity),
	}
}

type carDB struct {
	deps carDBDeps
	cars map[string]*CarEntity
}

func (c *carDB) InsertCar(ctx context.Context, car *CarEntity) error {
	if _, exists := c.cars[car.CarNumber]; exists {
		return fmt.Errorf("car %s already exists", car.CarNumber)
	}
	c.cars[car.CarNumber] = car
	return nil
}

func (c *carDB) PaintCar(ctx context.Context, carNumber string, newColor string) error {
	if car, exists := c.cars[carNumber]; exists {
		car.CurrentColor = newColor
		car.Painted = true
		return nil
	}
	return fmt.Errorf("unknown car ID %s", carNumber)
}

func (c *carDB) GetCar(ctx context.Context, carNumber string) (*CarEntity, error) {
	if car, exists := c.cars[carNumber]; exists {
		return car, nil
	}
	return nil, fmt.Errorf("unknown car ID %s", carNumber)
}
func (c *carDB) RemoveCar(ctx context.Context, carNumber string) (*CarEntity, error) {
	car, err := c.GetCar(ctx, carNumber)
	if err == nil {
		delete(c.cars, carNumber)
		return car, nil
	}
	return nil, err
}
