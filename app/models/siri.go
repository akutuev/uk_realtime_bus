package models

import "encoding/xml"

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
