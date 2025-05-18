package load

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type User struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	Addresses   []Address `json:"addresses"`
}

func LoadData(filePath string) ([]User, error) {
	if filePath == "" {
		filePath = "data/users_data.json"
	}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	byteDate, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var users []User
	err = json.Unmarshal(byteDate, &users)
	if err != nil {
		return nil, fmt.Errorf("Error un-marshalling JSON: %w", err)
	}
	return users, nil
}
