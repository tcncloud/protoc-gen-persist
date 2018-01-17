package persist_lib

import (
	"database/sql"
)

type SqlClientGetter func() (*sql.DB, error)

func NewSqlClientGetter(cli *sql.DB) SqlClientGetter {
	return func() (*sql.DB, error) {
		return cli, nil
	}
}

type Scanable interface {
	Scan(dest ...interface{}) error
}
type Runable interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Exec(string, ...interface{}) (sql.Result, error)
}
type Result struct {
	result sql.Result
	row    *sql.Row
	rows   *sql.Rows
	err    error
}

func newResultFromSqlResult(r sql.Result) *Result {
	return &Result{result: r}
}
func newResultFromRow(r *sql.Row) *Result {
	return &Result{row: r}
}
func newResultFromRows(r *sql.Rows) *Result {
	return &Result{rows: r}
}
func newResultFromErr(err error) *Result {
	return &Result{err: err}
}
func (r *Result) Do(fun func(Scanable) error) error {
	if r.err != nil {
		return r.err
	}
	if r.row != nil {
		if err := fun(r.row); err != nil {
			return err
		}
	}
	if r.rows != nil {
		defer r.rows.Close()
		for r.rows.Next() {
			if err := fun(r.rows); err != nil {
				return err
			}
		}
	}
	return nil
}

// returns sql.ErrNoRows if it did not scan into dest
func (r *Result) Scan(dest ...interface{}) error {
	if r.result != nil {
		return sql.ErrNoRows
	} else if r.row != nil {
		return r.row.Scan(dest...)
	} else if r.rows != nil {
		err := r.rows.Scan(dest...)
		if r.rows.Next() {
			r.rows.Close()
		}
		return err
	}
	return sql.ErrNoRows
}
func (r *Result) Err() error {
	return r.err
}

type ExampleTable1ForTestservice1 struct {
	TableId      int32
	Key          string
	Value        string
	InnerMessage []byte
	InnerEnum    int32
	StringArray  []string
	BytesField   []byte
	StartTime    interface{}
	TestField    []byte
	Myyenum      int32
	Testsenum    int32
	Mappedenum   interface{}
}

// this could be used in a query, so generate the getters/setters
func (p *ExampleTable1ForTestservice1) GetTableId() int32               { return p.TableId }
func (p *ExampleTable1ForTestservice1) SetTableId(param int32)          { p.TableId = param }
func (p *ExampleTable1ForTestservice1) GetKey() string                  { return p.Key }
func (p *ExampleTable1ForTestservice1) SetKey(param string)             { p.Key = param }
func (p *ExampleTable1ForTestservice1) GetValue() string                { return p.Value }
func (p *ExampleTable1ForTestservice1) SetValue(param string)           { p.Value = param }
func (p *ExampleTable1ForTestservice1) GetInnerMessage() []byte         { return p.InnerMessage }
func (p *ExampleTable1ForTestservice1) SetInnerMessage(param []byte)    { p.InnerMessage = param }
func (p *ExampleTable1ForTestservice1) GetInnerEnum() int32             { return p.InnerEnum }
func (p *ExampleTable1ForTestservice1) SetInnerEnum(param int32)        { p.InnerEnum = param }
func (p *ExampleTable1ForTestservice1) GetStringArray() []string        { return p.StringArray }
func (p *ExampleTable1ForTestservice1) SetStringArray(param []string)   { p.StringArray = param }
func (p *ExampleTable1ForTestservice1) GetBytesField() []byte           { return p.BytesField }
func (p *ExampleTable1ForTestservice1) SetBytesField(param []byte)      { p.BytesField = param }
func (p *ExampleTable1ForTestservice1) GetStartTime() interface{}       { return p.StartTime }
func (p *ExampleTable1ForTestservice1) SetStartTime(param interface{})  { p.StartTime = param }
func (p *ExampleTable1ForTestservice1) GetTestField() []byte            { return p.TestField }
func (p *ExampleTable1ForTestservice1) SetTestField(param []byte)       { p.TestField = param }
func (p *ExampleTable1ForTestservice1) GetMyyenum() int32               { return p.Myyenum }
func (p *ExampleTable1ForTestservice1) SetMyyenum(param int32)          { p.Myyenum = param }
func (p *ExampleTable1ForTestservice1) GetTestsenum() int32             { return p.Testsenum }
func (p *ExampleTable1ForTestservice1) SetTestsenum(param int32)        { p.Testsenum = param }
func (p *ExampleTable1ForTestservice1) GetMappedenum() interface{}      { return p.Mappedenum }
func (p *ExampleTable1ForTestservice1) SetMappedenum(param interface{}) { p.Mappedenum = param }

type Test_TestForTestservice1 struct {
	Id   int32
	Name string
}

// this could be used in a query, so generate the getters/setters
func (p *Test_TestForTestservice1) GetId() int32         { return p.Id }
func (p *Test_TestForTestservice1) SetId(param int32)    { p.Id = param }
func (p *Test_TestForTestservice1) GetName() string      { return p.Name }
func (p *Test_TestForTestservice1) SetName(param string) { p.Name = param }
