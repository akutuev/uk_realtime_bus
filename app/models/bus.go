package models

// TODO add some documentation
// This bus model is used for map providing full realtime bus information
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
