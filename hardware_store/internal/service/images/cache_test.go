package images

// import (
// 	"context"
// 	"hardware_store/internal/model/images"
// 	"hardware_store/internal/redis"
// 	"hardware_store/internal/storage/cache"
// 	"log/slog"
// 	"testing"
// 	"time"

// 	"github.com/alicebob/miniredis/v2"
// 	"github.com/google/uuid"
// 	goredis "github.com/redis/go-redis/v9"
// )

// func TestGetImageByProduct_CacheHit(t *testing.T) {
// 	ctx := context.Background()

// 	// Создаём минимальный Redis сервер для тестирования
// 	mr := miniredis.RunT(t)
// 	defer mr.Close()

// 	// Подключаемся к минимальному Redis
// 	rdb := goredis.NewClient(&goredis.Options{
// 		Addr: mr.Addr(),
// 	})

// 	// Создаём ImageCache
// 	logger := slog.Default()
// 	redisClient := &redis.Client{
// 		Client: rdb,
// 	}
// 	imageCache := cache.NewImageCache(
// 		redisClient,
// 		"img:",
// 		1*time.Hour,
// 	)

// 	// Создаём mock репозитория
// 	mockRepo := &mockImagesRepository{
// 		shouldFail: false,
// 	}

// 	service := NewImageService(mockRepo, imageCache, &mockTxManager{}, logger)

// 	productID := uuid.New()
// 	expectedImage := images.Images{
// 		ImageID: uuid.New(),
// 		// Image:   []byte("test image data"),
// 	}

// 	// Устанавливаем данные в кэш вручную для первого вызова
// 	mockRepo.expectedImage = expectedImage
// 	img, err := service.GetImageByProduct(ctx, productID)
// 	if err != nil {
// 		t.Fatalf("first call failed: %v", err)
// 	}

// 	// Второй вызов должен вернуться из кэша
// 	mockRepo.callCount = 0 // обнуляем счётчик
// 	img2, err := service.GetImageByProduct(ctx, productID)
// 	if err != nil {
// 		t.Fatalf("second call failed: %v", err)
// 	}

// 	if mockRepo.callCount != 0 {
// 		t.Errorf("expected cache hit, but repo was called %d times", mockRepo.callCount)
// 	}

// 	if img.ImageID != img2.ImageID {
// 		t.Errorf("expected same image from cache, got different")
// 	}
// }

// func TestGetImageByProduct_CacheMiss(t *testing.T) {
// 	ctx := context.Background()

// 	mr := miniredis.RunT(t)
// 	defer mr.Close()

// 	rdb := goredis.NewClient(&goredis.Options{
// 		Addr: mr.Addr(),
// 	})

// 	logger := slog.Default()
// 	redisClient := &redis.Client{
// 		Client: rdb,
// 	}
// 	imageCache := cache.NewImageCache(
// 		redisClient,
// 		"img:",
// 		1*time.Hour,
// 	)

// 	mockRepo := &mockImagesRepository{
// 		shouldFail: false,
// 	}

// 	service := NewImageService(mockRepo, imageCache, &mockTxManager{}, logger)

// 	productID := uuid.New()
// 	expectedImage := images.Images{
// 		ImageID: uuid.New(),
// 		Image:   []byte("test image"),
// 	}
// 	mockRepo.expectedImage = expectedImage

// 	img, err := service.GetImageByProduct(ctx, productID)
// 	if err != nil {
// 		t.Fatalf("GetImageByProduct failed: %v", err)
// 	}

// 	if mockRepo.callCount != 1 {
// 		t.Errorf("expected repo to be called once for cache miss, got %d calls", mockRepo.callCount)
// 	}

// 	if img.ImageID != expectedImage.ImageID {
// 		t.Errorf("expected image ID %s, got %s", expectedImage.ImageID, img.ImageID)
// 	}
// }

// func TestDeleteImage_InvalidatesCache(t *testing.T) {
// 	ctx := context.Background()

// 	mr := miniredis.RunT(t)
// 	defer mr.Close()

// 	rdb := goredis.NewClient(&goredis.Options{
// 		Addr: mr.Addr(),
// 	})

// 	logger := slog.Default()
// 	redisClient := &redis.Client{
// 		Client: rdb,
// 	}
// 	imageCache := cache.NewImageCache(
// 		redisClient,
// 		"img:",
// 		1*time.Hour,
// 	)

// 	mockRepo := &mockImagesRepository{
// 		shouldFail: false,
// 	}

// 	service := NewImageService(mockRepo, imageCache, &mockTxManager{}, logger)

// 	imageID := uuid.New()

// 	// Удаляем изображение
// 	err := service.DeleteImage(ctx, imageID)
// 	if err != nil {
// 		t.Fatalf("DeleteImage failed: %v", err)
// 	}

// 	if mockRepo.callCount != 1 {
// 		t.Errorf("expected repo.Delete to be called once, got %d calls", mockRepo.callCount)
// 	}
// }

// func TestUpdateImage_InvalidatesCache(t *testing.T) {
// 	ctx := context.Background()

// 	mr := miniredis.RunT(t)
// 	defer mr.Close()

// 	rdb := goredis.NewClient(&goredis.Options{
// 		Addr: mr.Addr(),
// 	})

// 	logger := slog.Default()
// 	redisClient := &redis.Client{
// 		Client: rdb,
// 	}
// 	imageCache := cache.NewImageCache(
// 		redisClient,
// 		"img:",
// 		1*time.Hour,
// 	)

// 	mockRepo := &mockImagesRepository{
// 		shouldFail: false,
// 	}

// 	service := NewImageService(mockRepo, imageCache, &mockTxManager{}, logger)

// 	imageID := uuid.New()
// 	newImageData := []byte("updated image data")

// 	// Обновляем изображение
// 	err := service.UpdateImage(ctx, imageID, newImageData)
// 	if err != nil {
// 		t.Fatalf("UpdateImage failed: %v", err)
// 	}

// 	if mockRepo.callCount != 1 {
// 		t.Errorf("expected repo.Update to be called once, got %d calls", mockRepo.callCount)
// 	}
// }

// // Mock структуры

// type mockImagesRepository struct {
// 	shouldFail    bool
// 	callCount     int
// 	expectedImage images.Images
// }

// func (m *mockImagesRepository) Insert(ctx context.Context, image images.Images) error {
// 	m.callCount++
// 	if m.shouldFail {
// 		return ErrMockFailure
// 	}
// 	return nil
// }

// func (m *mockImagesRepository) Update(ctx context.Context, image images.Images) error {
// 	m.callCount++
// 	if m.shouldFail {
// 		return ErrMockFailure
// 	}
// 	return nil
// }

// func (m *mockImagesRepository) Delete(ctx context.Context, imagesID uuid.UUID) error {
// 	m.callCount++
// 	if m.shouldFail {
// 		return ErrMockFailure
// 	}
// 	return nil
// }

// func (m *mockImagesRepository) GetByProduct(ctx context.Context, productID uuid.UUID) (images.Images, error) {
// 	m.callCount++
// 	if m.shouldFail {
// 		return images.Images{}, ErrMockFailure
// 	}
// 	return m.expectedImage, nil
// }

// func (m *mockImagesRepository) GetById(ctx context.Context, imagesID uuid.UUID) (images.Images, error) {
// 	m.callCount++
// 	if m.shouldFail {
// 		return images.Images{}, ErrMockFailure
// 	}
// 	return m.expectedImage, nil
// }

// func (m *mockImagesRepository) Attach(ctx context.Context, imagesID, productID uuid.UUID) error {
// 	m.callCount++
// 	if m.shouldFail {
// 		return ErrMockFailure
// 	}
// 	return nil
// }

// type mockTxManager struct{}

// func (m *mockTxManager) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
// 	return fn(ctx)
// }