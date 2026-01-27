package repository

import (
	"Market_backend/internal/common"
	"Market_backend/internal/common/types"
	"Market_backend/internal/product/dto"
	"Market_backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FlashDriveRepository struct {
	db *gorm.DB
}

func NewFlashDriveRepository() *FlashDriveRepository {
	return &FlashDriveRepository{db: common.DB}
}

func (r *FlashDriveRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *FlashDriveRepository) CreateFlashDrive(fd *models.FlashDrive) error {
	return r.db.Create(&fd).Error
}

func (r *FlashDriveRepository) DeleteFlashDrive(fdID uuid.UUID) error {
	return r.db.Where("id = ?", fdID).Delete(&models.FlashDrive{}).Error
}

// --------------------------------------------------------------
// Получение списка флешек с фильтрами
// --------------------------------------------------------------
func (r *FlashDriveRepository) GetFlashDrivesByFilter(filter dto.FlashDriveFilterDTO) ([]dto.AllFlashDrivesResponseDTO, error) {

	subQuery := r.db.
		Table("images").
		Select("url").
		Where("images.flash_drive_id = flash_drives.id").
		Order("images.created_at ASC").
		Limit(1)

	query := r.db.
		Table("flash_drives").
		Select(`
            flash_drives.id,
            flash_drives.name,
            flash_drives.brand,
            flash_drives.capacity_gb,
            flash_drives.usb_interface,
            flash_drives.retail_price,
            flash_drives.wholesale_price,
            (?) AS image_url
        `, subQuery)

	// ---- Filters ----
	if len(filter.Brands) > 0 && filter.Brands[0] != "" {
		query = query.Where("flash_drives.brand IN ?", filter.Brands)
	}

	if len(filter.CapacityGB) > 0 {
		query = query.Where("flash_drives.capacity_gb IN ?", filter.CapacityGB)
	}

	if len(filter.USBInterface) > 0 && filter.USBInterface[0] != "" {
		query = query.Where("flash_drives.usb_interface IN ?", filter.USBInterface)
	}

	if filter.PriceAsc {
		query = query.Order("flash_drives.retail_price ASC")
	} else {
		query = query.Order("flash_drives.retail_price DESC")
	}

	query = query.Limit(filter.Limit).Offset(filter.Offset)

	var drives []dto.AllFlashDrivesResponseDTO
	if err := query.Scan(&drives).Error; err != nil {
		return nil, err
	}

	return drives, nil
}

// --------------------------------------------------------------
// Получение флешки по ID + все изображения
// --------------------------------------------------------------
func (r *FlashDriveRepository) GetFlashDriveById(fdID uuid.UUID) (*dto.FlashDriveWithImagesDTO, error) {
	var fd models.FlashDrive
	if err := r.db.Preload("Images").First(&fd, "id = ?", fdID).Error; err != nil {
		return nil, err
	}

	// Собираем URL изображений
	var urls []string
	for _, img := range fd.Images {
		urls = append(urls, img.URL)
	}

	return &dto.FlashDriveWithImagesDTO{
		ID:              fd.ID,
		SKU:             fd.SKU,
		Name:            fd.Name,
		Brand:           fd.Brand,
		RetailPrice:     fd.RetailPrice,
		WholesalePrice:  fd.WholesalePrice,
		WholesaleMinQty: fd.WholesaleMinQty,
		Stock:           fd.Stock,

		CapacityGB:      fd.CapacityGB,
		USBInterface:    fd.USBInterface,
		FormFactor:      fd.FormFactor,
		ReadSpeed:       fd.ReadSpeed,
		WriteSpeed:      fd.WriteSpeed,
		ChipType:        fd.ChipType,
		OTGSupport:      fd.OTGSupport,
		BodyMaterial:    fd.BodyMaterial,
		Color:           fd.Color,
		WaterResistance: fd.WaterResistance,
		DustResistance:  fd.DustResistance,
		Shockproof:      fd.Shockproof,
		CapType:         fd.CapType,

		LengthMM:    fd.LengthMM,
		WidthMM:     fd.WidthMM,
		ThicknessMM: fd.ThicknessMM,
		WeightG:     fd.WeightG,

		Compatibility:   fd.Compatibility,
		OperatingTemp:   fd.OperatingTemp,
		StorageTemp:     fd.StorageTemp,
		CountryOfOrigin: fd.CountryOfOrigin,
		PackageContents: fd.PackageContents,
		WarrantyMonths:  fd.WarrantyMonths,
		Features:        fd.Features,

		ImageURLs: urls,
	}, nil
}

// Version для транзакций
func (r *FlashDriveRepository) GetFlashDriveByIdTx(tx *gorm.DB, fdID uuid.UUID) (*dto.FlashDriveWithImagesDTO, error) {
	var fd models.FlashDrive
	if err := tx.Preload("Images").First(&fd, "id = ?", fdID).Error; err != nil {
		return nil, err
	}

	var urls []string
	for _, img := range fd.Images {
		urls = append(urls, img.URL)
	}

	return &dto.FlashDriveWithImagesDTO{
		ID:              fd.ID,
		SKU:             fd.SKU,
		Name:            fd.Name,
		Brand:           fd.Brand,
		RetailPrice:     fd.RetailPrice,
		WholesalePrice:  fd.WholesalePrice,
		WholesaleMinQty: fd.WholesaleMinQty,
		Stock:           fd.Stock,

		CapacityGB:      fd.CapacityGB,
		USBInterface:    fd.USBInterface,
		FormFactor:      fd.FormFactor,
		ReadSpeed:       fd.ReadSpeed,
		WriteSpeed:      fd.WriteSpeed,
		ChipType:        fd.ChipType,
		OTGSupport:      fd.OTGSupport,
		BodyMaterial:    fd.BodyMaterial,
		Color:           fd.Color,
		WaterResistance: fd.WaterResistance,
		DustResistance:  fd.DustResistance,
		Shockproof:      fd.Shockproof,
		CapType:         fd.CapType,

		LengthMM:    fd.LengthMM,
		WidthMM:     fd.WidthMM,
		ThicknessMM: fd.ThicknessMM,
		WeightG:     fd.WeightG,

		Compatibility:   fd.Compatibility,
		OperatingTemp:   fd.OperatingTemp,
		StorageTemp:     fd.StorageTemp,
		CountryOfOrigin: fd.CountryOfOrigin,
		PackageContents: fd.PackageContents,
		WarrantyMonths:  fd.WarrantyMonths,
		Features:        fd.Features,

		ImageURLs: urls,
	}, nil
}

// --------------------------------------------------------------
// Обновление
// --------------------------------------------------------------
func (r *FlashDriveRepository) Update(fdID uuid.UUID, fd dto.FlashDriveUpdateDTO) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		updateData := map[string]interface{}{
			"name":              fd.Name,
			"brand":             fd.Brand,
			"retail_price":      fd.RetailPrice,
			"wholesale_price":   fd.WholesalePrice,
			"wholesale_min_qty": fd.WholesaleMinQty,
			"stock":             fd.Stock,

			"capacity_gb":      fd.CapacityGB,
			"usb_interface":    fd.USBInterface,
			"form_factor":      fd.FormFactor,
			"read_speed":       fd.ReadSpeed,
			"write_speed":      fd.WriteSpeed,
			"chip_type":        fd.ChipType,
			"otg_support":      fd.OTGSupport,
			"body_material":    fd.BodyMaterial,
			"color":            fd.Color,
			"water_resistance": fd.WaterResistance,
			"dust_resistance":  fd.DustResistance,
			"shockproof":       fd.Shockproof,
			"cap_type":         fd.CapType,

			"length_mm":    fd.LengthMM,
			"width_mm":     fd.WidthMM,
			"thickness_mm": fd.ThicknessMM,
			"weight_g":     fd.WeightG,

			"compatibility":     fd.Compatibility,
			"operating_temp":    fd.OperatingTemp,
			"storage_temp":      fd.StorageTemp,
			"country_of_origin": fd.CountryOfOrigin,
			"package_contents":  fd.PackageContents,
			"warranty_months":   fd.WarrantyMonths,
			"features":          fd.Features,
		}

		res := tx.Model(&models.FlashDrive{}).Where("id = ?", fdID).Updates(updateData)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

// --------------------------------------------------------------
// Работа с изображениями
// --------------------------------------------------------------
func (r *FlashDriveRepository) CreateImage(image *models.Image) error {
	return r.db.Create(image).Error
}

func (r *FlashDriveRepository) GetImagesByFlashDriveID(fdID uuid.UUID) ([]models.Image, error) {
	var images []models.Image
	err := r.db.Where("flash_drive_id = ?", fdID).Find(&images).Error
	return images, err
}

func (r *FlashDriveRepository) DeleteImages(fdID uuid.UUID) error {
	return r.db.Where("flash_drive_id = ?", fdID).Delete(&models.Image{}).Error
}

func (r *FlashDriveRepository) AddImages(fdID uuid.UUID, urls []string) error {
	for _, url := range urls {
		img := models.Image{
			ID:           uuid.New(),
			FlashDriveID: &fdID,
			URL:          url,
		}
		if err := r.db.Create(&img).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *FlashDriveRepository) CountOrders(flashID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&models.OrderItem{}).
		Where("product_id = ? AND product_type = ?", flashID, types.FlashDriver).
		Distinct("order_id"). // учитываем только уникальные заказы
		Count(&count).Error
	return int(count), err
}

func (r *FlashDriveRepository) UpdateStock(flashID uuid.UUID, newStock int) error {
	res := r.db.Model(&models.FlashDrive{}).Where("id = ?", flashID).Update("stock", newStock)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
