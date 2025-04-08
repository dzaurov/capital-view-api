// handlers/financial_statement_handlers.go
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

// GetFinancialStatementsByRegcode godoc
// @Summary Получить фин. отчеты компании по Regcode
// @Description Возвращает пагинированный список фин. отчетов (financial statements) для указанной компании.
// @Tags financial_statement
// @Produce json
// @Param regcode path string true "Regcode компании"
// @Param page query int false "Номер страницы" default(1) minimum(1)
// @Param limit query int false "Записей на странице" default(20) minimum(1) maximum(100)
// @Success 200 {object} models.PaginatedResponse{data=[]models.FinancialStatement} "Пагинированный список фин. отчетов"
// @Failure 400 {object} HTTPError "Неверный Regcode"
// @Failure 500 {object} HTTPError "Внутренняя ошибка сервера"
// @Router /financial-statements/by-regcode/{regcode} [get] // <-- Пример роута
func GetFinancialStatementsByRegcode(c *gin.Context) {
	regcode := c.Param("regcode")
	if regcode == "" {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("regcode не может быть пустым")))
		return
	}

	pagination := utils.GetPaginationParams(c)

	var statements []models.FinancialStatement
	var totalRecords int64

	// Базовый запрос
	queryBuilder := db.DB.Model(&models.FinancialStatement{}).Where("legal_entity_registration_number = ?", regcode)

	// Считаем общее количество
	if err := queryBuilder.Count(&totalRecords).Error; err != nil {
		log.Printf("Error counting financial statements for regcode %s: %v", regcode, err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		return
	}

	// Получаем данные для страницы + сортировка по убыванию года
	err := queryBuilder.Order("year desc").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&statements).Error
	if err != nil {
		log.Printf("Error finding financial statements for regcode %s with pagination: %v", regcode, err)
		c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		return
	}

	response := models.PaginatedResponse{
		TotalRecords: totalRecords,
		Page:         pagination.Page,
		Limit:        pagination.Limit,
		Data:         statements,
	}

	c.JSON(http.StatusOK, response)
}
