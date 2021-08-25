package spanner

import (
	"context"
	"reflect"
	"time"

	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_json "github.com/tomwangsvc/lib-svc/json"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_misc "github.com/tomwangsvc/lib-svc/misc"
)

type CarCustomerAssociation struct {
	CarId           string    `json:"car_id" spanner:"car_id"`
	CustomerId      string    `json:"customer_id" spanner:"customer_id"`
	DateCreated     time.Time `json:"date_created" spanner:"date_created"`
	DateRentalEnd   time.Time `json:"date_rental_end" spanner:"date_rental_end"`
	DateRentalStart time.Time `json:"date_rental_start" spanner:"DateRentalStart"`
	DateUpdated     time.Time `json:"date_updated" spanner:"date_updated"`
	Test            bool      `json:"test" spanner:"test"`
}

var (
	CarCustomerAssociationColumns       = lib_misc.StructTaggedFieldNames(reflect.TypeOf(CarCustomerAssociation{}), "spanner")
	CarCustomerAssociationFieldMetaData = lib_json.StructFieldMetadata(reflect.TypeOf(CarCustomerAssociation{}))
)

const (
	tableCarCustomerAssociation = "car_customer_association"
)

var ()

func (c client) TransformBrandClassAssociationToJson(ctx context.Context, carCustomerAssociation CarCustomerAssociation) ([]byte, error) {
	lib_log.Info(ctx, "Transforming", lib_log.FmtAny("carCustomerAssociation", carCustomerAssociation))

	carCustomerAssociationListJson, err := lib_json.GenerateJson(carCustomerAssociation, CarCustomerAssociationFieldMetaData, "")
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed generating json list")
	}

	lib_log.Info(ctx, "Transformed", lib_log.FmtInt("len(carCustomerAssociationListJson)", len(carCustomerAssociationListJson)))
	return carCustomerAssociationListJson, nil
}

func (c client) TransformBrandClassAssociationsToJson(ctx context.Context, carCustomerAssociations []CarCustomerAssociation) ([]byte, error) {
	lib_log.Info(ctx, "Transforming", lib_log.FmtInt("len(carCustomerAssociations)", len(carCustomerAssociations)))

	if len(carCustomerAssociations) == 0 {
		lib_log.Info(ctx, "Transformed")
		return nil, nil
	}
	var carCustomerAssociationsList []interface{}
	for _, v := range carCustomerAssociationsList {
		carCustomerAssociationsList = append(carCustomerAssociationsList, v)
	}
	carCustomerAssociationsListJson, err := lib_json.GenerateJsonList(carCustomerAssociationsList, CarCustomerAssociationFieldMetaData, "")
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed generating json list")
	}

	lib_log.Info(ctx, "Transformed", lib_log.FmtInt("len(carCustomerAssociationsListJson)", len(carCustomerAssociationsListJson)))
	return carCustomerAssociationsListJson, nil
}
