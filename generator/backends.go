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
	p.P("func() error {\n")
	for _, td := range typeDescs {
		p.P(
			"// set '%s.%s' in params\n",
			s.method.GetInputTypeName(),
			td.ProtoName,
		)
		if td.IsMapped {
			p.PA([]string{
				"if params.%s, err = (%s{}).ToSpanner(req.%s).SpannerValue(); err != nil {\n",
				"return err\n}\n",
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
					"return err\n}\n",
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
					"return err\n}\n",
					"params.%s = raw\n}\n",
				},
					td.Name, td.Name, td.GoTypeName,
					td.Name,
					td.Name,
				)
			}
		} else {
			p.P("params.%s = req.%s\n", td.Name, td.Name)
		}
	}
	p.P("return nil\n}()\n")

	return p.String()
}
func (s *SpannerStringer) RowType() string {
	return "*spanner.Row"
}
func (s *SpannerStringer) TranslateRowToResult() string {
	p := &Printer{}
	p.P("func() error {\n")
	for _, td := range s.method.GetTypeDescArrayForStruct(s.method.GetOutputTypeStruct()) {
		if td.IsMapped {
			p.PA([]string{
				"var %s_ = new(spanner.GenericColumnValue)\n",
				"if err := row.ColumnByName(\"%s\", %s_); err != nil {\nreturn err\n}\n{\n",
				"local := &%s{}\n",
				"if err := local.SpannerScan(%s_); err != nil {\n return err\n}\n",
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
					"if err := row.ColumnByName(\"%s\", &%s_); err != nil {\n return err\n}\n{\n",
					"local := make(%s, len(%s_))\n",
					"for i := range local {\nlocal[i] = new(%s)\n",
					"if err := proto.Unmarshal(%s_[i], local[i]); err != nil {\n return err\n}\n}\n",
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
					"if err := row.ColumnByName(\"%s\", &%s_); err != nil {\n return err\n}\n{\n",
					"local := new(%s)\n",
					"if err := proto.Unmarshal(%s_, local); err != nil {\n return err\n}\n",
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
		} else if td.IsRepeated {
			p.PA([]string{
				"var %s_ %s\n{\n",
				"local := make(%s, 0)\n",
				"if err := row.ColumnByName(\"%s\", &local); err != nil {\n return err\n}\n",
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
				"if err := row.ColumnByName(\"%s\", local); err != nil {\n return err\n}\n",
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
	p.P("return nil\n}()\n")
	return p.String()
}

type SqlStringer struct {
	method *Method
}

func (s *SqlStringer) MapRequestToParams() string {
	p := &Printer{}
	p.P("func() error {\n")
	typeDescs := s.method.GetTypeDescArrayForStruct(s.method.GetInputTypeStruct())
	for _, td := range typeDescs {
		if td.IsMapped {
			p.P("params.%s = (%s{}).ToSql(req.%s)\n", td.Name, td.GoName, td.Name)
		} else if td.IsMessage {
			p.PA([]string{
				"if req.%s == nil {\n req.%s = new(%s) \n}\n",
				"{\nraw, err := proto.Marshal(req.%s)\nif err != nil {\n return err\n}\n",
				"params.%s = raw\n}\n",
			},
				td.Name, td.Name, td.GoTypeName,
				td.Name,
				td.Name,
			)
		} else {
			p.P("params.%s = req.%s\n", td.Name, td.Name)
		}
	}
	p.P("return nil\n}()\n")
	return p.String()
}

func (s *SqlStringer) RowType() string {
	return "persist_lib.Scanable"
}

func (s *SqlStringer) TranslateRowToResult() string {
	p := &Printer{}
	outputFields := s.method.GetTypeDescArrayForStruct(s.method.GetOutputTypeStruct())
	p.P("func() error {\n")
	for _, td := range outputFields {
		if td.IsMessage {
			p.P("var %s_ []byte\n", td.Name)
		} else {
			p.P("var %s_ %s\n", td.Name, td.GoName)
		}
	}
	p.P("if err := row.Scan(\n")
	for _, td := range outputFields {
		p.P("&%s_,\n", td.Name)
	}
	// TODO: handle sql no rows
	p.P("); err != nil {\n return err \n}\n")
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
				"if err := proto.Unmarshal(%s_, converted); err != nil {\n return err\n}\n",
				"res.%s = converted\n}\n",
			},
				td.GoTypeName,
				td.Name,
				td.Name,
			)
		} else {
			p.P("res.%s = %s_\n", td.Name, td.Name)
		}
	}
	p.P("return nil\n}()\n")
	return p.String()
}
