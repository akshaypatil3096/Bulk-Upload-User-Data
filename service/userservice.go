package service

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/akshaypatil3096/Bulk-Upload-User-Data/model/user"
	"go.uber.org/zap"
)

var totalInsertedRecords int = 0

func Start(arg string) error {
	inserts, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}

	InsertData(inserts, inserts, 0)
	return err
}

func InsertData(totalInserts, remainingInserts, insertionIndex int) (err error) {
	now := time.Now()
	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	insertsDoneSignal := make(chan bool, 1)
	var csvFile *os.File
	if totalInsertedRecords > 0 {
		csvFile, err = os.OpenFile("user_data.csv", os.O_APPEND|os.O_RDWR, os.ModeAppend)
	} else {
		csvFile, err = os.Create("user_data.csv")
	}

	if err != nil {
		zap.S().Error(err.Error())
	}

	csvWrite := csv.NewWriter(csvFile)
	go func() {
		writeCSVFile(totalInserts, remainingInserts, insertionIndex, csvFile, csvWrite)
		insertsDoneSignal <- true
	}()

	select {
	case <-killSignal:
		fmt.Println("inserts done: ", false, " total number of records proceed: ", totalInsertedRecords)
		csvWrite.Flush()
		csvFile.Close()
		writeStatsFile(totalInserts, totalInsertedRecords, now)
	case done := <-insertsDoneSignal:
		fmt.Println("inserts done: ", done)
		writeStatsFile(totalInserts, totalInsertedRecords, now)
	}

	return nil
}

func ResumeInsertData() error {
	totalRequestedRecords, totalInsertsCompleted, err := readStatsFile()
	totalInsertedRecords = totalInsertsCompleted
	inserts := totalRequestedRecords - totalInsertsCompleted
	if err != nil {
		return err
	}

	insertionIndex := totalInsertsCompleted + 3
	err = InsertData(totalRequestedRecords, inserts, insertionIndex)
	return err
}

func readStatsFile() (int, int, error) {
	statsFile, err := os.Open("stats.txt")
	if err != nil {
		zap.S().Error(err.Error())
	}

	scanner := bufio.NewScanner(statsFile)
	scanner.Split(bufio.ScanLines)
	var textlines []string

	for scanner.Scan() {
		textlines = append(textlines, scanner.Text())
	}

	statsFile.Close()
	statsFileMap := make(map[string]string)
	for _, eachline := range textlines {
		splits := strings.Split(eachline, ":")
		statsFileMap[splits[0]] = splits[1]
	}

	totalRequestedRecordsStr := strings.ReplaceAll(statsFileMap["Total Number of Records Requested"], " ", "")
	totalInsertsCompletedStr := strings.ReplaceAll(statsFileMap["Total Number of Records Inserted"], " ", "")
	totalRequestedRecords, err := strconv.Atoi(totalRequestedRecordsStr)
	if err != nil {
		zap.S().Error(err.Error())
	}

	totalInsertsCompleted, err := strconv.Atoi(totalInsertsCompletedStr)
	if err != nil {
		zap.S().Error(err.Error())
	}

	return totalRequestedRecords, totalInsertsCompleted, err
}

func writeCSVFile(totalInserts, remainingInserts, insertionIndex int, csvFile *os.File, csvWrite *csv.Writer) {
	if totalInsertedRecords == 0 {
		writeHeaderRow(csvFile, csvWrite)
	}

	writeData(totalInserts, remainingInserts, insertionIndex, csvFile, csvWrite)
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

func writeData(totalInserts, remainingInserts, insertionIndex int, csvFile *os.File, csvWrite *csv.Writer) int {
	fmt.Println("inside writeData method")
	fmt.Println("total number of inserts are: ", totalInserts)
	fmt.Println("remaining inserts are: ", remainingInserts)
	fmt.Println("insertion index: ", insertionIndex)
	fmt.Println("totalInsertedRecords: ", totalInsertedRecords)

	for i := 0; totalInsertedRecords < totalInserts; i++ {
		var user user.User
		resp, err := http.Get("https://random-data-api.com/api/users/random_user")
		if err != nil {
			zap.S().Error(err.Error())
		}

		json.NewDecoder(resp.Body).Decode(&user)
		if user.Gender != "Male" && user.Gender != "Female" {
			continue
		}

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

	fmt.Println("totalInsertedRecords: ", totalInsertedRecords)
	return totalInsertedRecords
}
