package telegram

import (
	"bytes"
	"context"
	"embed"
	"text/template"

	"github.com/dfg007star/go_rocket/notification/internal/client/http"
	"github.com/dfg007star/go_rocket/notification/internal/model"
	"github.com/dfg007star/go_rocket/platform/pkg/logger"
	"go.uber.org/zap"
)

const chatID = 234586218

//go:embed templates/assembled_notification.tmpl
//go:embed templates/paid_notification.tmpl
var templateFS embed.FS

type OrderAssembledData struct {
	EventUuid    string
	OrderUuid    string
	UserUuid     string
	BuildTimeSec int64
}

type OrderPaidData struct {
	OrderUuid       string
	UserUuid        string
	PaymentMethod   string
	TransactionUuid string
}

var paidTemplate = template.Must(template.ParseFS(templateFS, "templates/paid_notification.tmpl"))
var assembledTemplate = template.Must(template.ParseFS(templateFS, "templates/assembled_notification.tmpl"))

type service struct {
	telegramClient http.TelegramClient
}

// NewService создает новый Telegram сервис
func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

// SendPaidNotification отправляет уведомление об оплате заказа
func (s *service) SendPaidNotification(ctx context.Context, uuid string, order model.OrderPaidEvent) error {
	message, err := s.buildPaidMessage(order)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

// buildPaidMessage создает сообщение об оплате заказа
func (s *service) buildPaidMessage(order model.OrderPaidEvent) (string, error) {
	data := OrderPaidData{
		OrderUuid:       order.OrderUuid,
		UserUuid:        order.UserUuid,
		PaymentMethod:   order.PaymentMethod,
		TransactionUuid: order.TransactionUuid,
	}

	var buf bytes.Buffer
	err := paidTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SendAssembledNotification отправляет уведомление о сборке заказа
func (s *service) SendAssembledNotification(ctx context.Context, uuid string, order model.ShipAssembledEvent) error {
	message, err := s.buildAssembledMessage(order)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int("chat_id", chatID), zap.String("message", message))
	return nil
}

// buildPaidMessage создает сообщение о сборке заказа
func (s *service) buildAssembledMessage(order model.ShipAssembledEvent) (string, error) {
	data := OrderAssembledData{
		EventUuid:    order.EventUuid,
		OrderUuid:    order.OrderUuid,
		UserUuid:     order.UserUuid,
		BuildTimeSec: order.BuildTimeSec,
	}

	var buf bytes.Buffer
	err := assembledTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
