package service

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/akshaypatil3096/Bulk-Upload-User-Data/model/user"
	"go.uber.org/zap"
)

var totalInsertedRecords int = 0

func InsertData(arg string) error {
	now := time.Now()
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	insertsDoneSignal := make(chan bool, 1)
	inserts, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}

	fmt.Println("total number of inserts are: ", inserts)
	csvFile, err := os.Create("user_data.csv")
	if err != nil {
		zap.S().Error(err.Error())
	}

	csvWrite := csv.NewWriter(csvFile)
	go func() {
		writeCSVFile(inserts, csvFile, csvWrite)
		insertsDoneSignal <- true
	}()

	select {
	case <-killSignal:
		fmt.Println("inserts done: ", false, " total number of records proceed: ", totalInsertedRecords)
		csvWrite.Flush()
		csvFile.Close()
		writeStatsFile(inserts, totalInsertedRecords, now)
	case done := <-insertsDoneSignal:
		fmt.Println("inserts done: ", done)
		writeStatsFile(inserts, totalInsertedRecords, now)
	}

	return nil
}

func writeCSVFile(inserts int, csvFile *os.File, csvWrite *csv.Writer) {
	writeHeaderRow(csvFile, csvWrite)
	writeData(inserts, csvFile, csvWrite)
	defer csvWrite.Flush()
}

func writeStatsFile(inserts, totalInsertedRecords int, now time.Time) {
	statsFile, err := os.Create("stats.txt")
	if err != nil {
		zap.S().Error(err.Error())
	}

	timeOfExecution := "Time of Execution: " + time.Now().String() + "\n"
	statsFile.WriteString(timeOfExecution)
	totalRecords := "Total Number of Records Requested: " + strconv.Itoa(inserts) + "\n"
	statsFile.WriteString(totalRecords)
	insertedRecords := "Total Number of Records Inserted: " + strconv.Itoa(totalInsertedRecords) + "\n"
	statsFile.WriteString(insertedRecords)
	timeDuration := "Time Duration: " + time.Since(now).String() + "\n"
	statsFile.WriteString(timeDuration)
	statsFile.Close()
}

func writeHeaderRow(csvFile *os.File, csvWrite *csv.Writer) {
	headerLevel1 := []string{"ID", "UID", "Password", "First Name", "Last Name", "Username", "Email", "Avatar", "Gender", "Phone Number", "Social Insurance Number", "DOB", "Employment", "", "Address", "", "", "", "", "", "", "", "CreditCard", "Subscription"}
	headerLevel2 := []string{"", "", "", "", "", "", "", "", "", "", "", "", "Title", "Key Skill", "City", "Street Name", "Street Address", "Zip Code", "State", "Country", "Coordinates", "", "CC Number", "Plan", "Status", "Payment Method", "Term"}
	headerLevel3 := []string{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "Lat", "Lng"}
	csvWrite.Write(headerLevel1)
	csvWrite.Write(headerLevel2)
	csvWrite.Write(headerLevel3)
}

func writeData(inserts int, csvFile *os.File, csvWrite *csv.Writer) int {
	fmt.Println("inside writeData method")
	for i := 0; i < inserts; i++ {
		var user user.User
		resp, err := http.Get("https://random-data-api.com/api/users/random_user")
		if err != nil {
			zap.S().Error(err.Error())
		}
		json.NewDecoder(resp.Body).Decode(&user)
		/* fmt.Println(user.Gender)
		if user.Gender != "Male" && user.Gender != "Female" {
			continue
		} */
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
		row = append(row, user.CreditCard.CCNumber)
		row = append(row, user.Subscription.Plan)
		row = append(row, user.Subscription.Status)
		row = append(row, user.Subscription.PaymentMethod)
		row = append(row, user.Subscription.Term)
		csvWrite.Write(row)
		totalInsertedRecords++
	}

	fmt.Println("inside writeData method totalInsertedRecords: ", totalInsertedRecords)
	return totalInsertedRecords
}
