package parser_test

import (
	"bytes"
	"fmt"
	"github.com/tcncloud/protoc-gen-persist/generator/parser"
	"testing"
)

func TestInsert(t *testing.T) {
	reader := bytes.NewBufferString(`
		INSERT INTO test_table
			(col1, col2, col3)
			VALUES
				(@field1, "string2", 3.3)
	`)
	p := parser.NewParser(reader)
	query, err := p.Parse()
	if err != nil {
		fmt.Printf("error in parser: %s", err.Error())
		t.FailNow()
	}
	fmt.Println(query.String())

}
