package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	handleRequests()
}

func getBuses(w http.ResponseWriter, r *http.Request) {

	fmt.Print(http.Dir("."))

	busesX34 := RequestXml("X34")
	busesX35 := RequestXml("X35")
	busesX2 := RequestXml("X2")
	busesX32 := RequestXml("X32")

	busesX34 = append(busesX34, busesX35...)
	busesX34 = append(busesX34, busesX2...)
	busesX34 = append(busesX34, busesX32...)

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(busesX34)
}

func handleRequests() {
	mux := http.NewServeMux()
	mux.HandleFunc("/buses", getBuses)
	mux.Handle("/", http.FileServer(http.Dir("./")))

	http.ListenAndServe(":8080", mux)
}

type Bus struct {
	ServiceName                 string  `json:"bus_service_name"`
	Latitude                    float64 `json:"latitude"`
	Longitude                   float64 `json:"longitude"`
	OriginAimedDepartureTime    string  `json:"originAimedDepartureTime"`
	DestinationAimedArrivalTime string  `json:"destinationAimedArrivalTime"`
	RecordedAtTime              string  `json:"recordedAtTime"`
	OriginName                  string  `json:"originName"`
	DestinationName             string  `json:"destinationName"`
	VehicleRef                  string  `json:"vehicleRef"`
	InService                   bool    `json:"inService"`
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

func RequestXml(busNumber string) []Bus {
	response, err := http.Get(fmt.Sprintf("https://data.bus-data.dft.gov.uk/api/v1/datafeed?operatorRef=THTR&lineRef=%s&api_key=xxx", busNumber))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("No response from request")
	}

	var result Siri
	if err := xml.Unmarshal([]byte(body), &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	busesXml := result.ServiceDelivery.VehicleMonitoringDelivery.VehicleActivity

	buses := []Bus{}
	for _, v := range busesXml {

		latitude, _ := strconv.ParseFloat(v.MonitoredVehicleJourney.VehicleLocation.Latitude, 64)
		longitude, _ := strconv.ParseFloat(v.MonitoredVehicleJourney.VehicleLocation.Longitude, 64)

		validUntilTime, _ := time.Parse(time.RFC3339, v.MonitoredVehicleJourney.DestinationAimedArrivalTime)
		timeNow := time.Now()

		inService := validUntilTime.After(timeNow)

		bs := Bus{
			ServiceName:                 fmt.Sprint(v.MonitoredVehicleJourney.PublishedLineName),
			Latitude:                    latitude,
			Longitude:                   longitude,
			OriginAimedDepartureTime:    normilizeTime(v.MonitoredVehicleJourney.OriginAimedDepartureTime),
			DestinationAimedArrivalTime: normilizeTime(v.MonitoredVehicleJourney.DestinationAimedArrivalTime),
			RecordedAtTime:              normilizeTime(v.RecordedAtTime),
			OriginName:                  normilizeName(v.MonitoredVehicleJourney.OriginName),
			DestinationName:             normilizeName(v.MonitoredVehicleJourney.DestinationName),
			VehicleRef:                  v.MonitoredVehicleJourney.VehicleRef,
			InService:                   inService,
		}

		buses = append(buses, bs)
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

type Siri struct {
	XMLName         xml.Name `xml:"Siri"`
	Text            string   `xml:",chardata"`
	ServiceDelivery struct {
		Text                      string `xml:",chardata"`
		ResponseTimestamp         string `xml:"ResponseTimestamp"`
		ProducerRef               string `xml:"ProducerRef"`
		VehicleMonitoringDelivery struct {
			Text                  string `xml:",chardata"`
			ResponseTimestamp     string `xml:"ResponseTimestamp"`
			RequestMessageRef     string `xml:"RequestMessageRef"`
			ValidUntil            string `xml:"ValidUntil"`
			ShortestPossibleCycle string `xml:"ShortestPossibleCycle"`
			VehicleActivity       []struct {
				Text                    string `xml:",chardata"`
				RecordedAtTime          string `xml:"RecordedAtTime"`
				ItemIdentifier          string `xml:"ItemIdentifier"`
				ValidUntilTime          string `xml:"ValidUntilTime"`
				MonitoredVehicleJourney struct {
					Text                    string `xml:",chardata"`
					LineRef                 string `xml:"LineRef"`
					DirectionRef            string `xml:"DirectionRef"`
					FramedVehicleJourneyRef struct {
						Text                   string `xml:",chardata"`
						DataFrameRef           string `xml:"DataFrameRef"`
						DatedVehicleJourneyRef string `xml:"DatedVehicleJourneyRef"`
					} `xml:"FramedVehicleJourneyRef"`
					PublishedLineName           string `xml:"PublishedLineName"`
					OperatorRef                 string `xml:"OperatorRef"`
					OriginRef                   string `xml:"OriginRef"`
					OriginName                  string `xml:"OriginName"`
					DestinationRef              string `xml:"DestinationRef"`
					DestinationName             string `xml:"DestinationName"`
					OriginAimedDepartureTime    string `xml:"OriginAimedDepartureTime"`
					DestinationAimedArrivalTime string `xml:"DestinationAimedArrivalTime"`
					VehicleLocation             struct {
						Text      string `xml:",chardata"`
						Longitude string `xml:"Longitude"`
						Latitude  string `xml:"Latitude"`
					} `xml:"VehicleLocation"`
					BlockRef   string `xml:"BlockRef"`
					VehicleRef string `xml:"VehicleRef"`
					Bearing    string `xml:"Bearing"`
				} `xml:"MonitoredVehicleJourney"`
				Extensions struct {
					Text           string `xml:",chardata"`
					VehicleJourney struct {
						Text        string `xml:",chardata"`
						Operational struct {
							Text          string `xml:",chardata"`
							TicketMachine struct {
								Text                     string `xml:",chardata"`
								TicketMachineServiceCode string `xml:"TicketMachineServiceCode"`
								JourneyCode              string `xml:"JourneyCode"`
							} `xml:"TicketMachine"`
						} `xml:"Operational"`
						VehicleUniqueId string `xml:"VehicleUniqueId"`
					} `xml:"VehicleJourney"`
				} `xml:"Extensions"`
			} `xml:"VehicleActivity"`
		} `xml:"VehicleMonitoringDelivery"`
	} `xml:"ServiceDelivery"`
}
