package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/akshaypatil3096/Bulk-Upload-User-Data/model"
	"go.uber.org/zap"
)

func main() {
	var userData model.User
	resp, err := http.Get("https://random-data-api.com/api/users/random_user")
	if err != nil {
		zap.S().Error(err.Error())
	}
	json.NewDecoder(resp.Body).Decode(&userData)
	csvFile, err := os.Create("user_data.csv")
	if err != nil {
		zap.S().Error(err.Error())
	}

	csvWrite := csv.NewWriter(csvFile)
	var row []string
	row = append(row, strconv.Itoa(int(userData.ID)))
	row = append(row, userData.UID)
	row = append(row, userData.Password)
	row = append(row, userData.FirstName)
	row = append(row, userData.LastName)
	row = append(row, userData.Username)
	row = append(row, userData.Email)
	row = append(row, userData.Avatar)
	row = append(row, userData.Gender)
	row = append(row, userData.PhoneNumber)
	row = append(row, userData.SocialInsuranceNumber)
	row = append(row, userData.DateOfBirth)
	csvWrite.Write(row)
	csvWrite.Flush()
}
