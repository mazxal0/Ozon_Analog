package service

import (
	"Market_backend/internal/common/types"
	"Market_backend/internal/common/utils"
	"Market_backend/internal/product/dto"
	"Market_backend/internal/product/repository"
	"Market_backend/internal/storage"
	"Market_backend/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"os"
	"strings"
)

type ProcessorService struct {
	procRepo *repository.ProcessorRepository
	storage  *storage.MinioStorage
}

func NewProcessorService(repo *repository.ProcessorRepository, storage *storage.MinioStorage) *ProcessorService {
	return &ProcessorService{
		procRepo: repo,
		storage:  storage,
	}
}

func (s *ProcessorService) DB() *repository.ProcessorRepository {
	return s.procRepo
}

func (s *ProcessorService) CreateProcessor(dto dto.ProcessorCreateDTO) (*models.Processor, error) {
	// создаем processor с генерацией SKU
	processor := &models.Processor{
		ID:                 uuid.New(),
		SKU:                utils.GenerateSKU(),
		Name:               dto.Name,
		Brand:              dto.Brand,
		RetailPrice:        dto.RetailPrice,
		WholesalePrice:     dto.WholesalePrice,
		WholesaleMinQty:    dto.WholesaleMinQty,
		Stock:              dto.Stock,
		Line:               dto.Line,
		Architecture:       dto.Architecture,
		Socket:             dto.Socket,
		BaseFrequency:      dto.BaseFrequency,
		TurboFrequency:     dto.TurboFrequency,
		Cores:              dto.Cores,
		Threads:            dto.Threads,
		L1Cache:            dto.L1Cache,
		L2Cache:            dto.L2Cache,
		L3Cache:            dto.L3Cache,
		Lithography:        dto.Lithography,
		TDP:                dto.TDP,
		Features:           dto.Features,
		MemoryType:         dto.MemoryType,
		MaxRAM:             dto.MaxRAM,
		MaxRAMFrequency:    dto.MaxRAMFrequency,
		IntegratedGraphics: dto.IntegratedGraphics,
		GraphicsModel:      dto.GraphicsModel,
		MaxTemperature:     dto.MaxTemperature,
		PackageContents:    dto.PackageContents,
		CountryOfOrigin:    dto.CountryOfOrigin,
	}

	// сохраняем Processor в БД
	if err := s.procRepo.CreateProcessor(processor); err != nil {
		return nil, err
	}

	// загружаем изображения в S3 и сохраняем ссылки в БД
	for _, fileHeader := range dto.Images {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		// Создаём временный файл
		tmpFile, err := os.CreateTemp("", "upload-*")
		if err != nil {
			continue
		}
		defer os.Remove(tmpFile.Name()) // удалить после загрузки

		// Копируем содержимое multipart.File в tmpFile
		if _, err := io.Copy(tmpFile, file); err != nil {
			tmpFile.Close()
			continue
		}
		tmpFile.Close()

		s3Key := fmt.Sprintf("processors/%s/%s", processor.ID, fileHeader.Filename)
		url, err := s.storage.Upload(context.Background(), s3Key, tmpFile.Name()) // передаём путь к файлу
		if err != nil {
			continue
		}

		// Сохраняем URL в БД
		img := &models.Image{
			ID:          uuid.New(),
			ProcessorID: &processor.ID,
			URL:         url,
		}
		err = s.procRepo.CreateImage(img)
		if err != nil {
			continue
		}
	}

	return processor, nil
}

func (s *ProcessorService) DeleteProcessor(procID uuid.UUID) error {
	return s.procRepo.GetDB().Transaction(func(tx *gorm.DB) error {

		// 1. Удаляем процессор из всех корзин
		if err := tx.
			Where("product_id = ? AND product_type = ?", procID, types.Processor).
			Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		// 2. Удаляем сам процессор
		if err := tx.
			Where("id = ?", procID).
			Delete(&models.Processor{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *ProcessorService) GetAllProcessors(filter dto.ProcessorFilterDTO) ([]dto.AllProcessorsResponseDTO, error) {
	return s.procRepo.GetProcessorsByFilter(filter)
}

func (s *ProcessorService) GetProcessorById(procID uuid.UUID) (*dto.ProcessorWithImagesDTO, error) {
	totalModel, err := s.procRepo.GetProcessorById(procID)
	if err != nil {
		return nil, err
	}
	totalModel.CountOrders, err = s.procRepo.CountOrders(procID)
	if err != nil {
		return nil, err
	}
	return totalModel, nil
}
func (s *ProcessorService) UpdateProcessor(procID uuid.UUID, procDto dto.ProcUpdate) error {
	ctx := context.Background()

	// 1) Обновляем поля процессора (без картинок)
	if err := s.procRepo.Update(procID, procDto); err != nil {
		return err
	}

	// helper: из публичного URL -> objectName в MinIO (то, что лежит после /{bucket}/)
	// пример URL: http://minio:9000/products/processors/<id>/222.jpg
	objectNameFromURL := func(u string) (string, bool) {
		marker := "/" + s.storage.Bucket + "/"
		i := strings.Index(u, marker)
		if i == -1 {
			return "", false
		}
		obj := u[i+len(marker):]
		if obj == "" {
			return "", false
		}
		return obj, true
	}

	// 2) Удаляем старые картинки, которые пользователь убрал (не входят в keep_image_urls)
	// Если keep_image_urls пустой — удалит все старые (логично, если юзер удалил всё).
	currentImages, err := s.procRepo.GetImagesByProcessorID(procID)
	if err != nil {
		return err
	}

	keepSet := make(map[string]struct{}, len(procDto.KeepImageURLs))
	for _, u := range procDto.KeepImageURLs {
		keepSet[u] = struct{}{}
	}

	for _, img := range currentImages {
		if _, ok := keepSet[img.URL]; ok {
			continue // оставляем
		}

		// удалить объект из MinIO
		if objName, ok := objectNameFromURL(img.URL); ok {
			_ = s.storage.Delete(ctx, objName) // не валим весь апдейт, даже если object уже удалён
		}

		// удалить запись из БД
		_ = s.procRepo.DeleteImageByID(img.ID)
	}

	// 3) Добавляем новые изображения (если пришли)
	for _, fileHeader := range procDto.Images {
		f, err := fileHeader.Open()
		if err != nil {
			continue
		}

		tmpFile, err := os.CreateTemp("", "upload-*")
		if err != nil {
			_ = f.Close()
			continue
		}

		// копируем в tmp
		if _, err := io.Copy(tmpFile, f); err != nil {
			_ = tmpFile.Close()
			_ = os.Remove(tmpFile.Name())
			_ = f.Close()
			continue
		}

		_ = tmpFile.Close()
		_ = f.Close()

		// ключ в MinIO (bucket = products)
		s3Key := fmt.Sprintf("processors/%s/%s", procID, fileHeader.Filename)

		url, err := s.storage.Upload(ctx, s3Key, tmpFile.Name())
		_ = os.Remove(tmpFile.Name())
		if err != nil {
			continue
		}

		img := &models.Image{
			ID:          uuid.New(),
			ProcessorID: &procID,
			URL:         url,
		}
		_ = s.procRepo.CreateImage(img)
	}

	return nil
}
