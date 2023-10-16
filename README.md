**Table of Contents**

- [How to start service and tests](#how-to-start-service-and-tests)
- [Naming Conventions for Sheet and Cell IDs](#naming-conventions-for-sheet-and-cell-ids)
- [Input Rules for Formulas](#input-rules-for-formulas)
  - [What operators and types does this support?](#what-operators-and-types-does-this-support)

## How to start service and tests
1. Running tests:
```
    docker compose up spreadsheets-test
```
2. Launching the application:
```
    docker compose up spreadsheets
```
I recommend using [Postman](https://www.postman.com/) to make requests to the API.

## Naming Conventions for Sheet and Cell IDs
<span style="color:red">⭑</span> The names are case-independent.

* It should start with a letter (a-z, A-Z) or an underscore (_);
* Subsequent characters can be letters (a-z, A-Z), digits (0-9), or underscores (_).
  
If you enter an invalid sheet or cell name, you will receive an error message:
```
POST: http://localhost:8080/api/v1/sheet1/1cell1
Body: {"value":"1+2"}

RESPONSE: 
{
    "value": "1+2",
    "result": "Invalid cell id!"
}
```
If you want to get the value of a cell that does not exist on the sheet, you will also get an error message:
```
GET: http://localhost:8080/api/v1/sheet1/cell12

RESPONSE: Cell cell12 is missing
```

## Input Rules for Formulas

Formulas can be entered with or without the "=" sign. Examples of formula input:

```json
{"value":"1+2"}    // OK
{"value":"=1+2"}   // OK
{"value":"=(1+2)"} // OK
{"value":"=2**2"}  // two squared | OK
{"value":"2.0+1"}  // OK
{"value":"2/0"}    // ERROR
{"value":"2str*2"} // ERROR
```
<span style="color:red">⭑</span>**There must necessarily be an operator between a variable and a number!**

### What operators and types does this support?
I use the govaluate package to calculate formulas. It supports the following operators:

* Modifiers: `+` `-` `/` `*` `&` `|` `^` `**` `%` `>>` `<<`
* Comparators: `>` `>=` `<` `<=` `==` `!=` `=~` `!~`
* Logical ops: `||` `&&`
* Numeric constants, as 64-bit floating point (`12345.678`)
* String constants (single quotes: `'foobar'`)
* Date constants (single quotes, using any permutation of RFC3339, ISO8601, ruby date, or unix date; date parsing is automatically tried with any string constant)
* Boolean constants: `true` `false`
* Parenthesis to control order of evaluation `(` `)`
* Arrays (anything separated by `,` within parenthesis: `(1, 2, 'foo')`)
* Prefixes: `!` `-` `~`
* Ternary conditional: `?` `:`
* Null coalescence: `??`

See [MANUAL.md](https://github.com/Knetic/govaluate/blob/master/MANUAL.md) for exacting details on what types each operator supports.

Full documentation: [govaluate](https://github.com/Knetic/govaluate)