package service

import (
	"Market_backend/internal/common/types"
	orderRepo "Market_backend/internal/order/repository"
	paymentRepo "Market_backend/internal/payment/repository"
	"Market_backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	paymentRepo *paymentRepo.PaymentRepository
	orderRepo   *orderRepo.OrderRepository
}

func NewPaymentService(paymentRepo *paymentRepo.PaymentRepository, orderRepo *orderRepo.OrderRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo, orderRepo: orderRepo}
}

// CreatePayment создаёт Payment и возвращает confirmation_url
func (s *PaymentService) CreatePayment(order *models.Order, method types.PaymentMethod) (*models.Payment, string, error) {
	payment := &models.Payment{
		ID:        uuid.New(),
		OrderID:   order.ID,
		Method:    method,
		Status:    types.PaymentStatusPending,
		Amount:    order.Total,
		Currency:  "RUB",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Сохраняем запись Payment в БД
	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, "", err
	}

	// Создание платежа через API ЮKassa
	confirmationURL, paymentID, err := createYooKassaPayment(payment)
	if err != nil {
		return nil, "", err
	}

	// Сохраняем paymentID, полученный от ЮKassa
	payment.PaymentID = paymentID
	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, "", err
	}

	return payment, confirmationURL, nil
}

// UpdatePaymentStatus обновляет статус Payment и связанного заказа
func (s *PaymentService) UpdatePaymentStatus(paymentID string, status types.PaymentStatus) error {
	payment, err := s.paymentRepo.GetByPaymentID(paymentID)
	if err != nil {
		return err
	}

	payment.Status = status
	payment.UpdatedAt = time.Now()
	if err := s.paymentRepo.Update(payment); err != nil {
		return err
	}

	if status == types.PaymentStatusSucceeded {
		return s.orderRepo.ChangeStatus(payment.OrderID, types.Paid)
	}

	return nil
}

type YooKassaRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	PaymentMethodData struct {
		Type string `json:"type"` // "bank_card" или "sbp"
	} `json:"payment_method_data"`
	Confirmation struct {
		Type      string `json:"type"`       // "redirect"
		ReturnURL string `json:"return_url"` // URL возврата после оплаты
	} `json:"confirmation"`
	Description string            `json:"description"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Test        bool              `json:"test"` // для тестового режима
}

type YooKassaResponse struct {
	ID           string `json:"id"`
	Confirmation struct {
		ConfirmationURL string `json:"confirmation_url"`
	} `json:"confirmation"`
}

func createYooKassaPayment(payment *models.Payment) (string, string, error) {
	shopID := os.Getenv("YKASSA_SHOP_ID")
	secretKey := os.Getenv("YKASSA_SECRET_KEY")
	testMode := os.Getenv("YKASSA_TEST_MODE") == "true"

	reqBody := YooKassaRequest{}
	reqBody.Amount.Value = fmt.Sprintf("%.2f", payment.Amount)
	reqBody.Amount.Currency = payment.Currency
	reqBody.PaymentMethodData.Type = string(payment.Method) // "bank_card" или "sbp"
	reqBody.Confirmation.Type = "redirect"
	reqBody.Confirmation.ReturnURL = "http://localhost:3000/profile?tab=orders" // URL возврата
	reqBody.Description = fmt.Sprintf("Order %s", payment.OrderID)
	reqBody.Metadata = map[string]string{"order_id": payment.OrderID.String()}
	reqBody.Test = testMode

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", err
	}

	url := "https://api.yookassa.ru/v3/payments"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", "", err
	}

	req.SetBasicAuth(shopID, secretKey)
	req.Header.Set("Content-Type", "application/json")

	// Idempotence-Key — уникальный ключ для предотвращения двойной оплаты
	idempotenceKey := uuid.New().String()
	req.Header.Set("Idempotence-Key", idempotenceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return "", "", fmt.Errorf("YooKassa error: %s", string(body))
	}

	var ykResp YooKassaResponse
	if err := json.Unmarshal(body, &ykResp); err != nil {
		return "", "", err
	}

	return ykResp.Confirmation.ConfirmationURL, ykResp.ID, nil
}
