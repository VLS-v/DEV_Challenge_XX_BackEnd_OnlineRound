package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	exprEval "spreadsheets/helpers/expressionevaluator"
	"spreadsheets/models"
	"spreadsheets/utils/saves"
	"strings"

	"github.com/labstack/echo/v4"
)

/*
	Some HTTP statuses
	200 - http.StatusOK
	201 - http.StatusCreated
	400 - http.StatusBadRequest
	404 - http.StatusNotFound
	422 - http.StatusUnprocessableEntity
*/

type Controller struct {
	SavesInstance *saves.Saves
}

func New(savesInstance *saves.Saves) *Controller {
	return &Controller{
		SavesInstance: savesInstance,
	}
}

// POST: "/api/v1/:sheet_id/:cell_id"
func (c *Controller) SetCellValue(ctx echo.Context) error {
	sheetId := strings.ToLower(ctx.Param("sheet_id"))
	cellId := strings.ToLower(ctx.Param("cell_id"))

	requestBody := new(models.SetCell)
	if err := ctx.Bind(requestBody); err != nil {
		return err
	}
	isSheetIdValid := isIdValid(sheetId)
	if isSheetIdValid == false {
		return ctx.JSON(http.StatusBadRequest, models.Cell{Value: requestBody.Value, Result: "Invalid sheet id!"})
	}

	isCellIdValid := isIdValid(cellId)
	if isCellIdValid == false {
		return ctx.JSON(http.StatusBadRequest, models.Cell{Value: requestBody.Value, Result: "Invalid cell id!"})
	}

	requestBody.Value = strings.ToLower(requestBody.Value)
	savesData := c.SavesInstance.SavesData

	result, err := exprEval.EvaluateExpression(requestBody.Value, cellId, savesData[sheetId])
	if err != nil {
		response := models.Cell{
			Value:  requestBody.Value,
			Result: "ERROR",
		}
		return ctx.JSON(http.StatusUnprocessableEntity, response)
	}

	newCellValue := models.Cell{
		Value:  requestBody.Value,
		Result: result,
	}
	savesData[sheetId], err = exprEval.RecursionUpdate(&newCellValue, cellId, savesData[sheetId])
	if err != nil {
		response := models.Cell{
			Value:  requestBody.Value,
			Result: "ERROR",
		}
		return ctx.JSON(http.StatusUnprocessableEntity, response)
	}

	err = c.SavesInstance.Write(savesData)
	if err != nil {
		textErr := fmt.Sprintf("Ooops. Failed to save the value, try again. Input value %v: %s\n", cellId, *savesData[sheetId][cellId])
		return ctx.String(http.StatusBadRequest, textErr)
	}

	response := models.Cell{
		Value:  requestBody.Value,
		Result: result,
	}

	return ctx.JSON(http.StatusCreated, response)
}

// GET: "/api/v1/:sheet_id/"
func (c *Controller) GetSheet(ctx echo.Context) error {
	sheetId := strings.ToLower(ctx.Param("sheet_id"))

	_, sheetExists := c.SavesInstance.SavesData[sheetId]
	if sheetExists {
		return ctx.JSON(http.StatusOK, c.SavesInstance.SavesData[sheetId])
	}
	return ctx.String(http.StatusNotFound, fmt.Sprintf("Sheet %s is missing", sheetId))
}

// GET: "/api/v1/:sheet_id/:cell_id"
func (c *Controller) GetCell(ctx echo.Context) error {
	sheetId := strings.ToLower(ctx.Param("sheet_id"))
	cellId := strings.ToLower(ctx.Param("cell_id"))

	_, sheetExists := c.SavesInstance.SavesData[sheetId]
	if !sheetExists {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("Sheet %s is missing", sheetId))
	}

	_, cellExists := c.SavesInstance.SavesData[sheetId][cellId]
	if !cellExists {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("Cell %s is missing", cellId))
	}

	return ctx.JSON(http.StatusOK, c.SavesInstance.SavesData[sheetId][cellId])
}

func isIdValid(input string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]*$", input)
	return match
}
