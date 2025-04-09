// handlers/company_handlers.go
package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"capital-view-api/db"
	"capital-view-api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCompanyDetailsByRegcode godoc
// @Summary Получить полную информацию о компании по Regcode
// @Description Получает детальную информацию о компании, включая участников, бенефициаров и фин. отчеты, по её точному Regcode.
// @Tags company
// @Produce json
// @Param regcode path string true "Regcode компании"
// @Success 200 {object} models.Registers "Полная информация о компании (с вложенными данными)"
// @Failure 400 {object} HTTPError "Неверный или отсутствующий Regcode"
// @Failure 404 {object} HTTPError "Компания не найдена"
// @Failure 500 {object} HTTPError "Внутренняя ошибка сервера"
// @Router /company/{regcode} [get] // <-- Новый роут
func GetCompanyDetailsByRegcode(c *gin.Context) {
	regcode := c.Param("regcode")
	if regcode == "" {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("regcode не может быть пустым")))
		return
	}
	log.Printf("GetCompanyDetailsByRegcode: Fetching details for regcode: %s", regcode)

	var company models.Registers // Используем основную модель Registers
	err := db.DB.
		// Предзагружаем все необходимые связанные данные
		Preload("Members").
		Preload("BeneficialOwners").
		Preload("FinancialStatements", func(db *gorm.DB) *gorm.DB {
			return db.Order("financial_statements.year DESC") // Сортируем отчеты
		}).
		Preload("FinancialStatements.IncomeStatement").
		Preload("FinancialStatements.BalanceSheet").
		Preload("FinancialStatements.CashFlowStatement").
		Where("regcode = ?", regcode). // Точный поиск по regcode
		First(&company).Error          // Ищем одну запись

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("GetCompanyDetailsByRegcode: Company not found for regcode %s", regcode)
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("компания с таким regcode не найдена")))
		} else {
			log.Printf("GetCompanyDetailsByRegcode: Error fetching company details for %s: %v", regcode, err)
			c.JSON(http.StatusInternalServerError, NewHTTPError(fmt.Errorf("ошибка получения данных компании: %w", err)))
		}
		return
	}

	log.Printf("GetCompanyDetailsByRegcode: Successfully fetched details for regcode %s", regcode)
	// Возвращаем найденный объект Registers со всеми предзагруженными данными
	c.JSON(http.StatusOK, company)
}
