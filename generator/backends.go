package generator

type BackendStringer interface {
	MapRequestToParams() string
	RowType() string
	TranslateRowToResult() string
}

type SpannerStringer struct {
	method *Method
}

func (s *SpannerStringer) MapRequestToParams() string {
	p := &Printer{}
	typeDescs := s.method.GetTypeDescArrayForStruct(s.method.GetInputTypeStruct())
	// if value is mapped, always use the mapped value
	// if value is primitive or repeated primitive, use it
	// else convert to []byte, or [][]byte for spanner
	p.P(
		"func %s(req *%s) (*persist_lib.%s, error) {\n",
		ToParamsFuncName(s.method),
		s.method.GetInputType(),
		NewPLInputName(s.method),
	)
	p.P("var err error\n _ = err\n")
	p.P("params := &persist_lib.%s{}\n", NewPLInputName(s.method))
	for _, td := range typeDescs {
		p.P(
			"// set '%s.%s' in params\n",
			s.method.GetInputTypeMinusPackage(),
			td.ProtoName,
		)
		if td.IsMapped {
			p.PA([]string{
				"if params.%s, err = (%s{}).ToSpanner(req.%s).SpannerValue(); err != nil {\n",
				"return nil, err\n}\n",
			},
				td.Name,
				td.GoName,
				td.Name,
			)
		} else if td.IsMessage {
			if td.IsRepeated {
				p.PA([]string{
					"{\nvar bytesOfBytes [][]byte\n",
					"for _, msg := range req.%s{\n",
					"raw, err := proto.Marshal(msg)\nif err != nil {\n",
					"return nil, err\n}\n",
					"bytesOfBytes = append(bytesOfBytes, raw)\n}\n",
					"params.%s = bytesOfBytes\n}\n",
				},
					td.Name,
					td.Name,
				)
			} else {
				p.PA([]string{
					"if req.%s == nil {\n req.%s = new(%s)\n}\n",
					"{\nraw, err := proto.Marshal(req.%s)\nif err != nil {\n",
					"return nil, err\n}\n",
					"params.%s = raw\n}\n",
				},
					td.Name, td.Name, td.GoTypeName,
					td.Name,
					td.Name,
				)
			}
		} else if td.IsEnum {
			p.P("params.%s = int32(req.%s)\n", td.Name, td.Name)
		} else {
			p.P("params.%s = req.%s\n", td.Name, td.Name)
		}
	}
	p.P("return params, nil\n}\n")

	return p.String()
}

func (s *SpannerStringer) RowType() string {
	return "*spanner.Row"
}
func (s *SpannerStringer) TranslateRowToResult() string {
	p := &Printer{}
	p.P(
		"func %s(row *spanner.Row) (*%s, error) {\n",
		FromScanableFuncName(s.method),
		s.method.GetOutputType(),
	)
	p.P("res := &%s{}\n", s.method.GetOutputType())
	for _, td := range s.method.GetTypeDescArrayForStruct(s.method.GetOutputTypeStruct()) {
		if td.IsMapped {
			p.PA([]string{
				"var %s_ = new(spanner.GenericColumnValue)\n",
				"if err := row.ColumnByName(\"%s\", %s_); err != nil {\nreturn nil, err\n}\n{\n",
				"local := &%s{}\n",
				"if err := local.SpannerScan(%s_); err != nil {\n return nil, err\n}\n",
				"res.%s = local.ToProto()\n}\n",
			},
				td.Name,
				td.ProtoName,
				td.Name,
				td.GoName,
				td.Name,
				td.Name,
			)
		} else if td.IsMessage {
			// this is super tacky.  But I can be sure I need this import at this point
			s.method.
				Service.
				File.ImportList.GetOrAddImport("proto", "github.com/golang/protobuf/proto")
			if td.IsRepeated {
				p.PA([]string{
					"var %s_ [][]byte\n",
					"if err := row.ColumnByName(\"%s\", &%s_); err != nil {\n return nil, err\n}\n{\n",
					"local := make(%s, len(%s_))\n",
					"for i := range local {\nlocal[i] = new(%s)\n",
					"if err := proto.Unmarshal(%s_[i], local[i]); err != nil {\n return nil, err\n}\n}\n",
					"res.%s = local\n}\n",
				},
					td.Name,
					td.ProtoName,
					td.Name,
					td.GoName,
					td.Name,
					td.GoTypeName,
					td.Name,
					td.Name,
				)
			} else {
				p.PA([]string{
					"var %s_ []byte\n",
					"if err := row.ColumnByName(\"%s\", &%s_); err != nil {\n return nil, err\n}\n{\n",
					"local := new(%s)\n",
					"if err := proto.Unmarshal(%s_, local); err != nil {\n return nil, err\n}\n",
					"res.%s = local\n}\n",
				},
					td.Name,
					td.ProtoName,
					td.Name,
					td.GoTypeName,
					td.Name,
					td.Name,
				)
			}
		} else if td.IsEnum {
			if td.IsRepeated {
				// TODO: UNSUPPORTED YET
			} else {
				// even though we scan them in as int32, we scan out of spanner as int64
				// they should always fit in an int32 though,
				p.PA([]string{
					"var %s_ int64\n",
					"if err := row.ColumnByName(\"%s\", &%s_); err != nil {\n return nil, err\n}\n",
					"res.%s = %s(%s_)\n",
				},
					td.Name,
					td.ProtoName, td.Name,
					td.Name, td.GoTypeName, td.Name,
				)
			}
		} else if td.IsRepeated {
			p.PA([]string{
				"var %s_ %s\n{\n",
				"local := make(%s, 0)\n",
				"if err := row.ColumnByName(\"%s\", &local); err != nil {\n return nil, err\n}\n",
				"for _, l := range local {\nif l.Valid {\n",
				"%s_ = append(%s_, l.%s)\n",
				"res.%s = %s_\n}\n}\n}\n",
			},
				td.Name,
				td.GoName,
				td.SpannerType,
				td.ProtoName,
				td.Name,
				td.Name,
				td.SpannerTypeFieldName,
				td.Name,
				td.Name,
			)
		} else {
			p.PA([]string{
				"var %s_ %s\n{\nlocal := &%s{}\n",
				"if err := row.ColumnByName(\"%s\", local); err != nil {\n return nil, err\n}\n",
				"if local.Valid {\n %s_ = local.%s\n}\n",
				"res.%s = %s_\n}\n",
			},
				td.Name,
				td.GoName,
				td.SpannerType,
				td.ProtoName,
				td.Name,
				td.SpannerTypeFieldName,
				td.Name,
				td.Name,
			)
		}
	}
	p.P("return res, nil\n}\n")
	return p.String()
}

type SqlStringer struct {
	method *Method
}

func (s *SqlStringer) MapRequestToParams() string {
	p := &Printer{}
	p.P(
		"func %s(req *%s) (*persist_lib.%s, error) {\n",
		ToParamsFuncName(s.method),
		s.method.GetInputType(),
		NewPLInputName(s.method),
	)
	p.P("params := &persist_lib.%s{}\n", NewPLInputName(s.method))

	typeDescs := s.method.GetTypeDescArrayForStruct(s.method.GetInputTypeStruct())
	for _, td := range typeDescs {
		if td.IsMapped {
			p.P("params.%s = (%s{}).ToSql(req.%s)\n", td.Name, td.GoName, td.Name)
		} else if td.IsMessage {
			p.PA([]string{
				"if req.%s == nil {\n req.%s = new(%s) \n}\n",
				"{\nraw, err := proto.Marshal(req.%s)\nif err != nil {\n return nil, err\n}\n",
				"params.%s = raw\n}\n",
			},
				td.Name, td.Name, td.GoTypeName,
				td.Name,
				td.Name,
			)
		} else if td.IsEnum {
			p.P("params.%s = int32(req.%s)\n", td.Name, td.Name)
		} else {
			p.P("params.%s = req.%s\n", td.Name, td.Name)
		}
	}
	p.P("return params, nil\n}\n")
	return p.String()
}

func (s *SqlStringer) RowType() string {
	return "persist_lib.Scanable"
}

func (s *SqlStringer) TranslateRowToResult() string {
	p := &Printer{}
	outputFields := s.method.GetTypeDescArrayForStruct(s.method.GetOutputTypeStruct())
	p.P(
		"func %s(row persist_lib.Scanable) (*%s, error) {\n",
		FromScanableFuncName(s.method),
		s.method.GetOutputType(),
	)
	p.P("res := &%s{}\n", s.method.GetOutputType())

	for _, td := range outputFields {
		if td.IsMessage {
			p.P("var %s_ []byte\n", td.Name)
		} else if td.IsEnum {
			p.P("var %s_ int32\n", td.Name)
		} else {
			p.P("var %s_ %s\n", td.Name, td.GoName)
		}
	}
	p.P("if err := row.Scan(\n")
	for _, td := range outputFields {
		p.P("&%s_,\n", td.Name)
	}
	p.P("); (err != nil && err != sql.ErrNoRows) {\n return nil, err \n}\n")
	for _, td := range outputFields {
		if td.IsMapped {
			p.P("res.%s = %s_.ToProto()\n", td.Name, td.Name)
		} else if td.IsMessage {
			// this is super tacky.  But I can be sure I need this import at this point
			s.method.
				Service.
				File.ImportList.GetOrAddImport("proto", "github.com/golang/protobuf/proto")
			p.PA([]string{
				"{\n var converted = new(%s)\n",
				"if err := proto.Unmarshal(%s_, converted); err != nil {\n return nil, err\n}\n",
				"res.%s = converted\n}\n",
			},
				td.GoTypeName,
				td.Name,
				td.Name,
			)
		} else if td.IsEnum {
			p.P("res.%s = %s(%s_)\n", td.Name, td.GoTypeName, td.Name)
		} else {
			p.P("res.%s = %s_\n", td.Name, td.Name)
		}
	}
	p.P("return res, nil\n}\n")
	return p.String()
}
