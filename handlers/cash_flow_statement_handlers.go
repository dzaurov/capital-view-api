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

// --- CashFlowStatement Handlers ---

// CreateCashFlowStatement godoc
// @Summary Create a new cash flow statement entry
// @Description Add a new cash flow statement record to the database
// @Tags cash-flow-statements
// @Accept json
// @Produce json
// @Param cashFlowStatement body models.CashFlowStatement true "Cash Flow Statement data to create"
// @Success 201 {object} models.CashFlowStatement
// @Failure 400 {object} HTTPError "Bad Request - Invalid input data"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /cash-flow-statements [post]
func CreateCashFlowStatement(c *gin.Context) {
	var input models.CashFlowStatement
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

// GetCashFlowStatements godoc
// @Summary Get all cash flow statement entries
// @Description Retrieve a list of all cash flow statement records
// @Tags cash-flow-statements
// @Produce json
// @Success 200 {array} models.CashFlowStatement
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /cash-flow-statements [get]
func GetCashFlowStatements(c *gin.Context) {
	var cashFlowStatements []models.CashFlowStatement
	result := db.DB.Find(&cashFlowStatements)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}
	c.JSON(http.StatusOK, cashFlowStatements)
}

// GetCashFlowStatement godoc
// @Summary Get a single cash flow statement entry by ID
// @Description Retrieve details of a specific cash flow statement record using its ID
// @Tags cash-flow-statements
// @Produce json
// @Param id path int true "Cash Flow Statement ID"
// @Success 200 {object} models.CashFlowStatement
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Cash Flow Statement not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /cash-flow-statements/{id} [get]
func GetCashFlowStatement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var cashFlowStatement models.CashFlowStatement
	result := db.DB.First(&cashFlowStatement, uint(id))

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("cash flow statement not found")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		}
		return
	}

	c.JSON(http.StatusOK, cashFlowStatement)
}

// UpdateCashFlowStatement godoc
// @Summary Update an existing cash flow statement entry
// @Description Modify the details of an existing cash flow statement record by ID
// @Tags cash-flow-statements
// @Accept json
// @Produce json
// @Param id path int true "Cash Flow Statement ID"
// @Param cashFlowStatement body models.CashFlowStatement true "Cash Flow Statement data to update"
// @Success 200 {object} models.CashFlowStatement
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format or input data"
// @Failure 404 {object} HTTPError "Not Found - Cash Flow Statement not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /cash-flow-statements/{id} [put]
func UpdateCashFlowStatement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	var existingCashFlowStatement models.CashFlowStatement
	if err := db.DB.First(&existingCashFlowStatement, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, NewHTTPError(errors.New("cash flow statement not found to update")))
		} else {
			c.JSON(http.StatusInternalServerError, NewHTTPError(err))
		}
		return
	}

	var input models.CashFlowStatement
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(err))
		return
	}

	result := db.DB.Model(&existingCashFlowStatement).Updates(input)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	c.JSON(http.StatusOK, existingCashFlowStatement)
}

// DeleteCashFlowStatement godoc
// @Summary Delete a cash flow statement entry by ID
// @Description Remove a cash flow statement record from the database using its ID
// @Tags cash-flow-statements
// @Produce json
// @Param id path int true "Cash Flow Statement ID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} HTTPError "Bad Request - Invalid ID format"
// @Failure 404 {object} HTTPError "Not Found - Cash Flow Statement not found"
// @Failure 500 {object} HTTPError "Internal Server Error"
// @Router /cash-flow-statements/{id} [delete]
func DeleteCashFlowStatement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewHTTPError(errors.New("invalid ID format")))
		return
	}

	result := db.DB.Delete(&models.CashFlowStatement{}, uint(id))

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, NewHTTPError(result.Error))
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, NewHTTPError(errors.New("cash flow statement not found to delete")))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cash Flow Statement deleted successfully"})
}
