// handlers/register_handlers.go
package handlers

import (
	"errors"
	"log"
	"net/http"

	"capital-view-api/db"     // <-- Убедитесь, что путь правильный
	"capital-view-api/models" // <-- Убедитесь, что путь правильный
	"capital-view-api/utils"  // <-- Импорт пагинации

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- Если HTTPError и NewHTTPError еще не определены глобально ---
type HTTPError struct {
	Error string `json:"error"`
}

func NewHTTPError(err error) HTTPError {
	return HTTPError{Error: err.Error()}
}

// ---------------------------------------------------------------

// GetRegisterByID godoc
// @Summary Получить информацию о компании по Regcode
// @Description Получает детальную информацию о компании по её Regcode.
// @Tags register
// @Produce json
// @Param regcode path string true "Regcode компании"
// @Success 200 {object} models.Registers "Информация о компании"
// @Failure 400 {object} HTTPError "Неверный Regcode"
// @Failure 404 {object} HTTPError "Компания не найдена"
// @Failure 500 {object} HTTPError "Внутренняя ошибка сервера"
// @Router /register/{regcode} [get]
func GetRegisterByID(c *gin.Context) {
	regcode := c.Param("regcode")
	if regcode == "" {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("regcode не может быть пустым")))
		return
	}

	var register models.Registers // <-- Используем модель Registers
	err := db.DB.Where("regcode = ?", regcode).First(&register).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("компания с таким regcode не найдена")))
		} else {
			log.Printf("Error finding register by regcode %s: %v", regcode, err)
			c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		}
		return
	}

	c.JSON(http.StatusOK, register)
}

// GetAllRegisters godoc
// @Summary Получить список всех записей регистра
// @Description Возвращает пагинированный список записей из таблицы registers.
// @Tags register
// @Produce json
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Записей на странице" default(20) minimum(1) maximum(100)
// @Success 200 {object} models.PaginatedResponse{data=[]models.Registers} "Пагинированный список записей"
// @Failure 500 {object} HTTPError "Внутренняя ошибка сервера"
// @Router /registers [get] // <-- Пример роута, измените на ваш
func GetAllRegisters(c *gin.Context) {
	pagination := utils.GetPaginationParams(c) // <-- Получаем параметры пагинации

	var registers []models.Registers
	var totalRecords int64

	// Базовый запрос
	queryBuilder := db.DB.Model(&models.Registers{})

	// Фильтрация (если нужна, добавьте .Where() сюда)
	// queryBuilder = queryBuilder.Where("some_condition = ?", some_value)

	// Считаем общее количество
	if err := queryBuilder.Count(&totalRecords).Error; err != nil {
		log.Printf("Error counting registers: %v", err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		return
	}

	// Получаем данные для страницы + сортировка (пример)
	err := queryBuilder.Order("name asc"). // <-- Пример сортировки
						Limit(pagination.Limit).
						Offset(pagination.Offset).
						Find(&registers).Error
	if err != nil {
		log.Printf("Error finding registers with pagination: %v", err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		return
	}

	// Формируем ответ
	response := models.PaginatedResponse{
		TotalRecords: totalRecords,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
		Data:         registers,
	}

	c.JSON(http.StatusOK, response)
}
