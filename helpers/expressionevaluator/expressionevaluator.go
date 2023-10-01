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

func EvaluateExpression(expression string, savesData models.Sheet) (string, error) {
	preparedExpression, _ := replaceVariablesInExpression(expression, &savesData)
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

func replaceVariablesInExpression(expression string, savesData *models.Sheet) (string, error) {
	variables, err := extractVariables(expression)
	if err != nil {
		return expression, err
	}

	for _, variable := range variables {
		_, variableExists := (*savesData)[variable]
		if !variableExists {
			return expression, errors.New(fmt.Sprintf("Variable %s is not exist!", variable))
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
		fmt.Println("Помилка парсингу виразу:", err)
		return nil, err
	}

	v := &visitor{}
	ast.Walk(v, expr)

	return v.variables, nil
}
