package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"spreadsheets/utils/saves"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestIsIdValid(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"var1", true},                        // Valid identifier
		{"Var_123", true},                     // Valid identifier
		{"_var", true},                        // Valid identifier
		{"123var", false},                     // Invalid: starts with a number
		{"var@", false},                       // Invalid: contains special character
		{"", false},                           // Invalid: empty string
		{"va r", false},                       // Invalid: contains whitespace
		{"va&*r", false},                      // Invalid: contains special character
		{"vAr_2", true},                       // Valid identifier
		{"Another_Var", true},                 // Valid identifier
		{"_variable", true},                   // Valid identifier
		{"has spaces", false},                 // Invalid: contains whitespace
		{"double__underscore", true},          // Valid identifier
		{"identifier!", false},                // Invalid: contains special character
		{"longer_identifier1234567890", true}, // Valid identifier
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result := isIdValid(test.input)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestSpreadsheetMethods(t *testing.T) {
	var testCases = []struct {
		name          string
		requestString string
		requestMethod string
		input         string
		expectBody    string
		expectCode    int
	}{
		// First, let's do the test examples from the problem condition.

		/*
			POST /api/v1/devchallenge-xx/var1 with {“value:”: “0”}
			Response: {“value:”: “0”, “result”: “0”}
		*/
		{
			name:          "1_SetCell",
			requestString: "/api/v1/dev_challenge/var1",
			requestMethod: http.MethodPost,
			input:         "0",
			expectBody:    "{\"value\":\"0\",\"result\":\"0\"}\n",
			expectCode:    http.StatusCreated,
		},
		/*
			POST /api/v1/devchallenge-xx/var1 with {“value:”: “1”}
			Response: {“value:”: “1”, “result”: “1”}
		*/
		{
			name:          "2_SetCell",
			requestString: "/api/v1/dev_challenge/var1",
			requestMethod: http.MethodPost,
			input:         "1",
			expectBody:    "{\"value\":\"1\",\"result\":\"1\"}\n",
			expectCode:    http.StatusCreated,
		},
		/*
			POST /api/v1/devchallenge-xx/var2 with {“value”: “2”}
			Response: {“value:”: “2”, “result”: “2”}
		*/
		{
			name:          "3_SetCell",
			requestString: "/api/v1/dev_challenge/var2",
			requestMethod: http.MethodPost,
			input:         "2",
			expectBody:    "{\"value\":\"2\",\"result\":\"2\"}\n",
			expectCode:    http.StatusCreated,
		},
		/*
			POST /api/v1/devchallenge-xx/var3 with {“value”: “=var1+var2”}
			Response: {“value”: “=var1+var2”, “result”: “3”}
		*/
		{
			name:          "4_SetCellValue",
			requestString: "/api/v1/dev_challenge/var3",
			requestMethod: http.MethodPost,
			input:         "=var1+var2",
			expectBody:    "{\"value\":\"=var1+var2\",\"result\":\"3\"}\n",
			expectCode:    http.StatusCreated,
		},
		/*
			POST /api/v1/devchallenge-xx/var4 with {“value”: “=var3+var4”}
			Response: {“value”: “=var3+var4”, “result”: “ERROR”}
		*/
		{
			name:          "5_SetCellValue",
			requestString: "/api/v1/dev_challenge/var4",
			requestMethod: http.MethodPost,
			input:         "=var3+var4",
			expectBody:    "{\"value\":\"=var3+var4\",\"result\":\"ERROR\"}\n",
			expectCode:    http.StatusUnprocessableEntity,
		},
		/*
			GET /api/v1/devchallenge-xx/var1
			Response: {“value”: “1”, result: “1”}
		*/
		{
			name:          "1_GetCell",
			requestString: "/api/v1/dev_challenge/var1",
			requestMethod: http.MethodGet,
			input:         "",
			expectBody:    "{\"value\":\"1\",\"result\":\"1\"}\n",
			expectCode:    http.StatusOK,
		},
		/*
			I missed the query:

				GET /api/v1/devchallenge-xx/var1
				Response: {“value”: “2”, result: “2”}

			Because it repeats with the previous one and the wrong result is expected. The value of var1 cannot be 1 and 2 at the same time.
		*/
		/*
			GET /api/v1/devchallenge-xx/var3
			Response: {“value”: “=var1+var2”, result: “3”}
		*/
		{
			name:          "2_GetCell",
			requestString: "/api/v1/dev_challenge/var3",
			requestMethod: http.MethodGet,
			input:         "",
			expectBody:    "{\"value\":\"=var1+var2\",\"result\":\"3\"}\n",
			expectCode:    http.StatusOK,
		},
		// The test for obtaining a non-existent cell.
		{
			name:          "3_GetCell",
			requestString: "/api/v1/dev_challenge/var_test",
			requestMethod: http.MethodGet,
			input:         "",
			expectBody:    "Cell var_test is missing",
			expectCode:    http.StatusNotFound,
		},
		// GET /api/v1/:sheet_id
		{
			name:          "1_GetSheet",
			requestString: "/api/v1/dev_challenge/",
			requestMethod: http.MethodGet,
			input:         "",
			expectBody:    "{\"var1\":{\"value\":\"1\",\"result\":\"1\"},\"var2\":{\"value\":\"2\",\"result\":\"2\"},\"var3\":{\"value\":\"=var1+var2\",\"result\":\"3\"}}\n",
			expectCode:    http.StatusOK,
		},
		// The test for obtaining a non-existent sheet.
		{
			name:          "2_GetSheet",
			requestString: "/api/v1/dev_challenge_test/",
			requestMethod: http.MethodGet,
			input:         "",
			expectBody:    "Sheet dev_challenge_test is missing",
			expectCode:    http.StatusNotFound,
		},
		// Next, I do my own tests.
		{
			name:          "1_OwnSetCell",
			requestString: "/api/v1/sheet_1/cell1",
			requestMethod: http.MethodPost,
			input:         "1+2",
			expectBody:    "{\"value\":\"1+2\",\"result\":\"3\"}\n",
			expectCode:    http.StatusCreated,
		},
		// Check if an error is received when adding a cell with a formula in which an undefined variable is used.
		{
			name:          "2_OwnSetCell",
			requestString: "/api/v1/sheet_1/cell2",
			requestMethod: http.MethodPost,
			input:         "=(2+2)*cell2",
			expectBody:    "{\"value\":\"=(2+2)*cell2\",\"result\":\"ERROR\"}\n",
			expectCode:    http.StatusUnprocessableEntity,
		},
		{
			name:          "3_OwnSetCell",
			requestString: "/api/v1/sheet_1/cell3",
			requestMethod: http.MethodPost,
			input:         "(3+2)**cell1",
			expectBody:    "{\"value\":\"(3+2)**cell1\",\"result\":\"125\"}\n",
			expectCode:    http.StatusCreated,
		},
		{
			name:          "4_OwnSetCell",
			requestString: "/api/v1/sheet_1/cell1",
			requestMethod: http.MethodPost,
			input:         "3*1",
			expectBody:    "{\"value\":\"3*1\",\"result\":\"3\"}\n",
			expectCode:    http.StatusCreated,
		},
		{
			name:          "5_OwnSetCell",
			requestString: "/api/v1/sheet_1/cell4",
			requestMethod: http.MethodPost,
			input:         "3*1.5",
			expectBody:    "{\"value\":\"3*1.5\",\"result\":\"4.5\"}\n",
			expectCode:    http.StatusCreated,
		},
		// Check if the value in cell3 is recalculated after changing the formula in cell1.
		{
			name:          "1_OwnGetSheet",
			requestString: "/api/v1/sheet_1/",
			requestMethod: http.MethodGet,
			input:         "",
			expectBody:    "{\"cell1\":{\"value\":\"3*1\",\"result\":\"3\"},\"cell3\":{\"value\":\"(3+2)**cell1\",\"result\":\"125\"},\"cell4\":{\"value\":\"3*1.5\",\"result\":\"4.5\"}}\n",
			expectCode:    http.StatusOK,
		},
		{
			name:          "1_OwnGetCell",
			requestString: "/api/v1/sheet_1/cell1",
			requestMethod: http.MethodGet,
			input:         "",
			expectBody:    "{\"value\":\"3*1\",\"result\":\"3\"}\n",
			expectCode:    http.StatusOK,
		},
		/*
			Checking the looping of the formula.

			Cell3 contains the formula "(3+2)*cell1". Let's try to write the formula "=cell3" to cell1.
			This causes an error because cell1 would depend on itself (cell1:"="(3+2)*cell1).
		*/
		{
			name:          "1_OwnSetCell",
			requestString: "/api/v1/sheet_1/cell1",
			requestMethod: http.MethodPost,
			input:         "=cell3",
			expectBody:    "{\"value\":\"=cell3\",\"result\":\"ERROR\"}\n",
			expectCode:    http.StatusUnprocessableEntity,
		},
	}

	savesInstance := saves.Saves{}

	err := savesInstance.Open("./")
	if err != nil {
		return
	}

	err = savesInstance.Load()
	if err != nil {
		return
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			c := New(&savesInstance)
			e.POST("/api/v1/:sheet_id/:cell_id", c.SetCellValue)
			e.GET("/api/v1/:sheet_id/", c.GetSheet)
			e.GET("/api/v1/:sheet_id/:cell_id", c.GetCell)

			var req *http.Request
			var requestBody *bytes.Reader
			if tc.requestMethod == http.MethodPost {
				jsonData, err := json.Marshal(struct {
					Value string `json:"value"`
				}{
					Value: tc.input,
				})

				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				requestBody = bytes.NewReader(jsonData)
				req = httptest.NewRequest(tc.requestMethod, tc.requestString, requestBody)
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			} else {
				req = httptest.NewRequest(tc.requestMethod, tc.requestString, nil)
			}

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectBody, rec.Body.String())
			assert.Equal(t, tc.expectCode, rec.Code)

		})
	}
	savesInstance.SavesFile.Close()
	os.Remove("./" + savesInstance.SavesFile.Name())
}
