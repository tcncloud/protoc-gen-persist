package generator

type BackendStringer interface {
	MapRequestToParams() string
	RowType() string
	TranslateRowToResult() string
}

type SpannerStringer struct {
	method *Method
}

//TYPECHANGE
func (s *SpannerStringer) MapRequestToParams() string {
	sName := s.method.Service.GetName()
	p := &Printer{}
	typeDescs := s.method.GetTypeDescArrayForStruct(s.method.GetInputTypeStruct())
	// if value is mapped, always use the mapped value
	// if value is primitive or repeated primitive, use it
	// else convert to []byte, or [][]byte for spanner
	p.Q(
		"func ", ToParamsFuncName(s.method), "(serv ", sName, "TypeMapping, req *", s.method.GetInputType(),
		") (*persist_lib.", NewPLInputName(s.method), ", error) {\n",
	)
	p.P("var err error\n _ = err\n")
	p.P("params := &persist_lib.%s{}\n", NewPLInputName(s.method))

	for _, td := range typeDescs {
		_, titleCased := getGoNamesForTypeMapping(td.Mapping, s.method.Service.File)

		p.P("// set '%s.%s' in params\n", s.method.GetInputTypeMinusPackage(), td.ProtoName)

		if td.IsMapped {

			mappingString := P("serv.", titleCased, "()")
			p.Q("if params.", td.Name, ", err = ", mappingString, ".ToSpanner(req.", td.Name, ").SpannerValue(); err != nil {\n")
			p.Q("return nil, err\n")
			p.Q("}\n")
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
	sName := s.method.Service.GetName()
	p.P(
		"func %s(serv %sTypeMapping, row *spanner.Row) (*%s, error) {\n",
		FromScanableFuncName(s.method),
		sName,
		s.method.GetOutputType(),
	)
	p.P("res := &%s{}\n", s.method.GetOutputType())
	for _, td := range s.method.GetTypeDescArrayForStruct(s.method.GetOutputTypeStruct()) {
		_, titleCased := getGoNamesForTypeMapping(td.Mapping, s.method.Service.File)
		if td.IsMapped {
			p.Q("var ", td.Name, "_ = new(spanner.GenericColumnValue)\n")
			p.Q("if err := row.ColumnByName(\"", td.ProtoName, "\", ", td.Name, "_); err != nil {\n")
			p.Q("\treturn nil, err\n")
			p.Q("}\n{\n")
			// TYPECHAGE

			p.Q("mapper := serv.", titleCased, "()\n")
			p.Q("local := mapper.Empty()\n")
			p.Q("if err := local.SpannerScan(", td.Name, "_); err != nil {\n")
			p.Q("\treturn nil, err\n")
			p.Q("}\n")
			p.Q("if err :=.ToProto(&res.", td.Name, "); err != nil {\n")
			p.Q("\treturn nil, err\n")
			p.Q("}\n")
			// p.Q("res.", td.Name, " = mapper.ToProto(local)\n")
			p.Q("}\n")
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
	sName := s.method.Service.GetName()
	p.Q(
		"func ", ToParamsFuncName(s.method), "(serv ", sName, "TypeMapping, req *", s.method.GetInputType(),
		") (*persist_lib.", NewPLInputName(s.method), ", error) {\n",
	)
	p.P("params := &persist_lib.%s{}\n", NewPLInputName(s.method))

	typeDescs := s.method.GetTypeDescArrayForStruct(s.method.GetInputTypeStruct())
	for _, td := range typeDescs {
		_, titleCased := getGoNamesForTypeMapping(td.Mapping, s.method.Service.File)
		if td.IsMapped {
			p.Q("{\n")
			p.Q("mapper := serv.", titleCased, "()\n")
			p.Q("params.", td.Name, " = mapper.ToSql(req.", td.Name, ")\n")
			p.Q("}\n")
			// p.Q("params.", td.Name, " = s.", sName, titleCased, "(req.", td.Name, ")\n")
			// p.P("params.%s = (%s{}).ToSql(req.%s)\n", td.Name, td.GoName, td.Name)
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
	sName := s.method.Service.GetName()
	p := &Printer{}
	outputFields := s.method.GetTypeDescArrayForStruct(s.method.GetOutputTypeStruct())
	p.P(
		"func %s(serv %sTypeMapping, row persist_lib.Scanable) (*%s, error) {\n",
		FromScanableFuncName(s.method),
		sName,
		s.method.GetOutputType(),
	)
	p.P("res := &%s{}\n", s.method.GetOutputType())
	for _, td := range outputFields {
		if td.IsMessage {
			p.P("var %s_ []byte\n", td.Name)
		} else if td.IsEnum {
			p.P("var %s_ int32\n", td.Name)
		} else if td.IsMapped {
			_, titleCased := getGoNamesForTypeMapping(td.Mapping, s.method.Service.File)
			p.Q(td.Name, "_ := serv.", titleCased, "().Empty()\n")
		} else {
			p.P("var %s_ %s\n", td.Name, td.GoName)
		}
	}
	p.P("if err := row.Scan(\n")
	for _, td := range outputFields {
		if !td.IsMapped {
			p.P("&")
		}
		p.P("%s_,\n", td.Name)
	}
	p.P("); (err != nil && err != sql.ErrNoRows) {\n return nil, err \n}\n")
	for _, td := range outputFields {
		// _, titleCased := getGoNamesForTypeMapping(td.Mapping, s.method.Service.File)
		if td.IsMapped {
			p.Q("if err := ", td.Name, "_.ToProto(&res.", td.Name, "); err != nil {\n")
			p.Q("\treturn nil, err\n")
			p.Q("}\n")
			// mappingString := P("serv.", titleCased, "()")
			// p.Q("res.", td.Name, " = ", mappingString, ".ToProto(", td.Name, "_)\n")
			// p.P("res.%s = %s_.ToProto()\n", td.Name, td.Name)
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
