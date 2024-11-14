package citiesreader

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/wojcikp/go-weather-go/weather-feed/internal"
)

type ICityReader interface {
	Read() ([]CityInput, error)
}

type CitiesReader struct{}

type CitiesReaderMock struct{}

type CityInput struct {
	City              string `json:"city,omitempty"`
	Lat               string `json:"lat,omitempty"`
	Lng               string `json:"lng,omitempty"`
	Country           string `json:"country,omitempty"`
	Iso2              string `json:"iso2,omitempty"`
	Admin_name        string `json:"admin_name,omitempty"`
	Capital           string `json:"capital,omitempty"`
	Population        string `json:"population,omitempty"`
	Population_proper string `json:"population_proper,omitempty"`
}

func NewReader() ICityReader {
	return &CitiesReader{}
}

func NewReaderMock() ICityReader {
	return &CitiesReaderMock{}
}

func (r CitiesReader) Read() ([]CityInput, error) {
	// dir, err := os.Getwd()
	// if err != nil {
	// 	return nil, err
	// }
	p := path.Join("/app", "assets", "pl172.json")
	// p := path.Join(dir, "..", "..", "assets", "pl172.json")
	file, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var input []CityInput
	err = json.Unmarshal(file, &input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (r CitiesReaderMock) Read() ([]CityInput, error) {
	return []CityInput{
		{"Warsaw", "52.2300", "21.0111", "Poland", "PL", "Mazowieckie", "primary", "1860281", "1860281"},
		{"Kraków", "50.0614", "19.9372", "Poland", "PL", "Małopolskie", "admin", "800653", "800653"},
		{"Łódź", "51.7769", "19.4547", "Poland", "PL", "Łódzkie", "admin", "690422", "670642"},
		{"Wrocław", "51.1100", "17.0325", "Poland", "PL", "Dolnośląskie", "admin", "672929", "672929"},
		{"Poznań", "52.4083", "16.9336", "Poland", "PL", "Wielkopolskie", "admin", "546859", "546859"},
	}, nil
}

func GetCitiesInput(r ICityReader) ([]internal.BaseCityInfo, error) {
	input, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read file with input cities list, err: %w", err)
	}

	var baseCityInfo []internal.BaseCityInfo
	for _, city := range input {
		baseCityInfo = append(baseCityInfo, internal.BaseCityInfo{
			Name:       city.City,
			Latitude:   city.Lat,
			Longtitude: city.Lng,
		})
	}

	return baseCityInfo, nil
}
