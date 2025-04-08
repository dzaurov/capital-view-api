// handlers/beneficial_owner_handlers.go
package handlers

import (
	"errors"
	"log"
	"net/http"

	"capital-view-api/db"
	"capital-view-api/models"
	"capital-view-api/utils"

	"github.com/gin-gonic/gin"
)

// GetBeneficialOwnersByRegcode godoc
// @Summary Получить бенефициаров компании по Regcode
// @Description Возвращает пагинированный список бенефициаров (beneficial owners) для указанной компании.
// @Tags beneficial_owner
// @Produce json
// @Param regcode path string true "Regcode компании"
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Записей на странице" default(20) minimum(1) maximum(100)
// @Success 200 {object} models.PaginatedResponse{data=[]models.BeneficialOwner} "Пагинированный список бенефициаров"
// @Failure 400 {object} HTTPError "Неверный Regcode"
// @Failure 500 {object} HTTPError "Внутренняя ошибка сервера"
// @Router /beneficial-owners/by-regcode/{regcode} [get] // <-- Пример роута
func GetBeneficialOwnersByRegcode(c *gin.Context) {
	regcode := c.Param("regcode")
	if regcode == "" {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("regcode не может быть пустым")))
		return
	}

	pagination := utils.GetPaginationParams(c)

	var owners []models.BeneficialOwner
	var totalRecords int64

	// Базовый запрос
	queryBuilder := db.DB.Model(&models.BeneficialOwner{}).Where("legal_entity_registration_number = ?", regcode)

	// Считаем общее количество
	if err := queryBuilder.Count(&totalRecords).Error; err != nil {
		log.Printf("Error counting beneficial owners for regcode %s: %v", regcode, err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		return
	}

	// Получаем данные для страницы
	err := queryBuilder.Limit(pagination.Limit).Offset(pagination.Offset).Find(&owners).Error
	if err != nil {
		log.Printf("Error finding beneficial owners for regcode %s with pagination: %v", regcode, err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		return
	}

	response := models.PaginatedResponse{
		TotalRecords: totalRecords,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
		Data:         owners,
	}

	c.JSON(http.StatusOK, response)
}
