package models

type User struct {
	ID                int      `json:"id" gorm:"primaryKey"`
	Name              string   `json:"name" gorm:"not null" validate:"required"`
	Email             string   `json:"email" gorm:"not null;unique" validate:"required"`
	Username          string   `json:"username" gorm:"not null;unique" validate:"required"`
	Password          string   `json:"-" gorm:"not null" validate:"required"`
	LostPasswordToken *string  `json:"lostPasswordToken,omitempty"`
	Clients           []Client `json:"Clients" gorm:"foreignKey:UserID"`
}

type Client struct {
	ID        int     `json:"id" gorm:"primaryKey"`
	Name      string  `json:"name" gorm:"not null"`
	Document  string  `json:"document" gorm:"not null;unique"`
	Email     string  `json:"email" gorm:"not null;unique"`
	PhoneID   int     `json:"phoneId"`
	Phone     Phone   `json:"phone" gorm:"foreignKey:PhoneID"`
	AddressID int     `json:"addressId"`
	Address   Address `json:"address" gorm:"foreignKey:AddressID"`
	UserID    int     `json:"userId"`
	User      User    `json:"user" gorm:"foreignKey:UserID"`
}

type Phone struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	CountryCode string `json:"countryCode"`
	AreaCode    string `json:"areaCode"`
	Number      string `json:"number"`
}

type Address struct {
	ID         int     `json:"id" gorm:"primaryKey"`
	Street     string  `json:"street"`
	Number     string  `json:"number"`
	Complement *string `json:"complement,omitempty"`
	District   string  `json:"district"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	Country    string  `json:"country"`
	Zipcode    string  `json:"zipcode"`
}

type Type struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Control struct {
	ID         int     `json:"id" gorm:"primaryKey"`
	TypeID     int     `json:"typeId"`
	Type       Type    `json:"type" gorm:"foreignKey:TypeID"`
	ClientID   int     `json:"clientId"`
	Client     Client  `json:"client" gorm:"foreignKey:ClientID"`
	Dates      *string `json:"dates,omitempty"`
	PackNumber int     `json:"packNumber"`
	PackValue  float64 `json:"packValue"`
	TotalValue float64 `json:"totalValue"`
	Payed      float64 `json:"payed"`
}