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

// --- IncomeStatement Handlers ---

// CreateIncomeStatement godoc
// @Summary Create a new income statement entry
// @Description Add a new income statement record to the database
// @Tags income-statements
// @Accept json
// @Produce json
// @Param incomeStatement body models.IncomeStatement true "Income Statement data to create"
// @Success 201 {object} models.IncomeStatement
// @Failure 400 {object} HTTPError "Bad Request - Invalid input data"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /income-statements [post]
func CreateIncomeStatement(c *gin.Context) {
	var input models.IncomeStatement
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

// GetIncomeStatements godoc
// @Summary Get all income statement entries
// @Description Retrieve a list of all income statement records
// @Tags income-statements
// @Produce json
// @Success 200 {array} models.IncomeStatement
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /income-statements [get]
func GetIncomeStatements(c *gin.Context) {
	var incomeStatements []models.IncomeStatement
	result := db.DB.Find(&incomeStatements)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}
	c.JSON(http.StatusOK, incomeStatements)
}

// GetIncomeStatement godoc
// @Summary Get a single income statement entry by ID
// @Description Retrieve details of a specific income statement record using its ID
// @Tags income-statements
// @Produce json
// @Param id path int true "Income Statement ID"
// @Success 200 {object} models.IncomeStatement
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Income Statement not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /income-statements/{id} [get]
func GetIncomeStatement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var incomeStatement models.IncomeStatement
	result := db.DB.First(&incomeStatement, uint(id))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("income statement not found")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		}
		return
	}

	c.JSON(http.StatusOK, incomeStatement)
}

// UpdateIncomeStatement godoc
// @Summary Update an existing income statement entry
// @Description Modify the details of an existing income statement record by ID
// @Tags income-statements
// @Accept json
// @Produce json
// @Param id path int true "Income Statement ID"
// @Param incomeStatement body models.IncomeStatement true "Income Statement data to update"
// @Success 200 {object} models.IncomeStatement
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format or input data"
// @Failure 404 {object} HTTPError "Not Found - Income Statement not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /income-statements/{id} [put]
func UpdateIncomeStatement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var existingIncomeStatement models.IncomeStatement
	if err := db.DB.First(&existingIncomeStatement, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("income statement not found to update")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		}
		return
	}

	var input models.IncomeStatement
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	result := db.DB.Model(&existingIncomeStatement).Updates(input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusOK, existingIncomeStatement)
}

// DeleteIncomeStatement godoc
// @Summary Delete an income statement entry by ID
// @Description Remove an income statement record from the database using its ID
// @Tags income-statements
// @Produce json
// @Param id path int true "Income Statement ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Income Statement not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /income-statements/{id} [delete]
func DeleteIncomeStatement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	result := db.DB.Delete(&models.IncomeStatement{}, uint(id))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, NewHTTPError(errors.New("income statement not found to delete")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Income Statement deleted successfully"})
}
