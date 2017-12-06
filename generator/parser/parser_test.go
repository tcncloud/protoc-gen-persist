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
}

func TestUpdateWithSet(t *testing.T) {
	reader := bytes.NewBufferString(`
UPDATE test_table set col1234='single quoted string',
affa
=
@some_field
,
      mamama = 300
PK(abba = 1, ringo_ =   false, MESSY = @gg )`)
	p := parser.NewParser(reader)
	query, err := p.Parse()
	if err != nil {
		fmt.Printf("error in parser: %s", err.Error())
		t.FailNow()
	}
	fmt.Printf("%v\n\n", query.String())
}

func TestDeleteKeyRange(t *testing.T) {
	reader := bytes.NewBufferString(`
DELETE FROM test_table START("a", 3.3, 3, @someting)
END(false, true, "sometih thing aawa", @another  ) kind(CC)
`)
	p := parser.NewParser(reader)
	query, err := p.Parse()
	if err != nil {
		fmt.Printf("error in parser: %s", err.Error())
		t.FailNow()
	}
	query.SetParams(map[string]string{
		"@someting": "\"what the\"",
		"@another":  "req.Field",
	})
	fmt.Printf("%v\n\n", query.String())
}

func TestDeleteSingle(t *testing.T) {
	reader := bytes.NewBufferString(`
delete from test_table values( @someting, @another, "a", 3, 4.0005, false, true) `)
	p := parser.NewParser(reader)
	query, err := p.Parse()
	if err != nil {
		fmt.Printf("error in parser: %s", err.Error())
		t.FailNow()
	}
	query.SetParams(map[string]string{
		"@someting": "\"what the\"",
		"@another":  "req.Field",
	})
	fmt.Printf("%v\n\n", query.String())
}

func TestSelect(t *testing.T) {
	reader := bytes.NewBufferString(`
SELECT * from table_name
WHERE a = @field_one AND b = @field_2
`)
	p := parser.NewParser(reader)
	query, err := p.Parse()
	if err != nil {
		fmt.Printf("error in parser: %s", err.Error())
		t.FailNow()
	}
	query.SetParams(map[string]string{
		"@someting": "\"what the\"",
		"@another":  "req.Field",
	})
	fmt.Printf("%v\n\n", query.String())
}
