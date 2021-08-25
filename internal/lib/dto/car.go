package dto

import (
	lib_pagination "github.com/tomwangsvc/lib-svc/pagination"
	lib_search "github.com/tomwangsvc/lib-svc/search"
)

type CarCreate struct {
	UserInput CarCreateUserInput
	Test      bool
}

type CarCreateUserInput struct {
	BrandName string `json:"brand_name"`
	ModelName string `json:"model_name"`
}

type CarsSearch struct {
	Filters         CarsSearchFilters
	IntegrationTest bool
	Pagination      lib_pagination.Pagination
}

type CarsSearchFilters struct {
	LinkedFilters []lib_search.LinkedFilter
	Test          bool `json:"test"`
}

type CarRead struct {
	Id                    string
	IntegrationTest, Test bool
}

type CarUpdate struct {
	Id        string
	UserInput CarUpdateUserInput
	Test      bool
}

type CarUpdateUserInput struct {
	BrandName *string `json:"brand_name,omitempty"`
	ModelName *string `json:"model_name,omitempty"`
}

type CarDelete struct {
	Id   string
	Test bool
}
