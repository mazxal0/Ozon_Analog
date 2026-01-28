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

type FlashDriveService struct {
	repo    *repository.FlashDriveRepository
	storage *storage.MinioStorage
}

func NewFlashDriveService(repo *repository.FlashDriveRepository, storage *storage.MinioStorage) *FlashDriveService {
	return &FlashDriveService{
		repo:    repo,
		storage: storage,
	}
}

func (s *FlashDriveService) DB() *repository.FlashDriveRepository {
	return s.repo
}

//
// CREATE
//

func (s *FlashDriveService) CreateFlashDrive(dto dto.FlashDriveCreateDTO) (*models.FlashDrive, error) {
	fd := &models.FlashDrive{
		ID:              uuid.New(),
		SKU:             utils.GenerateSKU(),
		Name:            dto.Name,
		Brand:           dto.Brand,
		RetailPrice:     dto.RetailPrice,
		WholesalePrice:  dto.WholesalePrice,
		WholesaleMinQty: dto.WholesaleMinQty,
		Stock:           dto.Stock,

		CapacityGB:      dto.CapacityGB,
		USBInterface:    dto.USBInterface,
		FormFactor:      dto.FormFactor,
		ReadSpeed:       dto.ReadSpeed,
		WriteSpeed:      dto.WriteSpeed,
		ChipType:        dto.ChipType,
		OTGSupport:      dto.OTGSupport,
		BodyMaterial:    dto.BodyMaterial,
		Color:           dto.Color,
		WaterResistance: dto.WaterResistance,
		DustResistance:  dto.DustResistance,
		Shockproof:      dto.Shockproof,
		CapType:         dto.CapType,

		LengthMM:    dto.LengthMM,
		WidthMM:     dto.WidthMM,
		ThicknessMM: dto.ThicknessMM,
		WeightG:     dto.WeightG,

		Compatibility:   dto.Compatibility,
		OperatingTemp:   dto.OperatingTemp,
		StorageTemp:     dto.StorageTemp,
		CountryOfOrigin: dto.CountryOfOrigin,
		PackageContents: dto.PackageContents,
		WarrantyMonths:  dto.WarrantyMonths,
		Features:        dto.Features,
	}

	// сохраняем сам продукт
	if err := s.repo.CreateFlashDrive(fd); err != nil {
		return nil, err
	}

	// загружаем изображения
	for _, fileHeader := range dto.ImageFiles {
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

		s3Key := fmt.Sprintf("flashdrives/%s/%s", fd.ID, fileHeader.Filename)
		url, err := s.storage.Upload(context.Background(), s3Key, tmpFile.Name())
		if err != nil {
			continue
		}

		img := &models.Image{
			ID:           uuid.New(),
			FlashDriveID: &fd.ID,
			URL:          url,
		}

		err = s.repo.CreateImage(img)

		if err != nil {
			return nil, err
		}
	}

	return fd, nil
}

//
// GET ALL
//

func (s *FlashDriveService) GetAllFlashDrives(filter dto.FlashDriveFilterDTO) ([]dto.AllFlashDrivesResponseDTO, error) {
	return s.repo.GetFlashDrivesByFilter(filter)
}

//
// GET BY ID
//

func (s *FlashDriveService) GetFlashDriveById(id uuid.UUID) (*dto.FlashDriveWithImagesDTO, error) {
	totalModel, err := s.repo.GetFlashDriveById(id)
	if err != nil {
		return nil, err
	}
	totalModel.CountOrders, err = s.repo.CountOrders(id)
	if err != nil {
		return nil, err
	}
	return totalModel, nil
}

//
// DELETE
//

func (s *FlashDriveService) DeleteFlashDrive(id uuid.UUID) error {
	return s.repo.GetDB().Transaction(func(tx *gorm.DB) error {

		// 1. Удаляем FlashDrive из всех корзин
		if err := tx.
			Where(
				"product_id = ? AND product_type = ?",
				id,
				types.FlashDriver,
			).
			Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		// 2. Удаляем сам FlashDrive
		if err := tx.
			Where("id = ?", id).
			Delete(&models.FlashDrive{}).Error; err != nil {
			return err
		}

		return nil
	})
}

//
// UPDATE
//

func (s *FlashDriveService) UpdateFlashDrive(id uuid.UUID, upd dto.FlashDriveUpdateDTO) error {
	ctx := context.Background()

	// 1) обновляем поля в БД
	if err := s.repo.Update(id, upd); err != nil {
		return err
	}

	// helper: URL -> objectName (после /{bucket}/)
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

	// 2) удаляем старые изображения, которых нет в keep_image_urls
	currentImages, err := s.repo.GetImagesByFlashDriveID(id)
	if err != nil {
		return err
	}

	keepSet := make(map[string]struct{}, len(upd.KeepImageURLs))
	for _, u := range upd.KeepImageURLs {
		keepSet[u] = struct{}{}
	}

	for _, img := range currentImages {
		if _, ok := keepSet[img.URL]; ok {
			continue
		}

		if objName, ok := objectNameFromURL(img.URL); ok {
			_ = s.storage.Delete(ctx, objName)
		}
		_ = s.repo.DeleteImageByID(img.ID)
	}

	// 3) загружаем новые изображения
	for _, fileHeader := range upd.ImageFiles {
		f, err := fileHeader.Open()
		if err != nil {
			continue
		}

		tmpFile, err := os.CreateTemp("", "upload-*")
		if err != nil {
			_ = f.Close()
			continue
		}

		if _, err = io.Copy(tmpFile, f); err != nil {
			_ = tmpFile.Close()
			_ = os.Remove(tmpFile.Name())
			_ = f.Close()
			continue
		}

		_ = tmpFile.Close()
		_ = f.Close()

		s3Key := fmt.Sprintf("flashdrives/%s/%s", id, fileHeader.Filename)
		url, err := s.storage.Upload(ctx, s3Key, tmpFile.Name())
		_ = os.Remove(tmpFile.Name())
		if err != nil {
			continue
		}

		img := &models.Image{
			ID:           uuid.New(),
			FlashDriveID: &id,
			URL:          url,
		}
		_ = s.repo.CreateImage(img)
	}

	return nil
}
