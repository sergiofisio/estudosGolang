package models

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Address struct {
	Street       string `json:"street"`
	Number       int    `json:"number"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	ZipCode      string `json:"zip_code"`
}

type Phone struct {
	PhoneType   string `json:"phone_type"`
	CountryCode string `json:"country_code"`
	AreaCode    string `json:"area_code"`
	PhoneNumber string `json:"phone_number"`
}

type Purchase struct {
	ProductName  string  `json:"product_name"`
	Amount       float64 `json:"amount"`
	PurchaseDate string  `json:"purchase_date"`
}