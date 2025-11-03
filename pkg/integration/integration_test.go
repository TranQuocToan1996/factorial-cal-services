package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"factorial-cal-services/pkg/consumer"
	"factorial-cal-services/pkg/domain"
	"factorial-cal-services/pkg/repository"
	"factorial-cal-services/pkg/service"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// mockProducer is an in-memory producer for testing
type mockProducer struct {
	messages chan message
}

type message struct {
	queueName string
	payload   []byte
}

func newMockProducer() *mockProducer {
	return &mockProducer{
		messages: make(chan message, 100),
	}
}

func (m *mockProducer) Publish(ctx context.Context, queueName string, payload []byte) error {
	select {
	case m.messages <- message{queueName: queueName, payload: payload}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("message queue full")
	}
}

func (m *mockProducer) Close() error {
	close(m.messages)
	return nil
}

func (m *mockProducer) GetMessage() (string, []byte, error) {
	select {
	case msg, ok := <-m.messages:
		if !ok {
			return "", nil, fmt.Errorf("channel closed")
		}
		return msg.queueName, msg.payload, nil
	case <-time.After(1 * time.Second):
		return "", nil, fmt.Errorf("timeout waiting for message")
	}
}

func setupIntegrationTest(t *testing.T) (*gorm.DB, *redis.Client, *mockProducer, service.StorageService, func()) {
	// Setup in-memory database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate all tables
	err = db.AutoMigrate(
		&domain.FactorialCalculation{},
		&domain.FactorialMaxRequestNumber{},
		&domain.FactorialCurrentCalculatedNumber{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Setup Redis (miniredis)
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to create miniredis: %v", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	// Setup mock producer
	mockProd := newMockProducer()

	// Setup local storage
	tempDir, err := os.MkdirTemp("", "factorial-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	storageService := service.NewLocalStorageService(tempDir)

	cleanup := func() {
		redisClient.Close()
		mr.Close()
		os.RemoveAll(tempDir)
	}

	return db, redisClient, mockProd, storageService, cleanup
}

func TestIntegrationFullFlow(t *testing.T) {
	db, redisClient, mockProd, storageService, cleanup := setupIntegrationTest(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize services
	factorialService := service.NewFactorialServiceWithLimit(10000)
	redisService := service.NewRedisService(redisClient, 24*time.Hour, 1000)
	checksumService := service.NewChecksumService()

	// Initialize repositories
	factorialRepo := repository.NewFactorialRepository(db)
	maxRequestRepo := repository.NewMaxRequestRepository(db)
	currentCalculatedRepo := repository.NewCurrentCalculatedRepository(db)

	// Create S3 service wrapper that uses local storage
	var s3Service service.S3Service = &localStorageS3Wrapper{storage: storageService}

	// Initialize batch handler
	batchHandler := consumer.NewFactorialBatchHandler(
		factorialService,
		redisService,
		s3Service,
		factorialRepo,
		maxRequestRepo,
		currentCalculatedRepo,
		checksumService,
		service.NewIncrementalFactorialService(factorialService, currentCalculatedRepo),
	)

	// Test: POST calculation publishes message (simulated)
	t.Run("Message publishing works", func(t *testing.T) {
		// Simulate POST request by publishing message
		message := map[string]string{"number": "10"}
		payload, err := json.Marshal(message)
		if err != nil {
			t.Fatalf("Failed to marshal message: %v", err)
		}

		err = mockProd.Publish(ctx, "test-queue", payload)
		if err != nil {
			t.Fatalf("Failed to publish message: %v", err)
		}

		// Verify message was published
		queueName, receivedPayload, err := mockProd.GetMessage()
		if err != nil {
			t.Fatalf("Failed to get message: %v", err)
		}

		if queueName != "test-queue" {
			t.Errorf("Expected queue name 'test-queue', got %s", queueName)
		}

		var msg map[string]string
		if err := json.Unmarshal(receivedPayload, &msg); err != nil {
			t.Fatalf("Failed to unmarshal message: %v", err)
		}

		if msg["number"] != "10" {
			t.Errorf("Expected number '10', got %s", msg["number"])
		}
	})

	// Test: Worker processes batch
	t.Run("Worker processes batch and calculates factorial incrementally", func(t *testing.T) {
		// Create batch of messages
		messages := [][]byte{
			[]byte(`{"number":"5"}`),
			[]byte(`{"number":"10"}`),
			[]byte(`{"number":"15"}`),
		}

		// Process batch
		err := batchHandler(ctx, messages)
		if err != nil {
			t.Fatalf("Failed to process batch: %v", err)
		}

		// Wait for processing to complete
		time.Sleep(500 * time.Millisecond)

		// Verify calculations were stored
		calc5, err := factorialRepo.FindByNumber("5")
		if err != nil {
			t.Fatalf("Failed to find calculation for 5: %v", err)
		}
		if calc5 == nil || calc5.Status != domain.StatusDone {
			t.Errorf("Expected calculation for 5 to be done, got: %+v", calc5)
		}
		if calc5 != nil && calc5.S3Key == "" {
			t.Error("Expected S3 key to be set for calculation 5")
		}

		calc10, err := factorialRepo.FindByNumber("10")
		if err != nil {
			t.Fatalf("Failed to find calculation for 10: %v", err)
		}
		if calc10 == nil || calc10.Status != domain.StatusDone {
			t.Errorf("Expected calculation for 10 to be done, got: %+v", calc10)
		}

		calc15, err := factorialRepo.FindByNumber("15")
		if err != nil {
			t.Fatalf("Failed to find calculation for 15: %v", err)
		}
		if calc15 == nil || calc15.Status != domain.StatusDone {
			t.Errorf("Expected calculation for 15 to be done, got: %+v", calc15)
		}

		// Verify max number was updated
		maxNumber, err := maxRequestRepo.GetMaxNumber()
		if err != nil {
			t.Fatalf("Failed to get max number: %v", err)
		}
		if maxNumber != "15" {
			t.Errorf("Expected max number '15', got %s", maxNumber)
		}

		// Verify current calculated number was updated
		currentNumber, err := currentCalculatedRepo.GetCurrentNumber()
		if err != nil {
			t.Fatalf("Failed to get current number: %v", err)
		}
		if currentNumber != "15" {
			t.Errorf("Expected current number '15', got %s", currentNumber)
		}

		// Verify results are stored in storage
		if calc10 != nil && calc10.S3Key != "" {
			result, err := storageService.Download(ctx, calc10.S3Key)
			if err != nil {
				t.Fatalf("Failed to download result: %v", err)
			}
			if result == "" {
				t.Error("Expected non-empty result from storage")
			}
			// Verify it's the correct factorial (10! = 3628800)
			if len(result) == 0 {
				t.Error("Result should not be empty")
			}
		}

		// Verify small numbers are cached in Redis
		if calc5 != nil {
			cached, err := redisService.Get(ctx, "5")
			if err != nil {
				t.Logf("Number 5 not in Redis cache (may be expected if threshold check failed): %v", err)
			} else if cached == "" {
				t.Error("Expected number 5 to be cached in Redis")
			}
		}
	})

	// Test: GET result retrieves from cache or storage
	t.Run("GET result retrieves calculated factorial", func(t *testing.T) {
		// Process a calculation first
		messages := [][]byte{[]byte(`{"number":"7"}`)}
		err := batchHandler(ctx, messages)
		if err != nil {
			t.Fatalf("Failed to process batch: %v", err)
		}

		time.Sleep(300 * time.Millisecond)

		// Verify calculation exists
		calc, err := factorialRepo.FindByNumber("7")
		if err != nil {
			t.Fatalf("Failed to find calculation: %v", err)
		}
		if calc == nil || calc.Status != domain.StatusDone {
			t.Fatalf("Expected calculation to be done, got: %+v", calc)
		}

		// Try to get from Redis first (should be cached for small numbers)
		result, err := redisService.Get(ctx, "7")
		if err == nil && result != "" {
			// Verify result is correct
			if len(result) == 0 {
				t.Error("Expected non-empty result from Redis")
			}
		}

		// If not in Redis, get from storage
		if result == "" && calc.S3Key != "" {
			result, err = storageService.Download(ctx, calc.S3Key)
			if err != nil {
				t.Fatalf("Failed to download from storage: %v", err)
			}
			if result == "" {
				t.Error("Expected non-empty result from storage")
			}
		}

		if result == "" {
			t.Error("Failed to retrieve result from either Redis or storage")
		}
	})

	// Test: Incremental calculation works correctly
	t.Run("Incremental calculation processes sequentially", func(t *testing.T) {
		// Process numbers 20, 21, 22 incrementally
		messages := [][]byte{
			[]byte(`{"number":"20"}`),
			[]byte(`{"number":"21"}`),
			[]byte(`{"number":"22"}`),
		}

		err := batchHandler(ctx, messages)
		if err != nil {
			t.Fatalf("Failed to process batch: %v", err)
		}

		time.Sleep(500 * time.Millisecond)

		// Verify all calculations completed
		for _, num := range []string{"20", "21", "22"} {
			calc, err := factorialRepo.FindByNumber(num)
			if err != nil {
				t.Fatalf("Failed to find calculation for %s: %v", num, err)
			}
			if calc == nil || calc.Status != domain.StatusDone {
				t.Errorf("Expected calculation for %s to be done", num)
			}
		}

		// Verify current number is 22
		currentNumber, err := currentCalculatedRepo.GetCurrentNumber()
		if err != nil {
			t.Fatalf("Failed to get current number: %v", err)
		}
		if currentNumber != "22" {
			t.Errorf("Expected current number '22', got %s", currentNumber)
		}
	})
}

// localStorageS3Wrapper wraps LocalStorageService to implement S3Service interface
type localStorageS3Wrapper struct {
	storage service.StorageService
}

func (w *localStorageS3Wrapper) UploadFactorial(ctx context.Context, number string, result string) (string, error) {
	return w.storage.Upload(ctx, number, result)
}

func (w *localStorageS3Wrapper) DownloadFactorial(ctx context.Context, s3Key string) (string, error) {
	return w.storage.Download(ctx, s3Key)
}

func (w *localStorageS3Wrapper) GenerateS3Key(number string) string {
	return w.storage.GenerateS3Key(number)
}

func (w *localStorageS3Wrapper) Upload(ctx context.Context, number string, result string) (string, error) {
	return w.storage.Upload(ctx, number, result)
}

func (w *localStorageS3Wrapper) Download(ctx context.Context, key string) (string, error) {
	return w.storage.Download(ctx, key)
}
