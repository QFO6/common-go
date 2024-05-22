package commongo

import "reflect"

// FormFields get fields names defined in Struct with "form" tag
func FormFields(m interface{}) []string {
	el := reflect.TypeOf(m).Elem()
	kNum := el.NumField()
	output := []string{}
	for i := 0; i < kNum; i++ {
		fld := el.Field(i)
		_, found := fld.Tag.Lookup("form")
		if found {
			output = append(output, fld.Name)
		}
	}
	return output
}

// ModelFields return model fields name
func ModelFields(m interface{}) []string {
	//m, _ := ModelReg(modelName)
	el := reflect.TypeOf(m).Elem()
	kNum := el.NumField()
	output := []string{}
	for i := 0; i < kNum; i++ {
		fld := el.Field(i)
		tag, found := fld.Tag.Lookup("bson")
		_, hidden := fld.Tag.Lookup("hidden")
		if !hidden && found && tag != ",inline" && tag != "-" {
			output = append(output, fld.Name)
		}
	}
	return output
}
