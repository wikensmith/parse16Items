package structs


type ResGetRefRules16 struct {
	Freight      Freight  `json:"Freight"`
	Error        int      `json:"Error"`
	Message      string   `json:"Message"`
	EtermStr     []string `json:"EtermStr"`
	EtermTraffic int      `json:"EtermTraffic"`
}

type Freight struct {
	SerialNo      string      `json:"SerialNo"`
	PassengerType int         `json:"PassengerType"`
	CabinMark     string      `json:"CabinMark"`
	Fare          float64     `json:"Fare"`
	Tax           float64     `json:"Tax"`
	Taxs          Taxs        `json:"Taxs"`
	Total         float64     `json:"Total"`
	Currency      string      `json:"Currency"`
	ValidityDate  string      `json:"ValidityDate"`
	FareTrips     []FareTrips `json:"FareTrips"`
	Regulations   interface{} `json:"Regulations"`
}
type Taxs struct {
	CN float64 `json:"CN"`
	MO float64 `json:"MO"`
	WN float64 `json:"WN"`
	YR float64 `json:"YR"`
}
type FareTrips struct {
	FromCity     string `json:"FromCity"`
	ToCity       string `json:"ToCity"`
	CabinMark    string `json:"CabinMark"`
	Date1        string `json:"Date1"`
	Date2        string `json:"Date2"`
	LuggageValue int    `json:"LuggageValue"`
	LuggageUnit  string `json:"LuggageUnit"`
}
