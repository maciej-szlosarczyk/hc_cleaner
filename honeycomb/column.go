package honeycomb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Column struct {
	ID          string    `json:"id"`
	KeyName     string    `json:"key_name"`
	Type        string    `json:"type"`
	Hidden      bool      `json:"hidden"`
	Description string    `json:"description"`
	LastWritten time.Time `json:"last_written"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type HoneycombDataset struct {
	Name   string
	ApiKey string
}

func DeleteColumnsWithPrefix(datasetName string, apiKey string, prefix string) error {
	dataset := HoneycombDataset{Name: datasetName, ApiKey: apiKey}
	allColumns, err := GetColumns(dataset)

	if err != nil {
		fmt.Println(err)
		return err
	}

	filteredColumns := FilterPrefixColumns(allColumns, prefix)
	fmt.Printf("Total columns to delete: %d.\n", len(filteredColumns))

	err = DeleteColumns(filteredColumns, dataset)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func DeleteInactiveColumns(datasetName string, apiKey string, since int64) error {
	dataset := HoneycombDataset{Name: datasetName, ApiKey: apiKey}
	allColumns, err := GetColumns(dataset)

	if err != nil {
		fmt.Println(err)
		return err
	}

	filteredColumns := FilterInactiveColumns(allColumns, since)
	fmt.Printf("Total columns to delete: %d.\n", len(filteredColumns))

	err = DeleteColumns(filteredColumns, dataset)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func GetColumns(dataset HoneycombDataset) ([]Column, error) {
	var allColumns []Column

	client := http.Client{}
	url := fmt.Sprintf("https://api.honeycomb.io/1/columns/%s", dataset.Name)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Honeycomb-Team", dataset.ApiKey)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return allColumns, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &allColumns)

	if err != nil {
		fmt.Println(err)
		return allColumns, err
	}

	return allColumns, err
}

func FilterPrefixColumns(columns []Column, prefix string) []Column {
	var inactiveColumns []Column

	for _, element := range columns {
		if strings.HasPrefix(element.KeyName, prefix) {
			inactiveColumns = append(inactiveColumns, element)
		}
	}

	return inactiveColumns
}

func FilterInactiveColumns(columns []Column, daysSince int64) []Column {
	cutoffDate := time.Now()
	since := time.Hour * 24 * time.Duration(-daysSince)
	cutoffDate = cutoffDate.Add(since)
	var inactiveColumns []Column

	for _, element := range columns {
		if element.LastWritten.Before(cutoffDate) {
			inactiveColumns = append(inactiveColumns, element)
		}
	}

	return inactiveColumns
}

func DeleteColumn(column Column, dataset HoneycombDataset) error {
	fmt.Printf("Deleting column: %s.\n", column.KeyName)
	client := http.Client{}
	url := fmt.Sprintf("https://api.honeycomb.io/1/columns/%s/%s", dataset.Name, column.ID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Honeycomb-Team", dataset.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	return nil
}

func DeleteColumns(columns []Column, dataset HoneycombDataset) error {
	for _, element := range columns {
		err := DeleteColumn(element, dataset)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
