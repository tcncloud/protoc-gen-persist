package parser

type TokenKind int

const (
	// special
	ILLEGAL TokenKind = iota
	EOF
	WS

	COMMA
	OPEN_PARAN
	CLOSE_PARAN

	//Operators
	EQUAL_OP
	GREATER_OP
	LESS_OP
	GREATER_EQUAL_OP
	LESS_EQUAL_OP

	IDENT_STRING
	IDENT_FLOAT
	IDENT_INT
	IDENT_FIELD
	IDENT_BOOL
	// kind of its own thing, it represents the literal column or table name
	IDENT_TABLE_OR_COL

	// Keywords
	INSERT
	UPDATE
	DELETE
	SELECT
	FROM
	INTO
	VALUES
	SET
	AND
	OR
	START
	END
	KIND
	CLOSED_OPEN_KIND
	CLOSED_CLOSED_KIND
	OPEN_OPEN_KIND
	OPEN_CLOSED_KIND
	PRIMARY_KEY
)

type Token struct {
	tk  TokenKind
	raw string
}

var TokenNames map[TokenKind]string = map[TokenKind]string{
	ILLEGAL:            "ILLEGAL",
	EOF:                "EOF",
	WS:                 "WS",
	COMMA:              "COMMA",
	OPEN_PARAN:         "OPEN_PARAN",
	CLOSE_PARAN:        "CLOSE_PARAN",
	EQUAL_OP:           "EQUAL_OP",
	GREATER_OP:         "GREATER_OP",
	LESS_OP:            "LESS_OP",
	GREATER_EQUAL_OP:   "GREATER_EQUAL_OP",
	LESS_EQUAL_OP:      "LESS_EQUAL_OP",
	IDENT_STRING:       "IDENT_STRING",
	IDENT_FLOAT:        "IDENT_FLOAT",
	IDENT_INT:          "IDENT_INT",
	IDENT_FIELD:        "IDENT_FIELD",
	IDENT_BOOL:         "IDENT_BOOL",
	IDENT_TABLE_OR_COL: "IDENT_TABLE_OR_COL",
	INSERT:             "INSERT",
	UPDATE:             "UPDATE",
	DELETE:             "DELETE",
	SELECT:             "SELECT",
	FROM:               "FROM",
	INTO:               "INTO",
	VALUES:             "VALUES",
	SET:                "SET",
	AND:                "AND",
	OR:                 "OR",
	START:              "START",
	END:                "END",
	KIND:               "KIND",
	CLOSED_OPEN_KIND:   "CLOSED_OPEN_KIND",
	CLOSED_CLOSED_KIND: "CLOSED_CLOSED_KIND",
	OPEN_OPEN_KIND:     "OPEN_OPEN_KIND",
	OPEN_CLOSED_KIND:   "OPEN_CLOSED_KIND",
	PRIMARY_KEY:        "PRIMARY_KEY",
}
var TokenKeyWords map[TokenKind][]string = map[TokenKind][]string{
	ILLEGAL:            []string{"ILLEGAL"},
	EOF:                []string{"EOF"},
	WS:                 []string{"any whitespace"},
	COMMA:              []string{","},
	OPEN_PARAN:         []string{"("},
	CLOSE_PARAN:        []string{")"},
	EQUAL_OP:           []string{"="},
	GREATER_OP:         []string{">"},
	LESS_OP:            []string{"<"},
	GREATER_EQUAL_OP:   []string{">="},
	LESS_EQUAL_OP:      []string{"<="},
	IDENT_STRING:       []string{"'abc 123'", "\"abc_123\""},
	IDENT_FLOAT:        []string{"1.00", "0.5"},
	IDENT_INT:          []string{"123"},
	IDENT_FIELD:        []string{"@fieldName", "@field_name"},
	IDENT_BOOL:         []string{"true", "false"},
	IDENT_TABLE_OR_COL: []string{"table_name", "columnName", "(notice lack of quotes)"},
	INSERT:             []string{"INSERT", "insert"},
	UPDATE:             []string{"UPDATE", "update"},
	DELETE:             []string{"DELETE", "delete"},
	SELECT:             []string{"SELECT", "select"},
	FROM:               []string{"FROM", "from"},
	INTO:               []string{"INTO", "into"},
	VALUES:             []string{"VALUES", "values"},
	SET:                []string{"SET", "set"},
	AND:                []string{"AND", "and"},
	OR:                 []string{"OR", "or"},
	START:              []string{"START", "start"},
	END:                []string{"END", "end"},
	KIND:               []string{"KIND", "kind"},
	CLOSED_OPEN_KIND: []string{
		"CLOSED_OPEN", "closed_open", "ClosedOpen",
		"closedOpen", "CO", "co",
	},
	CLOSED_CLOSED_KIND: []string{
		"CLOSED_CLOSED", "closed_closed", "ClosedClosed",
		"closedClosed", "CC", "cc",
	},
	OPEN_OPEN_KIND: []string{
		"OPEN_OPEN", "open_open", "OpenOpen",
		"openOpen", "OO", "oo",
	},
	OPEN_CLOSED_KIND: []string{
		"OPEN_CLOSED", "open_closed", "OpenClosed",
		"openClosed", "OC", "oc",
	},
	PRIMARY_KEY: []string{
		"PRIMARY_KEY", "primary_key", "PrimaryKey",
		"primaryKey", "PK", "pk",
	},
}
