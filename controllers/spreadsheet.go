package controllers

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	exprEval "spreadsheets/helpers/expressionevaluator"
	"spreadsheets/models"
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
	file  *os.File
	saves models.SavesData
}

func New(file *os.File, saves models.SavesData) *Controller {
	return &Controller{
		file:  file,
		saves: saves,
	}
}

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
	_, keyExists := c.saves[sheetId]
	if keyExists {
		return ctx.JSON(http.StatusOK, c.saves)
	}
	return ctx.String(http.StatusNotFound, fmt.Sprintf("Sheet %s is missing", sheetId))
}

// GET: "/api/v1/:sheet_id/:cell_id"
func (c *Controller) GetCell(ctx echo.Context) error {
	//sheetId := ctx.Param("sheet_id")
	cellId := ctx.Param("cell_id")

	/*
		POST /api/v1/devchallenge-xx/var1 with {“value:”: “1”}
		Response: {“value:”: “1”, “result”: “1”}
		POST /api/v1/devchallenge-xx/var2 with {“value”: “2”}
		Response: {“value:”: “2”, “result”: “2”}
		POST /api/v1/devchallenge-xx/var3 with {“value”: “=var1+var2”}
		Response: {“value”: “=var1+var2”, “result”: “3”}
	*/
	/* var1 := "1"
	var2 := "2"
	var3 := "=var1+var2" */
	expression := "=_ts2_2t+var2(3*var1)+2/var1+var4"
	expression = "2-var2+3*3+var3+var4"

	expression = strings.TrimPrefix(expression, "=")
	fmt.Println(expression, " TEST")

	/* variables, err := extractVariables(expression)
	if err != nil {
		return ctx.String(http.StatusNotFound, err.Error())
	}

	fmt.Println("Variables from expression:")
	for _, variable := range variables {
		fmt.Println(variable)
		fmt.Println(saves["sheet1"][variable])
		ctx.String(http.StatusOK, variable+"\n")
	} */

	t := models.Sheet{
		"var1": {Value: "1", Result: "1"},
		"var2": {Value: "2", Result: "2"},
		"var3": {Value: "=var1+var2"},
		"var4": {Value: "=var1+var2"},
		"var5": {Value: "value5"},
	}

	result, err := exprEval.EvaluateExpression(expression, t)
	if err != nil {
		fmt.Println("Помилка обчислення виразу:", err)
		return ctx.String(http.StatusOK, err.Error())
	}

	fmt.Println("Результат обчислення: ", result)

	response := &models.Cell{
		Value: cellId,
		//Result: sheetId,
	}
	ctx.String(http.StatusOK, "Hello, World!")
	return ctx.JSON(http.StatusOK, response)
}

/* func evaluateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return expression, err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return expression, err
	}
	return fmt.Sprintf("%v", result), nil
}

func replaceVariablesInExpression(expression string, variables []string, savesData map[string]string) (string, error) {
	// Прибираємо "=" з виразу
	expression = strings.TrimPrefix(expression, "=")

	for _, variable := range variables {
		// Перевіряємо, чи значення є формулою
		if strings.HasPrefix(savesData[variable], "=") {
			value := strings.TrimPrefix(savesData[variable], "=")
			expression = strings.ReplaceAll(expression, variable, value)
			vars, err := extractVariables(value)
			if err != nil {
				return expression, err
			}
			expression, err = replaceVariablesInExpression(expression, vars, savesData)
			if err != nil {
				return expression, err
			}
		} else {
			// Замінюємо назви змінних їх значеннями
			expression = strings.ReplaceAll(expression, variable, savesData[variable])
		}
	}
	return expression, nil
}

func evaluateExpression(expression string, variables map[string]float64) (float64, error) {
	// Прибираємо "=" з виразу
	expression = strings.TrimPrefix(expression, "=")

	// Копіюємо значення змінних для безпечного обчислення
	calculatedVariables := make(map[string]float64)
	for key, value := range variables {
		calculatedVariables[key] = value
	}

	// Ітеративно замінюємо назви змінних їх значеннями, доки можна замінювати
	for {
		replaced := false
		for variable, value := range calculatedVariables {
			// Замінюємо назви змінних їх значеннями
			expression = strings.ReplaceAll(expression, variable, fmt.Sprintf("%f", value))
		}

		// Обчислюємо вираз
		result, err := math.Eval(expression)
		if err != nil {
			return 0, err
		}

		// Якщо результат вже є числовим значенням, повертаємо його
		if !math.IsNaN(result) {
			return result, nil
		}

		// Якщо не відбулось заміни, завершуємо цикл
		if !replaced {
			break
		}
	}

	return 0, fmt.Errorf("неможливо обчислити вираз")
} */

/* func replaceVariables(node ast.Node, variables map[string]string) string {
	switch n := node.(type) {
	case *ast.Ident:
		if value, ok := variables[n.Name]; ok {
			return value
		}
		return n.Name
	case *ast.BinaryExpr:
		return replaceVariables(n.X, variables) + n.Op.String() + replaceVariables(n.Y, variables)
	default:
		return ""
	}
} */

/* type visitor struct {
	variables []string
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if ident, ok := node.(*ast.Ident); ok {
		v.variables = append(v.variables, ident.Name)
	}
	return v
}

func extractVariables(expression string) ([]string, error) {
	expr, err := parser.ParseExpr(expression)
	if err != nil {
		fmt.Println("Помилка парсингу виразу:", err)
		return nil, err
	}

	v := &visitor{}
	ast.Walk(v, expr)

	return v.variables, nil
} */

func isIdValid(input string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z_][a-zA-Z0-9_]*$", input)
	return match
}
