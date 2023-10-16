package expressionevaluator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"spreadsheets/models"
	"strings"

	"github.com/Knetic/govaluate"
)

func EvaluateExpression(expression string, cellId string, savesData models.Sheet) (string, error) {
	expression = strings.TrimPrefix(expression, "=")
	preparedExpression, err := replaceVariablesInExpression(expression, cellId, &savesData)
	if err != nil {
		return expression, err
	}
	expr, err := govaluate.NewEvaluableExpression(preparedExpression)
	if err != nil {
		return expression, err
	}
	result, err := expr.Evaluate(nil)
	if err != nil {
		return expression, err
	}

	return fmt.Sprintf("%v", result), nil
}

func replaceVariablesInExpression(expression string, cellId string, savesData *models.Sheet) (string, error) {
	variables, err := extractVariables(expression)
	if err != nil {
		return expression, err
	}

	errorLooping := errors.New("The formula calculation is looping!")

	for _, variable := range variables {
		_, variableExists := (*savesData)[variable]
		if !variableExists {
			return expression, errors.New(fmt.Sprintf("Variable %s is not exist!", variable))
		}
		variablesInSubFormula, _ := extractVariables((*savesData)[variable].Value)
		for _, subVariable := range variablesInSubFormula {
			if subVariable == cellId {
				return expression, errorLooping
			}
		}
		if variable == cellId {
			return expression, errorLooping
		}
		expression = strings.ReplaceAll(expression, variable, (*savesData)[variable].Result)
	}

	return expression, nil
}

type visitor struct {
	variables []string
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if ident, ok := node.(*ast.Ident); ok {
		v.variables = append(v.variables, ident.Name)
	}
	return v
}

func extractVariables(expression string) ([]string, error) {
	expression = strings.TrimPrefix(expression, "=")
	expr, err := parser.ParseExpr(expression)
	if err != nil {
		return nil, err
	}

	v := &visitor{}
	ast.Walk(v, expr)

	return v.variables, nil
}

func RecursionUpdate(newCellValue *models.Cell, cellId string, sheetData models.Sheet) (models.Sheet, error) {
	if sheetData == nil {
		sheetData = models.Sheet{}
	}
	sheetData[cellId] = newCellValue

	for cellId, cell := range sheetData {
		res, err := EvaluateExpression(cell.Value, cellId, sheetData)
		if err != nil {
			return nil, err
		}
		cell.Result = res
	}

	return sheetData, nil
}
