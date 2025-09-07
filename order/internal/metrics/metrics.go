package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "order-service"
)

// =============================================================================
// METER - ФАБРИКА ДЛЯ СОЗДАНИЯ МЕТРИК
// =============================================================================
//
// Meter в OpenTelemetry - это фабрика для создания инструментов измерения метрик.
// Каждый КОМПОНЕНТ должен иметь свой meter с уникальным именем.
//
// АРХИТЕКТУРА ВЗАИМОДЕЙСТВИЯ:
//
//  1. platform/metrics инициализирует MeterProvider:
//     platform.InitProvider() → otel.SetMeterProvider(meterProvider)
//
//  2. ufo/metrics создает свой Meter:
//     otel.Meter("ufo-service") → получает глобальный MeterProvider
//
//  3. Meter создает метрики через MeterProvider:
//     meter.Int64Counter() → meterProvider.createCounter()
//
//  4. Метрики отправляются через Reader в MeterProvider:
//     Counter.Add() → Reader.collect() → Exporter.export() → OTLP Collector
//
// СХЕМА КОМПОНЕНТОВ:
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                     GLOBAL OTEL REGISTRY                           │
// │  otel.SetMeterProvider(provider) ← platform/metrics                │
// │  otel.Meter(name) → provider     ← ufo/metrics                     │
// └─────────────────────────────────────────────────────────────────────┘
//
//	↓
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                    METER PROVIDER (один)                           │
// │  ┌─────────────────────┐  ┌─────────────────────┐                  │
// │  │   Reader            │  │   Exporter          │                  │
// │  │ - Периодически      │  │ - Отправляет в      │                  │
// │  │   читает метрики    │  │   OTLP Collector    │                  │
// │  │ - Агрегирует        │  │ - Форматирует       │                  │
// │  │   данные            │  │   протокол          │                  │
// │  └─────────────────────┘  └─────────────────────┘                  │
// └─────────────────────────────────────────────────────────────────────┘
//
//	↓
//
// ┌─────────────────────────────────────────────────────────────────────┐
// │                     METERS (много)                                 │
// │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐     │
// │  │ ufo-service     │  │ database        │  │ http-client     │     │
// │  │ - RequestsTotal │  │ - Connections   │  │ - Requests      │     │
// │  │ - SightingsTotal│  │ - QueryDuration │  │ - Errors        │     │
// │  │ - AnalysisTime  │  │ - PoolSize      │  │ - Duration      │     │
// │  └─────────────────┘  └─────────────────┘  └─────────────────┘     │
// └─────────────────────────────────────────────────────────────────────┘
//
// ВАЖНЫЕ ПРИНЦИПЫ:
//
// 1. MeterProvider ОДИН - управляет инфраструктурой отправки метрик
// 2. Meter МНОГО - один на каждый логический компонент (сервис, библиотека)
// 3. Meter получает MeterProvider из глобального registry OpenTelemetry
// 4. Все метрики из всех Meter'ов отправляются через один MeterProvider
// 5. В Prometheus метрики группируются по label'у otel_scope_name
//
// Meter предоставляет методы для создания различных типов метрик:
// - Counter - монотонно возрастающий счетчик
// - UpDownCounter - счетчик, который может увеличиваться и уменьшаться
// - Histogram - распределение значений с bucketing
// - Gauge - моментальное значение (через UpDownCounter или Callback)
//
// Важно: meter должен быть создан один раз и переиспользоваться в рамках компонента
var meter = otel.Meter(serviceName)

// =============================================================================
// ТИПЫ МЕТРИК В OPENTELEMETRY
// =============================================================================
//
// 1. COUNTER (Счетчик) - metric.Int64Counter
//    - Монотонно возрастающее значение (только увеличивается)
//    - Используется для: количество запросов, ошибок, событий
//    - Пример: общее количество HTTP запросов
//    - Методы: Add() - добавить положительное значение
//
// 2. UPDOWNCOUNTER (Двунаправленный счетчик) - metric.Int64UpDownCounter
//    - Может увеличиваться и уменьшаться
//    - Используется для: активные соединения, размер очереди, память
//    - Пример: количество активных gRPC соединений
//    - Методы: Add() - добавить (может быть отрицательным)
//
// 3. HISTOGRAM (Гистограмма) - metric.Float64Histogram
//    - Распределение наблюдений в bucket'ах
//    - Автоматически создает метрики: _count, _sum, _bucket
//    - Используется для: время ответа, размер запроса, задержки
//    - Пример: время выполнения HTTP запроса
//    - Методы: Record() - записать наблюдение
//
// 4. GAUGE (Датчик) - НЕТ отдельного типа в OpenTelemetry!
//    - В OpenTelemetry нет прямого аналога Prometheus Gauge
//    - Для gauge-подобных метрик используются:
//      а) UpDownCounter - когда значение контролируется приложением
//      б) Асинхронные Observable - когда значение нужно читать по требованию
//    - Примеры: температура CPU, использование памяти, размер кэша
//    - Для простых случаев используйте UpDownCounter как gauge

var (
	// OrdersTotal - COUNTER для подсчета созданных заказов
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: бизнес-метрика для отслеживания количество общих заказов
	// Лейблы: нет (простой счетчик без группировки)
	OrdersTotal metric.Int64Counter

	// OrdersRevenueTotal - COUNTER для подсчета суммарной выручки
	// Тип: Float64Counter (монотонно возрастающий)
	OrdersRevenueTotal metric.Float64Counter

	// RequestsTotal - COUNTER для подсчета общего количества запросов
	// Тип: Int64Counter (монотонно возрастающий)
	// Использование: подсчет всех gRPC запросов с разбивкой по методам и статусам
	// Лейблы: method (название метода), status (success/error)
	RequestsTotal metric.Int64Counter

	// RequestDuration - HISTOGRAM для измерения времени выполнения запросов
	// Тип: Float64Histogram (распределение значений)
	// Использование: SLA мониторинг - отслеживание времени ответа API
	// Позволяет строить percentile (p50, p95, p99) для анализа производительности
	RequestDuration metric.Float64Histogram
)

// InitMetrics инициализирует все метрики Order сервиса
// Должна быть вызвана один раз при старте приложения после инициализации OpenTelemetry провайдера
func InitMetrics() error {
	var err error

	// Создаем счетчик созданных заказов
	OrdersTotal, err = meter.Int64Counter(
		"orders_total",
		metric.WithDescription("Total number of created orders"),
	)
	if err != nil {
		return err
	}

	// Создаем счетчик суммарной выручки
	OrdersRevenueTotal, err = meter.Float64Counter(
		"orders_revenue_total",
		metric.WithDescription("Total number of created orders"),
	)
	if err != nil {
		return err
	}

	// Создаем счетчик запросов с описанием для документации
	RequestsTotal, err = meter.Int64Counter(
		"orders_requests_total",
		metric.WithDescription("Total number of Order service requests"),
	)
	if err != nil {
		return err
	}

	// Создаем гистограмму времени запросов с правильными bucket'ами для gRPC
	// Bucket'ы оптимизированы для времени отклика в диапазоне от микросекунд до секунд
	RequestDuration, err = meter.Float64Histogram(
		"orders_request_duration_seconds",
		metric.WithDescription("Duration of gRPC requests"),
		metric.WithUnit("s"),
		// Добавляем explicit bucket boundaries для более точного измерения gRPC запросов
		// 1ms, 2ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2s, 5s
		metric.WithExplicitBucketBoundaries(
			0.001, 0.002, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.0, 5.0,
		),
	)
	if err != nil {
		return err
	}

	return nil
}
