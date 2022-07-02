package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/akshaypatil3096/Bulk-Upload-User-Data/model"
	"go.uber.org/zap"
)

func main() {
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	now := time.Now()
	var userData []model.User
	csvFile, err := os.Create("user_data.csv")
	if err != nil {
		zap.S().Error(err.Error())
	}
	csvWrite := csv.NewWriter(csvFile)

	headerLevel1 := []string{"ID", "UID", "Password", "First Name", "Last Name", "Username", "Email", "Avatar", "Gender", "Phone Number", "Social Insurance Number", "DOB", "Employment", "", "Address", "", "", "", "", "", "", "", "CreditCard", "Subscription"}
	headerLevel2 := []string{"", "", "", "", "", "", "", "", "", "", "", "", "Title", "Key Skill", "City", "Street Name", "Street Address", "Zip Code", "State", "Country", "Coordinates", "", "CC Number", "Plan", "Status", "Payment Method", "Term"}
	headerLevel3 := []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "Lat", "Lng"}

	csvWrite.Write(headerLevel1)
	csvWrite.Write(headerLevel2)
	csvWrite.Write(headerLevel3)
	record := 0
	for i := 0; i < 25; i++ {
		var user model.User
		resp, err := http.Get("https://random-data-api.com/api/users/random_user")
		if err != nil {
			zap.S().Error(err.Error())
		}
		json.NewDecoder(resp.Body).Decode(&user)
		var row []string
		row = append(row, strconv.Itoa(int(user.ID)))
		row = append(row, user.UID)
		row = append(row, user.Password)
		row = append(row, user.FirstName)
		row = append(row, user.LastName)
		row = append(row, user.Username)
		row = append(row, user.Email)
		row = append(row, user.Avatar)
		row = append(row, user.Gender)
		row = append(row, user.PhoneNumber)
		row = append(row, user.SocialInsuranceNumber)
		row = append(row, user.DateOfBirth)
		row = append(row, user.Employment.Title)
		row = append(row, user.Employment.KeySkill)
		row = append(row, user.Address.City)
		row = append(row, user.Address.StreetName)
		row = append(row, user.Address.StreetAddress)
		row = append(row, user.Address.ZipCode)
		row = append(row, user.Address.State)
		row = append(row, user.Address.Country)
		row = append(row, strconv.FormatFloat(user.Address.Coordinates.Lat, 'E', -1, 64))
		row = append(row, strconv.FormatFloat(user.Address.Coordinates.Lng, 'E', -1, 64))
		row = append(row, user.CreditCard.CcNumber)
		row = append(row, user.Subscription.Plan)
		row = append(row, user.Subscription.Status)
		row = append(row, user.Subscription.PaymentMethod)
		row = append(row, user.Subscription.Term)

		csvWrite.Write(row)
		record++
		userData = append(userData, user)
	}

	defer csvWrite.Flush()

	statsFile, err := os.Create("stats.txt")
	if err != nil {
		zap.S().Error(err.Error())
	}
	timeOfExecution := "Time of Execution: " + time.Now().String() + "\n"
	statsFile.WriteString(timeOfExecution)
	records := "Total Number of Records Proceed: " + strconv.Itoa(record) + "\n"
	statsFile.WriteString(records)
	timeDuration := "Time Duration: " + time.Since(now).String() + "\n"
	statsFile.WriteString(timeDuration)
	statsFile.Close()

}
