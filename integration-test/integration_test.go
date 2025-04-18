package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type PVZRequest struct {
	City string `json:"city"`
}

type PVZResponse struct {
	ID               string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}

type ReceptionRequest struct {
	PVZID string `json:"pvzId"`
}

type ReceptionResponse struct {
	ID       string `json:"id"`
	DateTime string `json:"dateTime"`
	PVZID    string `json:"pvzId"`
	Status   string `json:"status"`
}

type ProductRequest struct {
	Type  string `json:"type"`
	PVZID string `json:"pvzId"`
}

type ProductResponse struct {
	ID          string `json:"id"`
	DateTime    string `json:"dateTime"`
	Type        string `json:"type"`
	ReceptionID string `json:"receptionId"`
}

type PVZInfo struct {
	PVZ        PVZResponse             `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions"`
}

type ReceptionWithProducts struct {
	Reception ReceptionResponse `json:"reception"`
	Products  []ProductResponse `json:"products"`
}

func getDummyToken(t *testing.T, role string) string {
	url := "http://localhost:8080/dummyLogin"
	body := []byte(fmt.Sprintf(`{"role": "%s"}`, role))
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Skipf("Skipping test: Failed to get dummy token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		t.Fatalf("Failed to decode token response: %v", err)
	}

	return tokenResp.Token
}

func createPVZ(t *testing.T, token, city string) string {
	url := "http://localhost:8080/pvz"
	body := []byte(fmt.Sprintf(`{"city": "%s"}`, city))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create PVZ: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var pvz PVZResponse
	if err := json.NewDecoder(resp.Body).Decode(&pvz); err != nil {
		t.Fatalf("Failed to decode PVZ response: %v", err)
	}

	return pvz.ID
}

func createReception(t *testing.T, token, pvzID string) string {
	url := "http://localhost:8080/receptions"
	body := []byte(fmt.Sprintf(`{"pvzId": "%s"}`, pvzID))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to create reception: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var reception ReceptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&reception); err != nil {
		t.Fatalf("Failed to decode reception response: %v", err)
	}

	return reception.ID
}

func addProduct(t *testing.T, token, pvzID, productType string) {
	url := "http://localhost:8080/products"
	body := []byte(fmt.Sprintf(`{"type": "%s", "pvzId": "%s"}`, productType, pvzID))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to add product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}
}

func closeReception(t *testing.T, token, pvzID string) {
	url := fmt.Sprintf("http://localhost:8080/pvz/%s/close_last_reception", pvzID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to close reception: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}

func getPVZInfo(t *testing.T, token, pvzID string) PVZInfo {
	url := "http://localhost:8080/pvz?page=1&limit=30"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to get PVZ info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var pvzList []PVZInfo
	if err := json.NewDecoder(resp.Body).Decode(&pvzList); err != nil {
		t.Fatalf("Failed to decode PVZ list: %v", err)
	}

	for _, pvzInfo := range pvzList {
		if pvzInfo.PVZ.ID == pvzID {
			return pvzInfo
		}
	}

	t.Fatalf("PVZ with ID %s not found", pvzID)
	return PVZInfo{}
}

func TestIntegrationFlow(t *testing.T) {
	// Шаг 1: Получаем токены
	moderatorToken := getDummyToken(t, "moderator")
	employeeToken := getDummyToken(t, "employee")

	// Шаг 2: Создаем ПВЗ
	pvzID := createPVZ(t, moderatorToken, "Москва")

	// Шаг 3: Создаем приёмку
	_ = createReception(t, employeeToken, pvzID)

	// Шаг 4: Добавляем 50 товаров
	for i := 0; i < 50; i++ {
		addProduct(t, employeeToken, pvzID, "электроника")
	}

	// Шаг 5: Проверяем результаты
	pvzInfo := getPVZInfo(t, moderatorToken, pvzID)

	if len(pvzInfo.Receptions) != 1 {
		t.Fatalf("Expected 1 reception, got %d", len(pvzInfo.Receptions))
	}

	productsCount := len(pvzInfo.Receptions[0].Products)
	if productsCount != 50 {
		t.Fatalf("Expected 50 products, got %d", productsCount)
	}
	// Шаг 6: Закрываем приёмку
	closeReception(t, employeeToken, pvzID)

}
