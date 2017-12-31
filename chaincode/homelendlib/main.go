package lib

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

// Appraiser describes fields necessary for appraiser
type Appraiser struct {
	ID        string `json:"ID"`
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Timestamp int    `json:"Timestamp"`
}

// Property describes structure of real estate
type Property struct {
	Hash         string `json:"Hash"`
	Address      string `json:"Address"`
	SellingPrice int    `json:"SellingPrice"`
	Timestamp    int    `json:"Timestamp"`
}

// InsuranceOffer describes fields of offer
type InsuranceOffer struct {
	Hash          string `json:"Hash"`
	InsuranceHash string `json:"InsuranceHash"`
	Amount        int    `json:"Amount"`
	Timestamp     int    `json:"Timestamp"`
}

// Request defines buy processing and contains
type Request struct {
	Hash              string           `json:"Hash"`
	PropertyHash      string           `json:"PropertyHash"`
	BuyerHash         string           `json:"BuyerHash"`
	SellerHash        string           `json:"SellerHash"`
	AppraiserHash     string           `json:"AppraiserHash"`
	CreditScore       string           `json:"CreditScore"`
	AppraiserPrice    string           `json:"AppraiserPrice"`
	AppraiserAmount   int              `json:"AppraiserAmount"`
	InsuranceHash     string           `json:"InsuranceHash"`
	InsuranceAmount   string           `json:"InsuranceAmount"`
	GovernmentResult1 string           `json:"GovernmentResult1"`
	GovernmentResult2 string           `json:"GovernmentResult2"`
	GovernmentResult3 string           `json:"GovernmentResult3"`
	InsuranceOffers   []InsuranceOffer `json:"InsuranceOffers"`
	Salary            int              `json:"Salary"`
	LoanAmount        int              `json:"LoanAmount"`
	Status            string           `json:"Status"`
	Timestamp         int              `json:"Timestamp"`
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

// InsuranceCompany describes fields necessary for insurance company
type InsuranceCompany struct {
	ID        string `json:"ID"`
	Name      string `json:"Firstname"`
	Address   string `json:"Lastname"`
	Timestamp int    `json:"Timestamp"`
}
