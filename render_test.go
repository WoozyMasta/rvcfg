package rvcfg

import "testing"

func TestRenderFile(t *testing.T) {
	t.Parallel()

	file := File{
		Statements: []Statement{
			{
				Kind: NodeClass,
				Class: &ClassDecl{
					Name: "CfgModels",
					Body: []Statement{
						{
							Kind: NodeClass,
							Class: &ClassDecl{
								Name: "box",
								Base: "Default",
								Body: []Statement{
									{
										Kind: NodeProperty,
										Property: &PropertyAssign{
											Name:  "skeletonName",
											Value: Value{Kind: ValueScalar, Raw: "\"Skeleton\""},
										},
										TrailingComments: []Comment{
											{Text: "// Memory points selections:"},
											{Text: "// lid_front_axis"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	got, err := RenderFile(file)
	if err != nil {
		t.Fatalf("RenderFile: %v", err)
	}

	want := "" +
		"class CfgModels\n" +
		"{\n" +
		"  class box: Default\n" +
		"  {\n" +
		"    skeletonName = \"Skeleton\";\n" +
		"    // Memory points selections:\n" +
		"    // lid_front_axis\n" +
		"  };\n" +
		"};\n"

	if string(got) != want {
		t.Fatalf("unexpected render output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestRenderFileWithOptions(t *testing.T) {
	t.Parallel()

	file := File{
		Statements: []Statement{
			{
				Kind: NodeClass,
				Class: &ClassDecl{
					Name: "Cfg",
					Body: []Statement{
						{
							Kind: NodeProperty,
							Property: &PropertyAssign{
								Name:  "value",
								Value: Value{Kind: ValueScalar, Raw: "1"},
							},
							TrailingComments: []Comment{
								{Text: "// trailing"},
							},
						},
					},
				},
			},
		},
	}

	got, err := RenderFileWithOptions(file, FormatOptions{
		IndentChar:       "\t",
		IndentSize:       1,
		PreserveComments: false,
	})
	if err != nil {
		t.Fatalf("RenderFileWithOptions: %v", err)
	}

	want := "" +
		"class Cfg\n" +
		"{\n" +
		"\tvalue = 1;\n" +
		"};\n"

	if string(got) != want {
		t.Fatalf("unexpected render output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}

func TestRenderFileWithOptionsWrapByArrayName(t *testing.T) {
	t.Parallel()

	file := File{
		Statements: []Statement{
			{
				Kind: NodeClass,
				Class: &ClassDecl{
					Name: "CfgSkeletons",
					Body: []Statement{
						{
							Kind: NodeClass,
							Class: &ClassDecl{
								Name: "TestSkel",
								Body: []Statement{
									{
										Kind: NodeArrayAssign,
										ArrayAssign: &ArrayAssign{
											Name: "SkeletonBones",
											Value: Value{
												Kind: ValueArray,
												Elements: []Value{
													{Kind: ValueScalar, Raw: "\"a\""},
													{Kind: ValueScalar, Raw: "\"\""},
													{Kind: ValueScalar, Raw: "\"b\""},
													{Kind: ValueScalar, Raw: "\"\""},
												},
											},
										},
									},
									{
										Kind: NodeArrayAssign,
										ArrayAssign: &ArrayAssign{
											Name: "sections",
											Value: Value{
												Kind: ValueArray,
												Elements: []Value{
													{Kind: ValueScalar, Raw: "\"x\""},
													{Kind: ValueScalar, Raw: "\"y\""},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	got, err := RenderFileWithOptions(file, FormatOptions{
		IndentChar:             " ",
		IndentSize:             2,
		MaxLineWidth:           0,
		MaxInlineArrayElements: 0,
		ArrayWrapByName: map[string]int{
			"SkeletonBones": 2,
		},
		PreserveComments: true,
	})
	if err != nil {
		t.Fatalf("RenderFileWithOptions: %v", err)
	}

	want := "" +
		"class CfgSkeletons\n" +
		"{\n" +
		"  class TestSkel\n" +
		"  {\n" +
		"    SkeletonBones[] =\n" +
		"    {\n" +
		"      \"a\", \"\",\n" +
		"      \"b\", \"\",\n" +
		"    };\n" +
		"    sections[] = {\"x\", \"y\"};\n" +
		"  };\n" +
		"};\n"

	if string(got) != want {
		t.Fatalf("unexpected render output\nwant:\n%s\ngot:\n%s", want, string(got))
	}
}
