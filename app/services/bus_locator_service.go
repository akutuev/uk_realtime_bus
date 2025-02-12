package services

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"uk_realtime_bus/config"
	"uk_realtime_bus/models"
)

type BusLocatorService interface {
	GetAllRealtimeBuses() []models.Bus
}

type busLocatorService struct {
	BusDataHost     string
	BusDataApiKey   string
	BusOperatorRef  string
	BusesToTrackMap map[string]string
}

func NewBusLocatorService(envSettings *config.EnvSettings) BusLocatorService {

	// refactor this
	busesToTrackMap := map[string]string{}
	for _, busName := range envSettings.BusesToTrackList {
		busesToTrackMap[busName] = ""
	}

	return &busLocatorService{
		BusDataHost:     envSettings.BusDataHost,
		BusDataApiKey:   envSettings.BusDataApiKey,
		BusOperatorRef:  envSettings.BusOperatorRef,
		BusesToTrackMap: busesToTrackMap,
	}
}

// TODO refactoring needed
func (s *busLocatorService) GetAllRealtimeBuses() []models.Bus {
	response, err := http.Get(fmt.Sprintf("https://%s/api/v1/datafeed?operatorRef=%s&api_key=%s", s.BusDataHost, s.BusOperatorRef, s.BusDataApiKey))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("No response from request")
	}

	var result models.Siri
	if err := xml.Unmarshal([]byte(body), &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	buses := []models.Bus{}
	for _, busXml := range result.ServiceDelivery.VehicleMonitoringDelivery.VehicleActivity {
		_, busToAdd := s.BusesToTrackMap[busXml.MonitoredVehicleJourney.LineRef]

		if busToAdd {
			latitude, _ := strconv.ParseFloat(busXml.MonitoredVehicleJourney.VehicleLocation.Latitude, 64)
			longitude, _ := strconv.ParseFloat(busXml.MonitoredVehicleJourney.VehicleLocation.Longitude, 64)

			validUntilTime, _ := time.Parse(time.RFC3339, busXml.MonitoredVehicleJourney.DestinationAimedArrivalTime)
			timeNow := time.Now()

			inService := validUntilTime.After(timeNow)

			bus := models.Bus{
				ServiceName:                 fmt.Sprint(busXml.MonitoredVehicleJourney.PublishedLineName),
				Latitude:                    latitude,
				Longitude:                   longitude,
				OriginAimedDepartureTime:    normilizeTime(busXml.MonitoredVehicleJourney.OriginAimedDepartureTime),
				DestinationAimedArrivalTime: normilizeTime(busXml.MonitoredVehicleJourney.DestinationAimedArrivalTime),
				RecordedAtTime:              normilizeTime(busXml.RecordedAtTime),
				OriginName:                  normilizeName(busXml.MonitoredVehicleJourney.OriginName),
				DestinationName:             normilizeName(busXml.MonitoredVehicleJourney.DestinationName),
				VehicleRef:                  busXml.MonitoredVehicleJourney.VehicleRef,
				InService:                   inService,
			}

			buses = append(buses, bus)
		}
	}

	return buses
}

func normilizeName(str string) string {
	val, ok := locationMappingList[str]
	if ok {
		return val
	}
	return str
}

func normilizeTime(str string) string {
	timeParsed, _ := time.Parse(time.RFC3339, str)

	return timeParsed.Format(time.TimeOnly)
}

var locationMappingList = map[string]string{
	"Great_Western_Park__ASDA":                           "GWP ASDA",
	"Didcot__Haydon_Road":                                "Didcot Center",
	"Didcot__Orchard_Centre":                             "Didcot Center",
	"Oxford_City_Centre__Westgate":                       "Oxford Center",
	"Didcot__Broadway":                                   "Didcot Center",
	"Newbury__Newbury_Wharf":                             "Newbury",
	"Wantage__Market_Place":                              "Wantage",
	"John_Radcliffe_Hospital__JR_Hospital_Main_Entrance": "Oxford JR hospital",
	"Didcot__Holly_Lane":                                 "GWP ASDA",
	"Grove__Mandhill_Close":                              "Wantage",
}
