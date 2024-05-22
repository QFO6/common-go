package commongo

import (
	"reflect"
	"strings"
)

type HtmlForm struct {
	Name     []string // [Name, Desc ...]
	Type     []string // [input, select, ....]
	Required []string // ["required", "", "", "required"]
	Help     []string
	Label    []string
	Value    []interface{}
	Options  map[string][][]interface{}          // options for field
	Extra    map[string]string                   // eg: disable
	Attrs    map[string]map[string]string        // attr for field
	OptAttrs map[string][]map[string]interface{} // data-link="abc" in option only for select or options
}

// ParseModelForm generate *HtmlForm
func ParseModelForm(m interface{}) {
	el := reflect.TypeOf(m).Elem()
	kNum := el.NumField()
	form := new(HtmlForm)
	for i := 0; i < kNum; i++ {
		fld := el.Field(i)
		fldTagStr, found := fld.Tag.Lookup("form")
		if found {
			//name, required, type, label, help_text
			// 0,   	1, 		2, 		3, 		4
			// init
			form.Name = append(form.Name, fld.Name)
			form.Label = append(form.Label, fld.Name)
			form.Type = append(form.Type, "input")
			form.Required = append(form.Required, "")
			form.Help = append(form.Help, "")
			fList := strings.Split(fldTagStr, ",")
			listLen := len(fList)

			// check required, index 1
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

	vHtml := reflect.ValueOf(m).Elem().FieldByName("HtmlForm")
	if vHtml.IsValid() {
		// tHtml.Set(reflect.ValueOf(hMap))
		vHtml.Set(reflect.ValueOf(form))
	}
	// check "GenOptions" exist or not
	gen := reflect.ValueOf(m).MethodByName("GenOptions")
	if gen.IsValid() {
		gen.Call([]reflect.Value{})
	}
}
