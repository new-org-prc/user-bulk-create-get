package entities

type Address struct {
	ID      int    `json:"address_id" gorm:"primaryKey"`
	UserID  string `json:"user_id"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}
