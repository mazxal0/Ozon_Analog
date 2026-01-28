package repository

import (
	"Market_backend/internal/common"
	"Market_backend/internal/common/types"
	"Market_backend/internal/product/dto"
	"Market_backend/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProcessorRepository struct {
	db *gorm.DB
}

func (r *ProcessorRepository) GetDB() *gorm.DB {
	return r.db
}

func NewProcessorRepository() *ProcessorRepository {
	return &ProcessorRepository{db: common.DB}
}

func (r *ProcessorRepository) CreateProcessor(proc *models.Processor) error {
	return r.db.Create(&proc).Error
}

func (r *ProcessorRepository) DeleteProcessor(procID uuid.UUID) error {
	return r.db.Where("id = ?", procID).Delete(&models.Processor{}).Error
}

func (r *ProcessorRepository) GetProcessorsByFilter(filter dto.ProcessorFilterDTO) ([]dto.AllProcessorsResponseDTO, error) {
	var result []dto.AllProcessorsResponseDTO

	db := r.db.Table("processors p").
		Select(`
			p.id,
			p.name,
			p.retail_price,
			p.wholesale_price,
			i.url AS image_url
		`).
		Joins(`
			LEFT JOIN LATERAL (
				SELECT url
				FROM images
				WHERE images.processor_id = p.id
				ORDER BY created_at ASC
				LIMIT 1
			) i ON true
		`)

	// фильтры
	if len(filter.Brands) > 0 && !(len(filter.Brands) == 1 && filter.Brands[0] == "") {
		db = db.Where("p.brand IN ?", filter.Brands)
	}
	if len(filter.Frequencies) > 0 {
		db = db.Where("p.base_frequency IN ?", filter.Frequencies)
	}
	if len(filter.Cores) > 0 {
		db = db.Where("p.cores IN ?", filter.Cores)
	}

	// сортировка
	if filter.PriceAsc {
		db = db.Order("p.retail_price ASC")
	} else {
		db = db.Order("p.retail_price DESC")
	}

	// пагинация
	db = db.Limit(filter.Limit).Offset(filter.Offset)

	// выполняем запрос
	if err := db.Scan(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ProcessorRepository) GetProcessorById(procId uuid.UUID) (*dto.ProcessorWithImagesDTO, error) {
	var proc models.Processor
	if err := r.db.Preload("Images").First(&proc, "id = ?", procId).Error; err != nil {
		return nil, err
	}

	// собираем URL изображений
	var urls []string
	for _, img := range proc.Images {
		urls = append(urls, img.URL)
	}

	// возвращаем DTO со всеми полями + ImageURLs
	return &dto.ProcessorWithImagesDTO{
		ID:                 proc.ID,
		SKU:                proc.SKU,
		Name:               proc.Name,
		Brand:              proc.Brand,
		RetailPrice:        proc.RetailPrice,
		WholesalePrice:     proc.WholesalePrice,
		WholesaleMinQty:    proc.WholesaleMinQty,
		Stock:              proc.Stock,
		Line:               proc.Line,
		Architecture:       proc.Architecture,
		Socket:             proc.Socket,
		BaseFrequency:      proc.BaseFrequency,
		TurboFrequency:     proc.TurboFrequency,
		Cores:              proc.Cores,
		Threads:            proc.Threads,
		L1Cache:            proc.L1Cache,
		L2Cache:            proc.L2Cache,
		L3Cache:            proc.L3Cache,
		Lithography:        proc.Lithography,
		TDP:                proc.TDP,
		Features:           proc.Features,
		MemoryType:         proc.MemoryType,
		MaxRAM:             proc.MaxRAM,
		MaxRAMFrequency:    proc.MaxRAMFrequency,
		IntegratedGraphics: proc.IntegratedGraphics,
		GraphicsModel:      proc.GraphicsModel,
		MaxTemperature:     proc.MaxTemperature,
		PackageContents:    proc.PackageContents,
		CountryOfOrigin:    proc.CountryOfOrigin,
		ImageURLs:          urls,
	}, nil
}

func (r *ProcessorRepository) GetProcessorByIdTx(tx *gorm.DB, procId uuid.UUID) (*dto.ProcessorWithImagesDTO, error) {
	var proc models.Processor
	if err := tx.Preload("Images").First(&proc, "id = ?", procId).Error; err != nil {
		return nil, err
	}

	// собираем URL изображений
	var urls []string
	for _, img := range proc.Images {
		urls = append(urls, img.URL)
	}

	// возвращаем DTO со всеми полями + ImageURLs
	return &dto.ProcessorWithImagesDTO{
		ID:                 proc.ID,
		SKU:                proc.SKU,
		Name:               proc.Name,
		Brand:              proc.Brand,
		RetailPrice:        proc.RetailPrice,
		WholesalePrice:     proc.WholesalePrice,
		WholesaleMinQty:    proc.WholesaleMinQty,
		Stock:              proc.Stock,
		Line:               proc.Line,
		Architecture:       proc.Architecture,
		Socket:             proc.Socket,
		BaseFrequency:      proc.BaseFrequency,
		TurboFrequency:     proc.TurboFrequency,
		Cores:              proc.Cores,
		Threads:            proc.Threads,
		L1Cache:            proc.L1Cache,
		L2Cache:            proc.L2Cache,
		L3Cache:            proc.L3Cache,
		Lithography:        proc.Lithography,
		TDP:                proc.TDP,
		Features:           proc.Features,
		MemoryType:         proc.MemoryType,
		MaxRAM:             proc.MaxRAM,
		MaxRAMFrequency:    proc.MaxRAMFrequency,
		IntegratedGraphics: proc.IntegratedGraphics,
		GraphicsModel:      proc.GraphicsModel,
		MaxTemperature:     proc.MaxTemperature,
		PackageContents:    proc.PackageContents,
		CountryOfOrigin:    proc.CountryOfOrigin,
		ImageURLs:          urls,
	}, nil
}

func (r *ProcessorRepository) Update(procId uuid.UUID, proc dto.ProcUpdate) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// обновляем только поля Processor, GORM сам определит, какие поля есть
		updateData := map[string]interface{}{
			"name":                proc.Name,
			"brand":               proc.Brand,
			"retail_price":        proc.RetailPrice,
			"wholesale_price":     proc.WholesalePrice,
			"wholesale_min_qty":   proc.WholesaleMinQty,
			"stock":               proc.Stock,
			"line":                proc.Line,
			"architecture":        proc.Architecture,
			"socket":              proc.Socket,
			"base_frequency":      proc.BaseFrequency,
			"turbo_frequency":     proc.TurboFrequency,
			"cores":               proc.Cores,
			"threads":             proc.Threads,
			"l1_cache":            proc.L1Cache,
			"l2_cache":            proc.L2Cache,
			"l3_cache":            proc.L3Cache,
			"lithography":         proc.Lithography,
			"tdp":                 proc.TDP,
			"features":            proc.Features,
			"memory_type":         proc.MemoryType,
			"max_ram":             proc.MaxRAM,
			"max_ram_frequency":   proc.MaxRAMFrequency,
			"integrated_graphics": proc.IntegratedGraphics,
			"graphics_model":      proc.GraphicsModel,
			"max_temperature":     proc.MaxTemperature,
			"package_contents":    proc.PackageContents,
			"country_of_origin":   proc.CountryOfOrigin,
		}

		res := tx.Model(&models.Processor{}).Where("id = ?", procId).Updates(updateData)
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil
	})
}

// Сохраняем изображение
func (r *ProcessorRepository) CreateImage(image *models.Image) error {
	return r.db.Create(image).Error
}

// Получение всех изображений процессора (если нужно)
func (r *ProcessorRepository) GetImagesByProcessorID(processorID uuid.UUID) ([]models.Image, error) {
	var images []models.Image
	err := r.db.Where("processor_id = ?", processorID).Find(&images).Error
	return images, err
}

// Удаляем все изображения процессора
func (r *ProcessorRepository) DeleteImages(procID uuid.UUID) error {
	return r.db.Where("processor_id = ?", procID).Delete(&models.Image{}).Error
}

// Добавляем новые изображения
func (r *ProcessorRepository) AddImages(procID uuid.UUID, urls []string) error {
	for _, url := range urls {
		img := models.Image{
			ID:          uuid.New(),
			ProcessorID: &procID,
			URL:         url,
		}
		if err := r.db.Create(&img).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *ProcessorRepository) CountOrders(procID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&models.OrderItem{}).
		Where("product_id = ? AND product_type = ?", procID, types.Processor).
		Distinct("order_id"). // учитываем только уникальные заказы
		Count(&count).Error
	return int(count), err
}

func (r *ProcessorRepository) UpdateStock(processorID uuid.UUID, newStock int) error {
	res := r.db.Model(&models.Processor{}).Where("id = ?", processorID).Update("stock", newStock)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ProcessorRepository) DeleteImageByID(imageID uuid.UUID) error {
	res := r.db.Where("id = ?", imageID).Delete(&models.Image{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
