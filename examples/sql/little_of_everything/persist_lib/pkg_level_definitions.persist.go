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
type ExampleTable1ForExampleService1 struct {
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
}

// this could be used in a query, so generate the getters/setters
func (p *ExampleTable1ForExampleService1) GetTableId() int32              { return p.TableId }
func (p *ExampleTable1ForExampleService1) SetTableId(param int32)         { p.TableId = param }
func (p *ExampleTable1ForExampleService1) GetKey() string                 { return p.Key }
func (p *ExampleTable1ForExampleService1) SetKey(param string)            { p.Key = param }
func (p *ExampleTable1ForExampleService1) GetValue() string               { return p.Value }
func (p *ExampleTable1ForExampleService1) SetValue(param string)          { p.Value = param }
func (p *ExampleTable1ForExampleService1) GetInnerMessage() []byte        { return p.InnerMessage }
func (p *ExampleTable1ForExampleService1) SetInnerMessage(param []byte)   { p.InnerMessage = param }
func (p *ExampleTable1ForExampleService1) GetInnerEnum() int32            { return p.InnerEnum }
func (p *ExampleTable1ForExampleService1) SetInnerEnum(param int32)       { p.InnerEnum = param }
func (p *ExampleTable1ForExampleService1) GetStringArray() []string       { return p.StringArray }
func (p *ExampleTable1ForExampleService1) SetStringArray(param []string)  { p.StringArray = param }
func (p *ExampleTable1ForExampleService1) GetBytesField() []byte          { return p.BytesField }
func (p *ExampleTable1ForExampleService1) SetBytesField(param []byte)     { p.BytesField = param }
func (p *ExampleTable1ForExampleService1) GetStartTime() interface{}      { return p.StartTime }
func (p *ExampleTable1ForExampleService1) SetStartTime(param interface{}) { p.StartTime = param }
func (p *ExampleTable1ForExampleService1) GetTestField() []byte           { return p.TestField }
func (p *ExampleTable1ForExampleService1) SetTestField(param []byte)      { p.TestField = param }
func (p *ExampleTable1ForExampleService1) GetMyyenum() int32              { return p.Myyenum }
func (p *ExampleTable1ForExampleService1) SetMyyenum(param int32)         { p.Myyenum = param }
func (p *ExampleTable1ForExampleService1) GetTestsenum() int32            { return p.Testsenum }
func (p *ExampleTable1ForExampleService1) SetTestsenum(param int32)       { p.Testsenum = param }

type Test_TestForExampleService1 struct {
	Id   int32
	Name string
}

// this could be used in a query, so generate the getters/setters
func (p *Test_TestForExampleService1) GetId() int32         { return p.Id }
func (p *Test_TestForExampleService1) SetId(param int32)    { p.Id = param }
func (p *Test_TestForExampleService1) GetName() string      { return p.Name }
func (p *Test_TestForExampleService1) SetName(param string) { p.Name = param }
