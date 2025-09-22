package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func sendReport(endpoint string, r Report) error {
	data, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("ошибка маршалинга отчета: %w", err)
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("ошибка при отправке: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("сервер вернул ошибку: %s", resp.Status)
	}
	return nil
}

func saveReport(path string, r Report) error {
	jsonData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling data: ", err)
		return err
	}
	err = os.WriteFile("person.json", jsonData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return err
	}
	fmt.Println("JSON data successfully written to person.json")
	return nil
}

func CreateReportURL(endpoint string, taskID string) (string, error) {
	return "", nil
}
