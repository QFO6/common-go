package commongo

import (
	"reflect"
	"strings"
)

// ModelForm struct for form info
type ModelForm struct {
	Name     []string // [Name, Desc ...]
	Type     []string // [input, select, ....]
	Required []string // [required, "", "", required]
	Help     []string // hint
	Label    []string
	Options  map[string][]map[string]interface{} // for radio, checkbox, select, {text:text, value: value}
	Value    []interface{}
	Extra    map[string]string      //data-link="http://www.google.com"
	ViewArgs map[string]interface{} // storge all args for view
}

// ParseModelForm generate *ModelForm
func ParseModel(m interface{}, excludeFields []string) {
	fLen := len(excludeFields)
	el := reflect.TypeOf(m).Elem()
	kNum := el.NumField()
	form := new(ModelForm)
	for i := 0; i < kNum; i++ {
		fld := el.Field(i)
		// check fields loop
		if excludeFields != nil && fLen != 0 {
			exclude := false
			for _, item := range excludeFields {
				if item == fld.Name {
					exclude = true
					break
				}
			}
			if exclude {
				continue
			}
		}

		fldTagStr, found := fld.Tag.Lookup("form")
		if found {
			//name, requied, type, label, help_text
			// 0,   	1, 		2, 		3, 		4
			form.Name = append(form.Name, fld.Name)
			form.Label = append(form.Label, fld.Name)
			form.Type = append(form.Type, "input")
			form.Required = append(form.Required, "")
			form.Help = append(form.Help, "")
			fList := strings.Split(fldTagStr, ",")
			listLen := len(fList)

			//check required, index 1
			if listLen > 1 {
				required := fList[1]
				form.Required[len(form.Required)-1] = required
			}
			if listLen > 2 {
				inputType := fList[2]

				if inputType != "" && inputType != "input" {
					form.Type[len(form.Type)-1] = inputType
				}
			}
			if listLen > 3 {
				label := fList[3]
				if label != "" {
					form.Label[len(form.Label)-1] = label
				}
			}
			if listLen > 4 {
				help := fList[4]
				if help != "" {
					form.Help[len(form.Label)-1] = help
				}
			}
			v := reflect.ValueOf(m).Elem().FieldByName(fld.Name)
			if v.IsValid() {
				form.Value = append(form.Value, v.Interface())
			} else {
				form.Value = append(form.Value, "")
			}
		}
	}

	vHtml := reflect.ValueOf(m).Elem().FieldByName("ModelForm")
	if vHtml.IsValid() {
		vHtml.Set(reflect.ValueOf(form))
	}
	// check "GenOptions" exist or not
	gen := reflect.ValueOf(m).MethodByName("GenOptions")
	if gen.IsValid() {
		gen.Call([]reflect.Value{})
	}
}

// SliceToMap convert []string to [{"text":abc, "value": abc},{}...]
func SliceToMap(list []string) (out []map[string]interface{}) {
	for _, item := range list {
		tmp := map[string]interface{}{"text": item, "value": item}
		out = append(out, tmp)
	}
	return
}
