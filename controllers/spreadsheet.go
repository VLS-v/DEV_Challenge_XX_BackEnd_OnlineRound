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
	Saves saves.Saves
}

func New(saves saves.Saves) *Controller {
	return &Controller{
		Saves: saves,
	}
}

/*
	POST /api/v1/devchallenge-xx/var1 with {“value:”: “1”}
	Response: {“value:”: “1”, “result”: “1”}
	POST /api/v1/devchallenge-xx/var2 with {“value”: “2”}
	Response: {“value:”: “2”, “result”: “2”}
	POST /api/v1/devchallenge-xx/var3 with {“value”: “=var1+var2”}
	Response: {“value”: “=var1+var2”, “result”: “3”}
*/
// POST: "/api/v1/:sheet_id/:cell_id"
func (c *Controller) SetCellValue(ctx echo.Context) error {
	sheetId := ctx.Param("sheet_id")
	cellId := ctx.Param("cell_id")

	isSheetIdValid := isIdValid(sheetId)
	if isSheetIdValid == false {
		return ctx.JSON(http.StatusBadRequest, models.CellResponse{Value: cellId, Result: "Invalid sheet id!"})
	}

	isCellIdValid := isIdValid(cellId)
	if isCellIdValid == false {
		return ctx.JSON(http.StatusBadRequest, models.CellResponse{Value: cellId, Result: "Invalid cell id!"})
	}
	return ctx.JSON(http.StatusCreated, nil)
	response := &models.Cell{
		Value: cellId,
		//Result: sheetId,
	}
	ctx.String(http.StatusOK, "Hello, World!")
	return ctx.JSON(http.StatusOK, response)

	return nil
}

// GET: "/api/v1/:sheet_id/"
func (c *Controller) GetSheet(ctx echo.Context) error {
	sheetId := ctx.Param("sheet_id")
	_, sheetExists := c.Saves.SavesData[sheetId]
	if sheetExists {
		return ctx.JSON(http.StatusOK, c.Saves.SavesData)
	}
	return ctx.String(http.StatusNotFound, fmt.Sprintf("Sheet %s is missing", sheetId))
}

// GET: "/api/v1/:sheet_id/:cell_id"
func (c *Controller) GetCell(ctx echo.Context) error {
	sheetId := ctx.Param("sheet_id")
	cellId := ctx.Param("cell_id")

	_, sheetExists := c.Saves.SavesData[sheetId]
	if !sheetExists {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("Sheet %s is missing", sheetId))
	}

	_, cellExists := c.Saves.SavesData[sheetId][cellId]
	if cellExists {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("Cell %s is missing", cellId))
	}

	expression := "=_ts2_2t+var2(3*var1)+2/var1+var4"
	expression = "2-var2+3*3+var3+var4"

	expression = strings.TrimPrefix(expression, "=")
	fmt.Println(expression, " TEST")

	s := models.Sheet{
		"var1": {Value: "1", Result: "1"},
		"var2": {Value: "2", Result: "2"},
	}

	t := models.SetCell{
		Value: "1",
	}

	result, err := exprEval.EvaluateExpression(t.Value, s)
	if err != nil {
		fmt.Println("Помилка обчислення виразу:", err)
		return ctx.String(http.StatusOK, err.Error())
	}
	s[cellId] = &models.Cell{Value: t.Value, Result: result}

	fmt.Println("Результат обчислення: ", result)
	err = c.Saves.Write()
	//fmt.Printf("Ooops. Failed to save the value, try again. Input value %v: %s\n", cellId, *s[cellId])
	if err != nil {
		delete(s, cellId)
		return ctx.String(http.StatusOK, fmt.Sprintf("Ooops. Failed to save the value, try again. Input value %v: %s\n", cellId, *s[cellId]))
	}

	response := &models.Cell{
		Value:  t.Value,
		Result: result,
	}
	ctx.String(http.StatusOK, "Hello, World!")
	return ctx.JSON(http.StatusOK, response)
}

func isIdValid(input string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]*$", input)
	return match
}
