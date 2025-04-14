package t212

import (
	"fmt"
	"time"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

type APIClient struct {
	Key string
}

type CreateReportPayload struct {
	DataIncluded map[string]bool `json:"dataIncluded"`
	TimeFrom     string          `json:"timeFrom"` // Use string for RFC3339 formatted time
	TimeTo       string          `json:"timeTo"`
}

type Export struct {
	ReportId uint `json:"reportId"`
}

type Report struct {
	Id           uint            `json:"reportId"`
	TimeFrom     string          `json:"timeFrom"`
	TimeTo       string          `json:"timeTo"`
	DataIncluded map[string]bool `json:"dataIncluded"`
	Status       string          `json:"status"`
	DownloadLink string          `json:"downloadLink"`
}

func (apiClient APIClient) CreateReport(fromDt *time.Time, toDt *time.Time) *uint {

	const url string = "https://live.trading212.com/api/v0/history/exports"

	dataIncluded := map[string]bool{
		"includeDividends":    true,
		"includeInterest":     true,
		"includeOrders":       true,
		"includeTransactions": true,
	}
	
	payload := CreateReportPayload{
		DataIncluded: dataIncluded,
		TimeFrom: fromDt.Format(time.RFC3339),
		TimeTo:   toDt.Format(time.RFC3339),
	}

	payloadEncoded, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadEncoded))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", apiClient.Key)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("  Response Status: %v\n", response.Status)

	if response.Status != "200 OK" {
		return nil
	}
	
	defer response.Body.Close()
	reponseEncoded, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var reponseBody Export
	err = json.Unmarshal(reponseEncoded, &reponseBody)
	if err != nil {
		panic(err)
	}

	return &reponseBody.ReportId
}

func (apiClient APIClient) ListReports() []Report {

	const url string = "https://live.trading212.com/api/v0/history/exports"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", apiClient.Key)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("  Response Status: %v\n", response.Status)

	if response.Status != "200 OK" {
		return nil
	}

	defer response.Body.Close()
	responseEncoded, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var reponseBody []Report
	err = json.Unmarshal(responseEncoded, &reponseBody)
	if err != nil {
		panic(err)
	}

	return reponseBody
}


func (report Report) Download() []byte {

	req, err := http.NewRequest("GET", report.DownloadLink, nil)
	if err != nil {
		panic(err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Printf("  Response Status: %v\n", response.Status)

	if response.Status != "200 OK" {
		return nil
	}

	defer response.Body.Close()
	responseEncoded, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	return responseEncoded
}
