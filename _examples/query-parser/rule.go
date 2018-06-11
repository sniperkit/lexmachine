package main

/*
	Refs:
	- https://regex-golang.appspot.com/assets/html/index.html # test reguler expressions
*/

// Select query - regular expression patterns
const (

	// `SELECT_BLOCK` speficies regex to extract selection a block of cells. eg. `A2:B4`
	SELECT_BLOCK = `([A-Za-z]+[0-9]+)\:([A-Za-z]+[0-9]+)`

	// `SELECT_CELL` speficies regex to extract selection a specific cell. eg. `A2`
	SELECT_CELL = `([A-Za-z]+)([0-9]+)`

	// `SELECT_CELL` speficies regex to extract selection of columns. eg `A:B`
	SELECT_COLS = `([A-Za-z]+)\:([A-Za-z]+)`

	// `SELECT_ROWS` speficies regex to extract selection of rows. eg `2:5`
	SELECT_ROWS = `([0-9]+)\:([0-9]+)`

	// `SELECT_COL` speficies regex to extract selection of specific column. eg `A`
	SELECT_COL = `([A-Za-z]+)`

	// `SELECT_ROW` speficies regex to extract selection of specific row. eg `5`
	SELECT_ROW = `([0-9]+)`

	// `SELECT_NUMERIC` speficies regex to extract selection of columns and rows matched by their key index.
	// eg. `cols[0:5],rows[1:7]`, `cols[:],rows[:]`, `cols[1,2],rows[:]`
	SELECT_NUMERIC = `((col|cols|rows|row))\[((:\d+)|(\d+\:)|(\d+\:\d+)|(\:)|(\d+(,\d+)*))(\])`

	// `SELECT_ALPHA_NUMERIC` speficies regex to extract selection of columns and rows matched by their key index or column name.
	// Examples:
	// - `col["name", "full_name"], rows[1,5,7,8]`,
	// - `rows[1,10], col["id", "name", "full_name", "description", "language", "stargazers_count", "forks_count"]`
	SELECT_ALPHA_NUMERIC = `((col|cols|rows|row))\s\[((:[a-zA-Z0-9-_\"\']+)|([a-zA-Z0-9-_\"\']+\:)|([a-zA-Z0-9-_\"\']+\:[a-zA-Z0-9-_\"\']+)|([a-zA-Z0-9-_\"\']+(,[a-zA-Z0-9-_\"\']+)*)|(\:)|([a-zA-Z0-9-_\"\']+(,[a-zA-Z0-9-_\"\']+)*))(\])`

	// `SELECT_FUNCTION` speficies regex to extract selection of columns and rows matched by arguments.
	// Examples:
	// - `SELECT(col=1, row=2)`
	// - `SELECT(cols=[:], rows=[::5])`
	// - `SELECT(cols=["name","full_name"], rows=[::100])`
	// - `SELECT(cols=[1,2,7,10], row=2)`
	// - `SELECT(cols=[1,2,7,10], row=2)`
	// - `SELECT(cols=[1,2,7,10], rows=[2,4,7])`
	// - `SELECT(cols=[:], rows=[2,4,7])`
	// - `SELECT(cols=["name","full_name"], rows=[2,4,7])`
	// - `SELECT(cols=[0:5], rows=[1:7])`,
	// - `SELECT(cols=[:],rows=[:])`,
	// - `SELECT(cols=[1,2],rows=[:])`
	// - `SELECT(cols=["name","full_name"], rows=[1,5,7,8])`,
	// - `SELECT(rows=[1,10], cols=["id", "name", "full_name", "description", "language", "stargazers_count", "forks_count"])`
	SELECT_FUNCTION = `((col|cols|rows|row))\s\=\s\(((:\d+)|(\d+\:)|(\d+\:\d+)|(\:)|(\d+(,\d+)*))(\s\))`

	//
	// ref.
	SELECT_UNICODE = `^[\\p{L}0-9]*$`

	//-- End
)
