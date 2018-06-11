package main

var (

	// queryStrings represents a slice of query strings to test
	queryStrings = []string{
		////////////////////////////////////////////////////////////////////////
		// Select a column, multiple columns by name or index/indicies (Axis-X)
		//
		// Schema:
		// - col(index)
		// - col(name)
		// - cols(lower:upper:cap)
		// - cols(index_1, index_3, index_n...)
		// - cols(name_1, name_5, name_n...)
		//
		////////////////////////////////////////////////////////////////////////

		////////////////////////////////////
		// UNIQUE COLUMN
		////////////////////////////////////

		// unique column by name
		`col["name"]`,

		// unique column by index
		`col[2]`,

		////////////////////////////////////
		// COLUMNS
		////////////////////////////////////

		// columns by indices list
		`cols[1,2]`,
		`cols[1,5,7]`,

		// columns by headers list name
		`cols["name", "fullname"]`,
		`cols["name", "full_name"]`,

		// columns by index range
		`cols[:]`,
		`cols[1:]`,
		`cols[1:7]`,
		`cols[:10]`,

		////////////////////////////////////
		// MIXED
		////////////////////////////////////
		// selection queries for columns and rows
		`cols[0:5],rows[1:7]`,
		`cols[:],rows[:]`,
		`cols[1,2],rows[:]`,
		`col[name, full_name],row[4]`,
		`col["name", "full_name"],rows[1,5,7,8]`,
		`col["name", "full_name"],rows[:10]`,
		`rows[:10],col["name", "full_name"]`,
		`rows[1,10],col["name", "fullname"]`,
		`rows[1,10],col["id", "name", "full_name", "description", "language", "stargazers_count", "forks_count"]`,

		////////////////////////////////////////////////////////////////////////
		// Select a row, multiple rows by index/indicies (Axis-Y)
		//
		// Schema:
		// - row(index)
		// - rows(lower:upper:cap)
		// - rows(index_1, index_3, index_n...)
		//
		// Important: row(s) are ONLY selected with NUMERIC index
		////////////////////////////////////////////////////////////////////////

		////////////////////////////////////
		// UNIQUE ROW
		////////////////////////////////////

		// unique row by index
		`row[1]`,

		////////////////////////////////////
		// ROWS
		////////////////////////////////////

		// rows by indices list
		`rows[1,2]`,
		`rows[1,5,7]`,

		// rows by index range
		`rows[:]`,
		`rows[1:]`,
		`rows[1:7]`,
		`rows[:25]`,
	}

	//-- End
)
