package citiesreader

import (
	"reflect"
	"testing"

	"github.com/wojcikp/go-weather-go/weather-feed/internal"
)

func Test_GetCitiesInput(t *testing.T) {
	mockReader := NewReaderMock()
	want := []internal.BaseCityInfo{
		{Name: "Warsaw", Latitude: "52.2300", Longtitude: "21.0111"},
		{Name: "Kraków", Latitude: "50.0614", Longtitude: "19.9372"},
		{Name: "Łódź", Latitude: "51.7769", Longtitude: "19.4547"},
		{Name: "Wrocław", Latitude: "51.1100", Longtitude: "17.0325"},
		{Name: "Poznań", Latitude: "52.4083", Longtitude: "16.9336"},
	}
	t.Run("test GetCitiesInput()", func(t *testing.T) {
		got, err := GetCitiesInput(mockReader)
		if err != nil {
			t.Errorf("error occured: %v", err)
		}
		equal := reflect.DeepEqual(got, want)
		if !equal {
			t.Errorf("%v\n!=\n%v", got, want)
		}
	})
}
