package handler

import (
	"Market_backend/internal/common/utils"
	"Market_backend/internal/product/dto"
	"Market_backend/internal/product/service"
	"errors"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FlashDriveHandler struct {
	service *service.FlashDriveService
}

func NewFlashDriveHandler(service *service.FlashDriveService) *FlashDriveHandler {
	return &FlashDriveHandler{service: service}
}

func (h *FlashDriveHandler) CreateFlashDrive(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form"})
	}

	files := form.File["images"]
	if len(files) > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "max 5 images allowed"})
	}

	get := func(key string) string {
		if vals, ok := form.Value[key]; ok && len(vals) > 0 {
			return vals[0]
		}
		return ""
	}

	dtoFD := dto.FlashDriveCreateDTO{
		SKU:             get("sku"),
		Name:            get("name"),
		Brand:           get("brand"),
		RetailPrice:     utils.ParseFloat(get("retail_price")),
		WholesalePrice:  utils.ParseFloat(get("wholesale_price")),
		WholesaleMinQty: utils.ParseInt(get("wholesale_min_qty")),
		Stock:           utils.ParseInt(get("stock")),

		CapacityGB:      utils.ParseInt(get("capacity_gb")),
		USBInterface:    get("usb_interface"),
		FormFactor:      get("form_factor"),
		ReadSpeed:       utils.ParseInt(get("read_speed")),
		WriteSpeed:      utils.ParseInt(get("write_speed")),
		ChipType:        get("chip_type"),
		OTGSupport:      utils.ParseBool(get("otg_support")),
		BodyMaterial:    get("body_material"),
		Color:           get("color"),
		WaterResistance: utils.ParseBool(get("water_resistance")),
		DustResistance:  utils.ParseBool(get("dust_resistance")),
		Shockproof:      utils.ParseBool(get("shockproof")),
		CapType:         get("cap_type"),

		LengthMM:    utils.ParseFloat(get("length_mm")),
		WidthMM:     utils.ParseFloat(get("width_mm")),
		ThicknessMM: utils.ParseFloat(get("thickness_mm")),
		WeightG:     utils.ParseFloat(get("weight_g")),

		Compatibility:   get("compatibility"),
		OperatingTemp:   get("operating_temp"),
		StorageTemp:     get("storage_temp"),
		CountryOfOrigin: get("country_of_origin"),
		PackageContents: get("package_contents"),
		WarrantyMonths:  utils.ParseInt(get("warranty_months")),
		Features:        get("features"),

		ImageFiles: files,
	}

	fd, err := h.service.CreateFlashDrive(dtoFD)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"flash_drive": fd})
}

func (h *FlashDriveHandler) DeleteFlashDrive(c *fiber.Ctx) error {
	idStr := c.Params("flashId")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	err = h.service.DeleteFlashDrive(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *FlashDriveHandler) GetAllFlashDrives(c *fiber.Ctx) error {
	// Разбираем бренды
	brands := strings.Split(c.Query("brands", ""), ",")

	// Разбираем объемы памяти
	capacityStr := strings.Split(c.Query("capacities", ""), ",")
	var capacities []int
	for _, s := range capacityStr {
		if n, err := strconv.Atoi(s); err == nil {
			capacities = append(capacities, n)
		}
	}

	// Разбираем интерфейсы USB
	usbInterfaces := strings.Split(c.Query("usb_interfaces", ""), ",")

	// Лимит и оффсет
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	// Сортировка по цене
	priceAsc := c.QueryBool("price_asc")

	// Собираем DTO для фильтра
	filter := dto.FlashDriveFilterDTO{
		Brands:       brands,
		CapacityGB:   capacities,
		USBInterface: usbInterfaces,
		PriceAsc:     priceAsc,
		Limit:        limit,
		Offset:       offset,
	}

	// Получаем данные через сервис
	list, err := h.service.GetAllFlashDrives(filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"flash_drives": list})
}

func (h *FlashDriveHandler) GetFlashDriveById(c *fiber.Ctx) error {
	idStr := c.Params("flashId")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	fd, err := h.service.GetFlashDriveById(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"flash_drive": fd})
}

func (h *FlashDriveHandler) UpdateFlashDrive(c *fiber.Ctx) error {
	idStr := c.Params("flashId")
	flashID, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	form, err := c.MultipartForm()
	if err != nil && errors.Is(err, fiber.ErrUnprocessableEntity) {
		return c.Status(400).JSON(fiber.Map{"error": "invalid form"})
	}

	var files []*multipart.FileHeader
	if form != nil {
		files = form.File["images"]
	}

	get := func(key string) string {
		if form != nil {
			if vals, ok := form.Value[key]; ok && len(vals) > 0 {
				return vals[0]
			}
		}
		return ""
	}

	dtoFD := dto.FlashDriveUpdateDTO{
		Name:            get("name"),
		Brand:           get("brand"),
		RetailPrice:     utils.ParseFloat(get("retail_price")),
		WholesalePrice:  utils.ParseFloat(get("wholesale_price")),
		WholesaleMinQty: utils.ParseInt(get("wholesale_min_qty")),
		Stock:           utils.ParseInt(get("stock")),

		CapacityGB:      utils.ParseInt(get("capacity_gb")),
		USBInterface:    get("usb_interface"),
		FormFactor:      get("form_factor"),
		ReadSpeed:       utils.ParseInt(get("read_speed")),
		WriteSpeed:      utils.ParseInt(get("write_speed")),
		ChipType:        get("chip_type"),
		OTGSupport:      utils.ParseBool(get("otg_support")),
		BodyMaterial:    get("body_material"),
		Color:           get("color"),
		WaterResistance: utils.ParseBool(get("water_resistance")),
		DustResistance:  utils.ParseBool(get("dust_resistance")),
		Shockproof:      utils.ParseBool(get("shockproof")),
		CapType:         get("cap_type"),

		LengthMM:    utils.ParseFloat(get("length_mm")),
		WidthMM:     utils.ParseFloat(get("width_mm")),
		ThicknessMM: utils.ParseFloat(get("thickness_mm")),
		WeightG:     utils.ParseFloat(get("weight_g")),

		Compatibility:   get("compatibility"),
		OperatingTemp:   get("operating_temp"),
		StorageTemp:     get("storage_temp"),
		CountryOfOrigin: get("country_of_origin"),
		PackageContents: get("package_contents"),
		WarrantyMonths:  utils.ParseInt(get("warranty_months")),
		Features:        get("features"),

		ImageFiles: files,
	}

	err = h.service.UpdateFlashDrive(flashID, dtoFD)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "success"})
}
