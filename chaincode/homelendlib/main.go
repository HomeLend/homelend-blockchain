package lib

// HomelendChaincode basic struct to provide an API
type HomelendChaincode struct {
}

// Property describes structure of real estate
type Property struct {
	Hash         string `json:"Hash"`
	Address      string `json:"Address"`
	SellingPrice int    `json:"SellingPrice"`
	Timestamp    int    `json:"Timestamp"`
}

// Request defines buy processing and contains
type Request struct {
	Hash         string `json:"Hash"`
	PropertyHash string `json:"Name"`
	BuyerHash    string `json:"Buyer"`
	SellerHash   string `json:"Seller"`
	CreditScore  string `json:"CreditScore"`
	Salary       int    `json:"TotalSupply"`
	LoanAmount   int    `json:"LoanAmount"`
	Status       string `json:"Status,omitempty"`
	Timestamp    int    `json:"Timestamp"`
}

// Bank describes fields of Bank
type Bank struct {
	Hash        string `json:"Hash"`
	Name        string `json:"Name"`
	TotalSupply int    `json:"TotalSupply"`
	Timestamp   int    `json:"Timestamp"`
}

// Seller structure describes the seller fields
type Seller struct {
	ID        string `json:"ID"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Timestamp int    `json:"Timestamp"`
}

// Buyer describes fields necessary for buyer
type Buyer struct {
	ID           string `json:"ID"`
	Firstname    string `json:"Firstname"`
	Lastname     string `json:"Lastname"`
	IDNumber     string `json:"IDNumber"`
	IDBase64     string `json:"IDBase64"`
	SalaryBase64 string `json:"SalaryBase64"`
	Timestamp    int    `json:"Timestamp"`
}

// Appraiser describes fields necessary for appraiser
type Appraiser struct {
	ID        string `json:"ID"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Timestamp int    `json:"Timestamp"`
}

// InsuranceCompany describes fields necessary for insurance company
type InsuranceCompany struct {
	ID        string `json:"ID"`
	Name      string `json:"Firstname"`
	Address   string `json:"Lastname"`
	Timestamp int    `json:"Timestamp"`
}
