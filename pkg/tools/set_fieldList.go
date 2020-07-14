package tool

import (
	"fmt"
	"nuryanto2121/dynamic_rest_api_go/models"
	"strings"
)

func SetFieldList(FieldSource []models.ParamFunction, DefineColumn models.DefineColumn, Len int, isList bool) (FieldQuery string, Field string, DefineSize string, FieldWhere string) {
	var (
		dfColumn []string
		// dtType   int
	)

	if Len == 0 || Len > len(FieldSource) {
		Len = len(FieldSource)
	}

	if DefineColumn.ColumnField != "" {
		dfColumn = strings.Split(DefineColumn.ColumnField, ",")
	}
	if len(dfColumn) > 0 {
		if dfColumn[0] == "no" {
			// i := len(dfColumn) - 1

			dfColumn = append(dfColumn[:0], dfColumn[1:]...) //append(dfColumn[1:]...)
			fmt.Printf("%v", dfColumn)
		}
	}
	if len(FieldSource) > 0 {
		for i, Key := range FieldSource {
			dataType := strings.ToLower(Key.DataType)
			if i <= Len {
				Field += fmt.Sprintf(",%s", Key.ParameterName)
				if len(dfColumn) > 0 && isList {
					//LIST
					if isList {
						if strings.Contains(dataType, "timestamp") || strings.Contains(dataType, "datetime") {
							FieldQuery += fmt.Sprintf(",TO_CHAR(%s,'DD/MM/YYYY HH24:MI') as %s", Key.ParameterName, Key.ParameterName)
						} else if strings.Contains(dataType, "date") {
							FieldQuery += fmt.Sprintf(",TO_CHAR(%s::DATE,'DD/MM/YYYY') as %s", Key.ParameterName, Key.ParameterName)
						} else {
							FieldQuery += fmt.Sprintf(",%s", Key.ParameterName)
						}
					}
					//LIST LookUp
					for x := 0; x < len(dfColumn); x++ {
						if strings.ToLower(dfColumn[x]) == strings.ToLower(Key.ParameterName) {
							if !isList {
								if strings.Contains(dataType, "timestamp") || strings.Contains(dataType, "datetime") {
									if strings.Contains(strings.ToLower(dfColumn[x]), "as") {
										As := strings.Split(dfColumn[x], "AS")
										FieldQuery += fmt.Sprintf(",TO_CHAR(%s,'DD/MM/YYYY HH24:MI') as %s", As[0], As[1])
									} else {
										FieldQuery += fmt.Sprintf(",TO_CHAR(%s,'DD/MM/YYYY HH24:MI') as %s", Key.ParameterName, Key.ParameterName)
									}

								} else if strings.Contains(dataType, "date") {
									if strings.Contains(strings.ToLower(dfColumn[x]), "as") {
										As := strings.Split(dfColumn[x], "AS")
										FieldQuery += fmt.Sprintf(",TO_CHAR(%s::DATE,'DD/MM/YYYY') as %s", As[0], As[1])
									} else {
										FieldQuery += fmt.Sprintf(",TO_CHAR(%s::DATE,'DD/MM/YYYY') as %s", Key.ParameterName, Key.ParameterName)
									}

								} else {
									FieldQuery += fmt.Sprintf(",%s", dfColumn[x])
								}
							}

							// Define Size
							if strings.ToLower(Key.ParameterName) == "character varying" || strings.ToLower(Key.ParameterName) == "character" || strings.ToLower(Key.ParameterName) == "char" || strings.ToLower(Key.ParameterName) == "text" {
								switch {
								case Key.MaxLength.Int64 <= 19:
									DefineSize += ",S"
								case Key.MaxLength.Int64 >= 20 && Key.MaxLength.Int64 <= 50:
									DefineSize += ",M"
								case Key.MaxLength.Int64 >= 50:
									DefineSize += ",L"
								default:
									DefineSize += ",S"
								}
							} else {
								DefineSize += ",S"
							}
							sType := ""
							if strings.Contains(dataType, "timestamp") || strings.Contains(dataType, "date") {
								sType = "T"
							} else {
								sType = "S"
							}
							FieldWhere += fmt.Sprintf(",%s:%s", dfColumn[x], sType)
						}
					}

				} else {
					if strings.Contains(dataType, "timestamp") || strings.Contains(dataType, "datetime") {
						FieldQuery += fmt.Sprintf(",TO_CHAR(%s,'DD/MM/YYYY HH24:MI') as %s", Key.ParameterName, Key.ParameterName)
					} else if strings.Contains(dataType, "date") {
						FieldQuery += fmt.Sprintf(",TO_CHAR(%s::DATE,'DD/MM/YYYY') as %s", Key.ParameterName, Key.ParameterName)
					} else {
						FieldQuery += fmt.Sprintf(",%s", Key.ParameterName)
					}

					// Define Size
					if strings.ToLower(Key.ParameterName) == "character varying" || strings.ToLower(Key.ParameterName) == "character" || strings.ToLower(Key.ParameterName) == "char" || strings.ToLower(Key.ParameterName) == "text" {
						switch {
						case Key.MaxLength.Int64 <= 19:
							DefineSize += ",S"
						case Key.MaxLength.Int64 >= 20 && Key.MaxLength.Int64 <= 50:
							DefineSize += ",M"
						case Key.MaxLength.Int64 >= 50:
							DefineSize += ",L"
						default:
							DefineSize += ",S"
						}
					} else {
						DefineSize += ",S"
					}

					sType := ""
					if strings.Contains(dataType, "timestamp") || strings.Contains(dataType, "date") {
						sType = "T"
					} else {
						sType = "S"
					}
					FieldWhere += fmt.Sprintf(",%s:%s", Key.ParameterName, sType)
				}

				// if Key.Precision != "" && Key.Scale != "" {
				// 	dtType = 1
				// } else if strings.Contains(dataType, "timestamp") || strings.Contains(dataType, "date") {
				// 	dtType = 2
				// } else {
				// 	dtType = 3
				// }

			}
		}
		Field = Field[strings.Index(Field, `,`)+1:]
		FieldQuery = FieldQuery[strings.Index(FieldQuery, `,`)+1:]
		DefineSize = DefineSize[strings.Index(DefineSize, `,`)+1:]
		FieldWhere = FieldWhere[strings.Index(FieldWhere, `,`)+1:]
	}
	return FieldQuery, Field, DefineSize, FieldWhere
}
