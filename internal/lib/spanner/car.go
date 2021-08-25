package spanner

import (
	"car-svc/internal/lib/constants"
	"car-svc/internal/lib/dto"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_json "github.com/tomwangsvc/lib-svc/json"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_misc "github.com/tomwangsvc/lib-svc/misc"
	lib_pagination "github.com/tomwangsvc/lib-svc/pagination"
	lib_spanner "github.com/tomwangsvc/lib-svc/spanner"
	"google.golang.org/api/iterator"
)

type Car struct {
	BrandName   string           `json:"brand_name" spanner:"brand_name"`
	CarId       string           `json:"car_id" spanner:"car_id"`
	DateCreated time.Time        `json:"date_created" spanner:"date_created"`
	DateUpdated spanner.NullTime `json:"date_updated" spanner:"date_updated"`
	ModelName   string           `json:"model_name" spanner:"model_name"`
	Test        bool             `json:"test" spanner:"test"`
}

const (
	tableCar = "car"
)

var (
	CarColumns       = lib_misc.StructTaggedFieldNames(reflect.TypeOf(Car{}), "spanner")
	CarFieldMetaData = lib_json.StructFieldMetadata(reflect.TypeOf(Car{}))
)

func (c client) TransformCarToJson(ctx context.Context, car Car) ([]byte, error) {
	lib_log.Info(ctx, "Transforming", lib_log.FmtAny("car", car))

	ca, err := lib_json.GenerateJson(car, CarFieldMetaData, "")
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed generating response")
	}

	lib_log.Info(ctx, "Transformed", lib_log.FmtInt("len(ca)", len(ca)))
	return ca, nil
}

func (c client) TransformCarsToJson(ctx context.Context, cars []Car) ([]byte, error) {
	lib_log.Info(ctx, "Transforming", lib_log.FmtInt("len(cars)", len(cars)))

	if len(cars) == 0 {
		lib_log.Info(ctx, "Transformed")
		return nil, nil
	}
	var carsList []interface{}
	for _, v := range cars {
		carsList = append(carsList, v)
	}
	carsListJson, err := lib_json.GenerateJsonList(carsList, CarFieldMetaData, "")
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed generating json list")
	}

	lib_log.Info(ctx, "Transformed", lib_log.FmtInt("len(carsListJson)", len(carsListJson)))
	return carsListJson, nil
}

func (c client) CreateCar(ctx context.Context, carCreate dto.CarCreate) (string, error) {
	lib_log.Info(ctx, "Creating", lib_log.FmtAny("carCreate", carCreate), lib_log.FmtAny("c.config", c.config))

	var car Car
	if _, err := c.spannerClient.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		count, err := readCount(ctx, tx, spanner.Statement{
			SQL: fmt.Sprintf(`
				SELECT count(car_id) AS count
				FROM %s
				WHERE brand_name = @brand_name
				AND model_name = @model_name
			`,
				tableCar,
			),
			Params: map[string]interface{}{
				"brand_name": carCreate.UserInput.BrandName,
				"model_name": carCreate.UserInput.ModelName,
			}})
		if err != nil {
			return lib_errors.Wrap(err, "Failed counting cars")
		}
		if count > 0 {
			return lib_errors.NewCustom(http.StatusConflict, "Already exist")
		}

		car = newCar(carCreate)
		mutCar, err := spanner.InsertStruct(tableCar, car)
		if err != nil {
			return lib_errors.Wrap(err, "Failed creating mutCar for car")
		}

		if err := tx.BufferWrite([]*spanner.Mutation{mutCar}); err != nil {
			return lib_errors.Wrap(err, "Failed creating car")
		}

		return nil

	}); err != nil {
		return "", lib_spanner.WrapError(err, "Failed executing read write transaction")
	}

	lib_log.Info(ctx, "Created", lib_log.FmtAny("car", car))
	return car.CarId, nil
}

func newCar(carCreate dto.CarCreate) Car {
	return Car{
		BrandName:   carCreate.UserInput.BrandName,
		CarId:       uuid.New().String(),
		DateCreated: spanner.CommitTimestamp,
		ModelName:   carCreate.UserInput.ModelName,
		Test:        carCreate.Test,
	}
}

func (c client) SearchCars(ctx context.Context, carsSearch dto.CarsSearch) ([]Car, *lib_pagination.Pagination, error) {
	lib_log.Info(ctx, "Searching", lib_log.FmtAny("carsSearch", carsSearch))

	sqlFilters, params, err := lib_spanner.GenerateSqlWhereAndParamsForSearchV2(carsSearch.Filters.LinkedFilters)
	if err != nil {
		return nil, nil, lib_errors.Wrap(err, "Failed generating sql where and params for search")
	}
	sqlString := fmt.Sprintf(`
		SELECT %s
		FROM %s
		%s
		ORDER BY date_created %s
		LIMIT %d
		OFFSET %d
		`,
		strings.Join(CarColumns, ", "),
		tableCar,
		sqlFilters,
		carsSearch.Pagination.Order,
		carsSearch.Pagination.Limit,
		carsSearch.Pagination.Offset,
	)

	stmt := spanner.Statement{
		SQL:    sqlString,
		Params: params,
	}

	ro := c.spannerClient.ReadOnlyTransaction()
	defer ro.Close()

	iter := ro.Query(ctx, stmt)
	defer iter.Stop()

	lib_log.Info(ctx, "Reading", lib_log.FmtAny("stmt", stmt))

	var cars []Car
	for {
		row, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				break
			}
			return nil, nil, lib_errors.Wrap(err, "Failed iterating car")
		}

		var car Car
		if err := row.ToStruct(&car); err != nil {
			return nil, nil, lib_errors.Wrap(err, "Failed reading car")
		}

		cars = append(cars, car)
	}

	pagination, err := readCountForPagination(ctx, ro, carsSearch.Pagination, spanner.Statement{
		SQL: fmt.Sprintf(`
			SELECT count(car_id) AS count
			FROM %s
			%s
		`,
			tableCar,
			sqlFilters,
		),
		Params: params,
	})
	if err != nil {
		return nil, nil, lib_errors.Wrap(err, "Failed reading count for pagination")
	}
	ro.Close()

	lib_log.Info(ctx, "Read", lib_log.FmtInt("len(cars)", len(cars)), lib_log.FmtAny("pagination", pagination))
	return cars, pagination, nil
}

func readCountForPagination(ctx context.Context, r lib_spanner.Reader, pagination lib_pagination.Pagination, stmt spanner.Statement) (*lib_pagination.Pagination, error) {
	lib_log.Info(ctx, "Reading", lib_log.FmtAny("stmt", stmt))
	count, err := readCount(ctx, r, stmt)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed reading count")
	}
	pagination.Total = &count
	timeNow := time.Now().UTC()
	pagination.ReadTimestamp = &timeNow
	lib_log.Info(ctx, "Read", lib_log.FmtAny("pagination", pagination))
	return &pagination, nil
}

func readCount(ctx context.Context, r lib_spanner.Reader, stmt spanner.Statement) (int64, error) {
	lib_log.Info(ctx, "Reading", lib_log.FmtAny("stmt", stmt))
	iter := r.Query(ctx, stmt)
	defer iter.Stop()
	var count int64
	row, err := iter.Next()
	if err != nil {
		if err == iterator.Done {
			lib_log.Info(ctx, "Expected one row, got none, will return error", lib_log.FmtError(err))
			return 0, lib_errors.New("Query returned no rows")
		}
		return 0, lib_errors.Wrap(err, "Failed reading row from iter")
	}
	if err := row.ColumnByName("count", &count); err != nil {
		return 0, lib_errors.Wrap(err, "Failed unpacking count into int64")
	}
	iter.Stop()
	lib_log.Info(ctx, "Read", lib_log.FmtAny("count", count))
	return count, nil
}

func (c client) ReadCar(ctx context.Context, carRead dto.CarRead) (*Car, error) {
	lib_log.Info(ctx, "Reading", lib_log.FmtAny("carRead", carRead))

	car, err := readCar(ctx, c.spannerClient.Single(), carRead.Id)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed reading car")
	}

	lib_log.Info(ctx, "Read", lib_log.FmtAny("car", car))
	return car, nil
}

func readCar(ctx context.Context, reader lib_spanner.Reader, carId string) (*Car, error) {
	lib_log.Info(ctx, "reading", lib_log.FmtString("carId", carId))

	var car Car
	if err := lib_spanner.ReadById(ctx, reader, tableCar, CarColumns, carId, &car); err != nil {
		return nil, lib_errors.Wrap(err, "Failed reading car")
	}

	lib_log.Info(ctx, "read", lib_log.FmtAny("car", car))
	return &car, nil
}

func (c client) UpdateCar(ctx context.Context, carUpdate dto.CarUpdate) error {
	lib_log.Info(ctx, "Updating", lib_log.FmtAny("carUpdate", carUpdate))

	if _, err := c.spannerClient.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		car, err := readCar(ctx, tx, carUpdate.Id)
		if err != nil {
			return lib_errors.Wrap(err, "Failed reading car")
		}

		if car.Test != carUpdate.Test {
			return lib_errors.NewCustom(http.StatusUnprocessableEntity, constants.UnprocessableEntityAccessForbiddenByTest)
		}

		if err := tx.BufferWrite([]*spanner.Mutation{spanner.UpdateMap(tableCar, newCarUpdateMap(carUpdate))}); err != nil {
			return lib_errors.Wrap(err, "Failed creating car")
		}

		lib_log.Info(ctx, "Updated", lib_log.FmtAny("carUpdate", carUpdate))

		return nil
	}); err != nil {
		return lib_spanner.WrapError(err, "Failed executing read write transaction")
	}

	return nil
}

func newCarUpdateMap(carUpdate dto.CarUpdate) map[string]interface{} {
	carUpdateMap := map[string]interface{}{
		"car_id":       carUpdate.Id,
		"date_updated": spanner.CommitTimestamp,
	}
	if carUpdate.UserInput.BrandName != nil {
		carUpdateMap["brand_name"] = *carUpdate.UserInput.BrandName
	}
	if carUpdate.UserInput.ModelName != nil {
		carUpdateMap["model_name"] = *carUpdate.UserInput.ModelName
	}

	return carUpdateMap
}

func (c client) DeleteCar(ctx context.Context, carDelete dto.CarDelete) error {
	lib_log.Info(ctx, "Deleting", lib_log.FmtAny("carDelete", carDelete))

	if _, err := c.spannerClient.ReadWriteTransaction(ctx, func(ctx context.Context, tx *spanner.ReadWriteTransaction) error {
		car, err := readCar(ctx, tx, carDelete.Id)
		if err != nil {
			return lib_errors.Wrap(err, "Failed reading car")
		}

		if car.Test != carDelete.Test {
			return lib_errors.NewCustom(http.StatusUnprocessableEntity, constants.UnprocessableEntityAccessForbiddenByTest)
		}

		if err := tx.BufferWrite([]*spanner.Mutation{spanner.Delete(tableCar, spanner.Key{carDelete.Id})}); err != nil {
			return lib_errors.Wrap(err, "Failed deleting car")
		}

		lib_log.Info(ctx, "Deleted", lib_log.FmtAny("carDelete", carDelete))

		return nil
	}); err != nil {
		return lib_spanner.WrapError(err, "Failed executing read write transaction")
	}

	return nil
}
