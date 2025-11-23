package service

import (
	"context"
	"eduVix_backend/internal/common/utils"
	"eduVix_backend/internal/product/dto"
	"eduVix_backend/internal/product/repository"
	"eduVix_backend/internal/storage"
	"eduVix_backend/models"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
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
	return s.procRepo.DeleteProcessor(procID)
}

func (s *ProcessorService) GetAllProcessors(filter dto.ProcessorFilterDTO) ([]dto.AllProcessorsResponseDTO, error) {
	return s.procRepo.GetProcessorsByFilter(filter)
}

func (s *ProcessorService) GetProcessorById(procID uuid.UUID) (*dto.ProcessorWithImagesDTO, error) {
	return s.procRepo.GetProcessorById(procID)
}

func (s *ProcessorService) UpdateProcessor(procID uuid.UUID, procDto dto.ProcUpdate) error {
	// Получаем существующий процессор
	//processor, err := s.procRepo.GetProcessorById(procID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//// Обновляем поля
	//processor.Name = dto.Name
	//processor.Brand = dto.Brand
	//processor.RetailPrice = dto.RetailPrice
	//processor.WholesalePrice = dto.WholesalePrice
	//processor.WholesaleMinQty = dto.WholesaleMinQty
	//processor.Stock = dto.Stock
	//processor.Line = dto.Line
	//processor.Architecture = dto.Architecture
	//processor.Socket = dto.Socket
	//processor.BaseFrequency = dto.BaseFrequency
	//processor.TurboFrequency = dto.TurboFrequency
	//processor.Cores = dto.Cores
	//processor.Threads = dto.Threads
	//processor.L1Cache = dto.L1Cache
	//processor.L2Cache = dto.L2Cache
	//processor.L3Cache = dto.L3Cache
	//processor.Lithography = dto.Lithography
	//processor.TDP = dto.TDP
	//processor.Features = dto.Features
	//processor.MemoryType = dto.MemoryType
	//processor.MaxRAM = dto.MaxRAM
	//processor.MaxRAMFrequency = dto.MaxRAMFrequency
	//processor.IntegratedGraphics = dto.IntegratedGraphics
	//processor.GraphicsModel = dto.GraphicsModel
	//processor.MaxTemperature = dto.MaxTemperature
	//processor.PackageContents = dto.PackageContents
	//processor.CountryOfOrigin = dto.CountryOfOrigin

	// Обновляем процессор в БД через репозиторий
	if err := s.procRepo.Update(procID, procDto); err != nil {
		return err
	}

	// Если пришли новые изображения (файлы)
	for _, fileHeader := range procDto.Images {
		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		tmpFile, err := os.CreateTemp("", "upload-*")
		if err != nil {
			continue
		}
		defer os.Remove(tmpFile.Name())

		if _, err := io.Copy(tmpFile, file); err != nil {
			tmpFile.Close()
			continue
		}
		tmpFile.Close()

		s3Key := fmt.Sprintf("processors/%s/%s", procID, fileHeader.Filename)
		url, err := s.storage.Upload(context.Background(), s3Key, tmpFile.Name())
		if err != nil {
			continue
		}

		img := &models.Image{
			ID:          uuid.New(),
			ProcessorID: &procID,
			URL:         url,
		}
		s.procRepo.CreateImage(img)
	}

	return nil
}
