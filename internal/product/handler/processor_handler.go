package handler

import (
	"Market_backend/internal/common/utils"
	"Market_backend/internal/product/dto"
	"Market_backend/internal/product/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"mime/multipart"
	"strconv"
	"strings"
)

type ProcessorHandler struct {
	service *service.ProcessorService
}

func NewProcessorHandler(service *service.ProcessorService) *ProcessorHandler {
	return &ProcessorHandler{service: service}
}

func (h *ProcessorHandler) CreateProcessor(c *fiber.Ctx) error {
	// парсим форму с multipart
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form"})
	}

	files := form.File["images"] // массив до 5 файлов
	if len(files) > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "max 5 images allowed"})
	}

	getValue := func(key string) string {
		if vals, ok := form.Value[key]; ok && len(vals) > 0 {
			return vals[0]
		}
		return ""
	}

	// создаем processor через сервис
	procDto := dto.ProcessorCreateDTO{
		Name:               getValue("name"),
		Brand:              getValue("brand"),
		RetailPrice:        utils.ParseFloat(getValue("retail_price")),
		WholesalePrice:     utils.ParseFloat(getValue("wholesale_price")),
		WholesaleMinQty:    utils.ParseInt(getValue("wholesale_min_qty")),
		Stock:              utils.ParseInt(getValue("stock")),
		Line:               getValue("line"),
		Architecture:       getValue("architecture"),
		Socket:             getValue("socket"),
		BaseFrequency:      utils.ParseFloat(getValue("base_frequency")),
		TurboFrequency:     utils.ParseFloat(getValue("turbo_frequency")),
		Cores:              utils.ParseInt(getValue("cores")),
		Threads:            utils.ParseInt(getValue("threads")),
		L1Cache:            getValue("l1_cache"),
		L2Cache:            getValue("l2_cache"),
		L3Cache:            getValue("l3_cache"),
		Lithography:        getValue("lithography"),
		TDP:                utils.ParseInt(getValue("tdp")),
		Features:           getValue("features"),
		MemoryType:         getValue("memory_type"),
		MaxRAM:             getValue("max_ram"),
		MaxRAMFrequency:    getValue("max_ram_frequency"),
		IntegratedGraphics: utils.ParseBool(getValue("integrated_graphics")),
		GraphicsModel:      getValue("graphics_model"),
		MaxTemperature:     utils.ParseInt(getValue("max_temperature")),
		PackageContents:    getValue("package_contents"),
		CountryOfOrigin:    getValue("country_of_origin"),
		Images:             files,
	}

	proc, err := h.service.CreateProcessor(procDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"processor": proc})
}

func (h *ProcessorHandler) DeleteProcessor(c *fiber.Ctx) error {
	procIDStr := c.Params("procId")
	if len(procIDStr) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	procID, err := uuid.Parse(procIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	err = h.service.DeleteProcessor(procID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"processor": nil})
}

func (h *ProcessorHandler) GetAllProcessors(c *fiber.Ctx) error {
	// получаем query-параметры
	brands := strings.Split(c.Query("brands", ""), ",")
	freqStrs := strings.Split(c.Query("frequencies", ""), ",")
	coresStrs := strings.Split(c.Query("cores", ""), ",")

	var frequencies []float64
	for _, f := range freqStrs {
		if v, err := strconv.ParseFloat(f, 64); err == nil {
			frequencies = append(frequencies, v)
		}
	}

	var cores []int
	for _, c := range coresStrs {
		if v, err := strconv.Atoi(c); err == nil {
			cores = append(cores, v)
		}
	}

	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	priceAsc := c.QueryBool("price_asc")

	filter := dto.ProcessorFilterDTO{
		Brands:      brands,
		Frequencies: frequencies,
		Cores:       cores,
		PriceAsc:    priceAsc,
		Limit:       limit,
		Offset:      offset,
	}

	processors, err := h.service.GetAllProcessors(filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"processors": processors})
}

func (h *ProcessorHandler) GetProcessorById(c *fiber.Ctx) error {
	procIDStr := c.Params("procId")
	if len(procIDStr) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	procID, err := uuid.Parse(procIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	proc, err := h.service.GetProcessorById(procID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"processor": proc})
}

func (h *ProcessorHandler) UpdateProcessor(c *fiber.Ctx) error {
	procIDStr := c.Params("procId")
	procID, err := uuid.Parse(procIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid processor ID"})
	}

	form, err := c.MultipartForm()
	if err != nil && err != fiber.ErrUnprocessableEntity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid form"})
	}

	var files []*multipart.FileHeader
	if form != nil {
		files = form.File["images"]
	}

	getValue := func(key string) string {
		if form != nil {
			if vals, ok := form.Value[key]; ok && len(vals) > 0 {
				return vals[0]
			}
		}
		return ""
	}

	procDto := dto.ProcUpdate{
		Name:               getValue("name"),
		Brand:              getValue("brand"),
		RetailPrice:        utils.ParseFloat(getValue("retail_price")),
		WholesalePrice:     utils.ParseFloat(getValue("wholesale_price")),
		WholesaleMinQty:    utils.ParseInt(getValue("wholesale_min_qty")),
		Stock:              utils.ParseInt(getValue("stock")),
		Line:               getValue("line"),
		Architecture:       getValue("architecture"),
		Socket:             getValue("socket"),
		BaseFrequency:      utils.ParseFloat(getValue("base_frequency")),
		TurboFrequency:     utils.ParseFloat(getValue("turbo_frequency")),
		Cores:              utils.ParseInt(getValue("cores")),
		Threads:            utils.ParseInt(getValue("threads")),
		L1Cache:            getValue("l1_cache"),
		L2Cache:            getValue("l2_cache"),
		L3Cache:            getValue("l3_cache"),
		Lithography:        getValue("lithography"),
		TDP:                utils.ParseInt(getValue("tdp")),
		Features:           getValue("features"),
		MemoryType:         getValue("memory_type"),
		MaxRAM:             getValue("max_ram"),
		MaxRAMFrequency:    getValue("max_ram_frequency"),
		IntegratedGraphics: utils.ParseBool(getValue("integrated_graphics")),
		GraphicsModel:      getValue("graphics_model"),
		MaxTemperature:     utils.ParseInt(getValue("max_temperature")),
		PackageContents:    getValue("package_contents"),
		CountryOfOrigin:    getValue("country_of_origin"),
		// Можно добавить ImageURLs через JSON, если есть
		Images: files,
	}

	err = h.service.UpdateProcessor(procID, procDto)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
}
