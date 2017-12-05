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
	fmt.Printf("%v\n\n", query.String())

}

func TestUpdate(t *testing.T) {
	reader := bytes.NewBufferString(`
		update test_table
			(col1234, affa, mamama)
		VALUES
			('single quoted string', @some_field, 300)
		PRIMARY_KEY(abba = 1, ringo_=false , MESSY   =  @gg
		)
	`)
	p := parser.NewParser(reader)
	query, err := p.Parse()
	if err != nil {
		fmt.Printf("error in parser: %s", err.Error())
		t.FailNow()
	}
	fmt.Printf("%v\n\n", query.String())

	for _, t := range query.Tokens() {
		fmt.Printf("%s,  %s\n", t.Name(), t.Raw())
	}
	fmt.Printf("\nFields:\n")
	for _, t := range query.Fields() {
		fmt.Printf("%s\n", t)
	}
	fmt.Printf("\nARGS:\n")
	for _, t := range query.Args() {
		fmt.Printf("%s,  %s\n", t.Name(), t.Raw())
	}
}
