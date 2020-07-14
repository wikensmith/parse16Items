package structs

type DETRStruct struct {
	UsedHistoryPrice float64
	UsedFare     float64            `json:"UsedFare"`
	Taxs         map[string]float64 `json:"Taxs"`
	Data         *Data_             `json:"Data"`
	CostInfo     *CostInfo          `json:"CostInfo"`
	Error        int                `json:"Error"`
	Message      string             `json:"Message"`
	EtermStr     []string           `json:"EtermStr"`
	EtermTraffic int                `json:"EtermTraffic"`
}
type TripInfos_ struct {
	TripCode       string      `json:"TripCode"`
	TripNo         string      `json:"TripNo"`
	Airline        string      `json:"Airline"`
	FlightNo       string      `json:"FlightNo"`
	Cabin          string      `json:"Cabin"`
	CabinMark      interface{} `json:"CabinMark"`
	FromCity       string      `json:"FromCity"`
	FromAirport    string      `json:"FromAirport"`
	ToCity         string      `json:"ToCity"`
	ToAirport      string      `json:"ToAirport"`
	FormTerminal   string      `json:"FormTerminal"`
	ToTerminal     string      `json:"ToTerminal"`
	TicketNoStatus string      `json:"TicketNoStatus"`
	FlightDate     string      `json:"FlightDate"`
	FareBasis      string      `json:"FareBasis"`
	DepartureTime  string      `json:"DepartureTime"`
	ArrivalTime    interface{} `json:"ArrivalTime"`
	Luggage        string      `json:"Luggage"`
	TicketNo       string      `json:"TicketNo"`
	Pnr            string      `json:"Pnr"`
}
type Data_ struct {
	PassengerName string        `json:"PassengerName"`
	Endorsement   string        `json:"Endorsement"`
	OldTicketNo   string        `json:"OldTicketNo"`
	TripInfos     []*TripInfos_ `json:"TripInfos"`
	Itinerary     interface{}   `json:"itinerary"`
}
type TripList struct {
	Airline     string `json:"Airline"`
	FromAirport string `json:"FromAirport"`
	ToAirport   string `json:"ToAirport"`
	Share       bool   `json:"Share"`
	FlightNo    string `json:"FlightNo"`
	Cabin       string `json:"Cabin"`
	FlyDate     string `json:"FlyDate"`
}
type TripPriceList struct {
	Airline    string      `json:"Airline"`
	FromCity   string      `json:"FromCity"`
	ToCity     string      `json:"ToCity"`
	Share      bool        `json:"Share"`
	FlightNo   string      `json:"FlightNo"`
	Cabin      string      `json:"Cabin"`
	FlyDate    string      `json:"FlyDate"`
	QValue     float64     `json:"QValue"`
	SValue     float64     `json:"SValue"`
	OtherValue interface{} `json:"OtherValue"`
	Value      float64     `json:"Value"`
	Mileage    int         `json:"Mileage"`
}
type CostInfo struct {
	Currency      string             `json:"Currency"`
	ROEValue      float64            `json:"ROEValue"`
	NUCValue      float64            `json:"NUCValue"`
	CNFee         float64            `json:"CNFee"`
	YQFee         float64            `json:"YQFee"`
	YRFee         float64            `json:"YRFee"`
	Taxs          map[string]float64 `json:"Taxs"`
	Asd           string             `json:"asd"`
	Price         float64            `json:"Price"`
	Tax           float64            `json:"Tax"`
	AgencyFee     float64            `json:"AgencyFee"`
	EXCH          interface{}        `json:"EXCH"`
	CONJTKT       string             `json:"CONJTKT"`
	IssueDate     string             `json:"IssueDate"`
	Pnr           string             `json:"Pnr"`
	TripList      []*TripList        `json:"TripList"`
	TripPriceList []*TripPriceList   `json:"TripPriceList"`
}

