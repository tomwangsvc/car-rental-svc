package schema

const (
	Car        = "car.json"
	CarCreate  = "car_create.json"
	Cars       = "cars.json"
	CarsSearch = "cars_search.json"
	CarUpdate  = "car_update.json"
)

func SupportedSchema() []string {
	return []string{
		Car,
		CarCreate,
		CarsSearch,
		Cars,
		CarUpdate,
	}
}
