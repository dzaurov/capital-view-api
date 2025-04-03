package handlers

import (
	"capital-view-api/db"     // Adjust import path if needed
	"capital-view-api/models" // Adjust import path if needed
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- BalanceSheet Handlers ---

// CreateBalanceSheet godoc
// @Summary Create a new balance sheet entry
// @Description Add a new balance sheet record to the database
// @Tags balance-sheets
// @Accept json
// @Produce json
// @Param balanceSheet body models.BalanceSheet true "Balance Sheet data to create"
// @Success 201 {object} models.BalanceSheet
// @Failure 400 {object} HTTPError "Bad Request - Invalid input data"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /balance-sheets [post]
func CreateBalanceSheet(c *gin.Context) {
	var input models.BalanceSheet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	result := db.DB.Create(&input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GetBalanceSheets godoc
// @Summary Get all balance sheet entries
// @Description Retrieve a list of all balance sheet records
// @Tags balance-sheets
// @Produce json
// @Success 200 {array} models.BalanceSheet
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /balance-sheets [get]
func GetBalanceSheets(c *gin.Context) {
	var balanceSheets []models.BalanceSheet
	result := db.DB.Find(&balanceSheets)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}
	c.JSON(http.StatusOK, balanceSheets)
}

// GetBalanceSheet godoc
// @Summary Get a single balance sheet entry by ID
// @Description Retrieve details of a specific balance sheet record using its ID
// @Tags balance-sheets
// @Produce json
// @Param id path int true "Balance Sheet ID"
// @Success 200 {object} models.BalanceSheet
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Balance Sheet not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /balance-sheets/{id} [get]
func GetBalanceSheet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var balanceSheet models.BalanceSheet
	result := db.DB.First(&balanceSheet, uint(id))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("balance sheet not found")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		}
		return
	}

	c.JSON(http.StatusOK, balanceSheet)
}

// UpdateBalanceSheet godoc
// @Summary Update an existing balance sheet entry
// @Description Modify the details of an existing balance sheet record by ID
// @Tags balance-sheets
// @Accept json
// @Produce json
// @Param id path int true "Balance Sheet ID"
// @Param balanceSheet body models.BalanceSheet true "Balance Sheet data to update"
// @Success 200 {object} models.BalanceSheet
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format or input data"
// @Failure 404 {object} HTTPError "Not Found - Balance Sheet not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /balance-sheets/{id} [put]
func UpdateBalanceSheet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var existingBalanceSheet models.BalanceSheet
	if err := db.DB.First(&existingBalanceSheet, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("balance sheet not found to update")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		}
		return
	}

	var input models.BalanceSheet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	result := db.DB.Model(&existingBalanceSheet).Updates(input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusOK, existingBalanceSheet)
}

// DeleteBalanceSheet godoc
// @Summary Delete a balance sheet entry by ID
// @Description Remove a balance sheet record from the database using its ID
// @Tags balance-sheets
// @Produce json
// @Param id path int true "Balance Sheet ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Balance Sheet not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /balance-sheets/{id} [delete]
func DeleteBalanceSheet(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	result := db.DB.Delete(&models.BalanceSheet{}, uint(id))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, NewHTTPError(errors.New("balance sheet not found to delete")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Balance Sheet deleted successfully"})
}
