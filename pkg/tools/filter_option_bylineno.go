package tool

import (
	"nuryanto2121/dynamic_rest_api_go/models"
	"strings"
)

func FilterOptionList(OptionList []models.OptionDB, LineNo int, MethodAPI string) (Out []models.OptionDB) {

	for _, Key := range OptionList {
		if Key.LineNo == LineNo {
			Out = append(Out, Key)
		}
	}
	if MethodAPI != "" {
		var Outs []models.OptionDB
		for _, Key := range Out {
			if strings.ToLower(Key.MethodApi) == strings.ToLower(MethodAPI) {
				Outs = append(Outs, Key)
			}
		}

		return Outs

	}
	return Out
}

func FilterParamterList(ParamList []models.ParamFunction, ParamName string) (Out models.ParamFunction) {

	for _, Key := range ParamList {
		if Key.ParameterName == ParamName {
			return Key
		}
	}
	return Out
}
