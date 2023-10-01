package expressionevaluator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"spreadsheets/models"
	"strings"

	"github.com/Knetic/govaluate"
)

func EvaluateExpression(expression string, savesData models.Sheet) (string, error) {
	/* variables, err := extractVariables(expression)
	if err != nil {
		return expression, err
	} */
	preparedExpression, _ := replaceVariablesInExpression(expression, &savesData)
	/* var preparedExpression string
	var variables []string
	var err error
	variables, err = extractVariables(expression)
	if err != nil {
		return expression, err
	}
	for len(variables) > 0 {
		preparedExpression, err = replaceVariablesInExpression(expression, variables, &savesData)
		variables, err = extractVariables(expression)
		if err != nil {
			return expression, err
		}
	} */
	/* fmt.Println("Variables from expression:")
	for _, variable := range variables {
		fmt.Println(variable)
		//fmt.Println(saves["sheet1"][variable])
	} */
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

func replaceVariablesInExpression(expression string /* variables []string, */, savesData *models.Sheet) (string, error) {
	// Прибираємо "=" з виразу
	//expression = strings.TrimPrefix(expression, "=")

	variables, err := extractVariables(expression)
	if err != nil {
		return expression, err
	}
	//2-var2+3*3+var3+var4
	// var2->1, var3->var1+var2->2+1, var4->var1+var2->2+1
	for _, variable := range variables {
		// Перевіряємо, чи значення є формулою
		if strings.HasPrefix((*savesData)[variable].Value, "=") {
			value := strings.TrimPrefix((*savesData)[variable].Value, "=")
			expression = strings.ReplaceAll(expression, variable, value)
			/* vars, err := extractVariables(value)
			if err != nil {
				return expression, err
			} */

			expression, err = replaceVariablesInExpression(expression, savesData)
			if err != nil {
				return expression, err
			}

			variables, err = extractVariables(expression)
			if err != nil {
				return expression, err
			}

			if len(variables) == 0 {
				return expression, nil
			}

			/* expression = strings.ReplaceAll(expression, variable, extractedExpression)
			return expression, nil */
			/* extractedExpression, err := replaceVariablesInExpression(expression, savesData)
			if err != nil {
				return expression, err
			}
			expression = strings.ReplaceAll(expression, variable, extractedExpression)
			return expression, nil */
			/* expression, err = replaceVariablesInExpression(expression, vars, &savesData)
			if err != nil {
				return expression, err
			} */
		} else {
			// Замінюємо назви змінних їх значеннями
			expression = strings.ReplaceAll(expression, variable, (*savesData)[variable].Value)
		}
	}
	return expression, nil
}

type visitor struct {
	variables []string
}

// Visit викликається для кожного вузла в AST
func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if ident, ok := node.(*ast.Ident); ok {
		// Додамо ім'я змінної до списку
		v.variables = append(v.variables, ident.Name)
	}
	return v
}

func extractVariables(expression string) ([]string, error) {
	// Створимо токенізатор та парсер
	//fs := token.NewFileSet()
	expression = strings.TrimPrefix(expression, "=")
	expr, err := parser.ParseExpr(expression)
	if err != nil {
		fmt.Println("Помилка парсингу виразу:", err)
		return nil, err
	}

	// Створимо візитора для обходу AST
	v := &visitor{}
	ast.Walk(v, expr)

	return v.variables, nil
}
