// handlers/member_handlers.go
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

// GetMembersByRegcode godoc
// @Summary Получить участников компании по Regcode
// @Description Возвращает пагинированный список участников (members) для указанной компании.
// @Tags member
// @Produce json
// @Param regcode path string true "Regcode компании"
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Записей на странице" default(20) minimum(1) maximum(100)
// @Success 200 {object} models.PaginatedResponse{data=[]models.Member} "Пагинированный список участников"
// @Failure 400 {object} HTTPError "Неверный Regcode"
// @Failure 500 {object} HTTPError "Внутренняя ошибка сервера"
// @Router /members/by-regcode/{regcode} [get] // <-- Пример роута
func GetMembersByRegcode(c *gin.Context) {
	regcode := c.Param("regcode")
	if regcode == "" {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("regcode не может быть пустым")))
		return
	}

	pagination := utils.GetPaginationParams(c)

	var members []models.Member
	var totalRecords int64

	// Базовый запрос с фильтром по regcode
	queryBuilder := db.DB.Model(&models.Member{}).Where("legal_entity_registration_number = ?", regcode)
	// Или использовать at_legal_entity_registration_number? Зависит от вашей логики.
	// queryBuilder := db.DB.Model(&models.Member{}).Where("at_legal_entity_registration_number = ?", regcode)

	// Считаем общее количество
	if err := queryBuilder.Count(&totalRecords).Error; err != nil {
		log.Printf("Error counting members for regcode %s: %v", regcode, err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		return
	}

	// Получаем данные для страницы
	err := queryBuilder.Limit(pagination.Limit).Offset(pagination.Offset).Find(&members).Error
	if err != nil {
		log.Printf("Error finding members for regcode %s with pagination: %v", regcode, err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		return
	}

	// Формируем ответ
	response := models.PaginatedResponse{
		TotalRecords: totalRecords,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
		Data:         members,
	}

	c.JSON(http.StatusOK, response)
}
