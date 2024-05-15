package citiesreader

import (
	"encoding/json"
	"os"

	"github.com/wojcikp/go-weather-go/weather-feed/internal"
)

type IReader interface {
	Read() ([]CityInput, error)
}

type CitiesReader struct {
	FilePath string
}

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

func NewReader(path string) *CitiesReader {
	return &CitiesReader{}
}

func (r CitiesReader) Read() ([]CityInput, error) {
	file, err := os.ReadFile(r.FilePath)
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

func GetCitiesInput(r IReader) ([]internal.BaseCityInfo, error) {
	input, err := r.Read()
	if err != nil {
		return nil, err
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
