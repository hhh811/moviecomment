package utils

import (
	"fmt"
	"html/template"
	"reflect"
	"strings"

	"time"

	"github.com/astaxie/beego/validation"
)

func init() {
	//TODO form init

	initCommonField()
}

type FieldCreater func(*FieldSet)

type FieldFilter func(*FieldSet)

var customCreaters = make(map[string]FieldCreater)

var customFilters = make(map[string]FieldFilter)

type HtmlLazyField func() template.HTML

func RegisterFieldCreater(name string, field FieldCreater) {
	customCreaters[name] = field
}

func RegisterFieldFilter(name string, field FieldFilter) {
	customFilters[name] = field
}

func (f HtmlLazyField) String() string {
	return string(f())
}

type FieldSet struct {
	Label       template.HTML
	Field       HtmlLazyField
	Id          string
	Name        string
	LabelText   string
	Value       interface{}
	Help        string
	Error       string
	Type        string
	Kind        string
	Placeholder string
	Attrs       string
	FormElm     reflect.Value
	// TODO localer
}

type FormSets struct {
	FieldList []*FieldSet
	Fields    map[string]*FieldSet
	inited    bool
	form      interface{}
	errs      map[string]*validation.Error
}

// create formSets for generate label/field html code
func NewFormSets(form interface{}, errs map[string]*validation.Error) *FormSets {
	fSets := new(FormSets)
	fSets.errs = errs
	fSets.Fields = make(map[string]*FieldSet)
	// TODO locale

	val := reflect.ValueOf(form)

	panicAssertStructPtr(val)

	elm := val.Elem()

	// TODO helps labels and places

outFor:
	for i := 0; i < elm.NumField(); i++ {
		f := elm.Field(i)
		fT := elm.Type().Field(i)

		name := fT.Name
		value := f.Interface()
		fTyp := "text"

		switch f.Kind() {
		case reflect.Bool:
			fTyp = "checkbox"
		default:
			switch value.(type) {
			case time.Time:
				fTyp = "datetime"
			}
		}

		fName := name

		var attrm map[string]string

		// parse struct tag settings
		for _, v := range strings.Split(fT.Tag.Get("form"), ";") {
			v = strings.TrimSpace(v)
			if v == "-" {
				continue outFor
			} else if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
				tN := v[:i]
				v = strings.TrimSpace(v[i+1 : len(v)-1])
				switch tN {
				case "type":
					fTyp = v
				case "name":
					fName = v
				case "attr":
					if attrm == nil {
						attrm = make(map[string]string)
					}
					parts := strings.SplitN(v, ",", 2)
					if len(parts) > 1 {
						attrm[parts[0]] = parts[1]
					} else {
						attrm[v] = v
					}
				}
			}
		}

		var attrs string
		for k, v := range attrm {
			attrs += fmt.Sprintf(` %s="%s"`, k, v)
		}

		// set field id
		fId := elm.Type().Name() + "-" + fName

		// one field in form, one FieldSet
		var fSet FieldSet

		fSet.Id = fId
		fSet.Name = fName
		fSet.Value = value
		fSet.Attrs = attrs
		fSet.FormElm = elm

		if i := strings.IndexRune(fTyp, ','); i != -1 {
			fSet.Type = fTyp[:i]
			fSet.Kind = fTyp[i+1:]
			fTyp = fSet.Type
		} else {
			fSet.Type = fTyp
			fSet.Kind = fTyp
		}

		// TODO get field label text
		// TODO field help
		// TODO field palce holder
		// TODO error string

		// create label html
		switch fTyp {
		case "checkbox", "hidden":
		default:
			fSet.Label = template.HTML(fmt.Sprintf(`
				<label class="control-label" for="%s">%s</label>`, fSet.Id, fSet.LabelText))
		}
		// 设置默认 Field() 用于生成默认 html
		if creater, ok := customCreaters[fTyp]; ok {
			creater(&fSet)
		}

		if filter, ok := customFilters[fTyp]; ok {
			filter(&fSet)
		}

		if fSet.Field == nil {
			fSet.Field = func() template.HTML { return "" }
		}

		fSets.FieldList = append(fSets.FieldList, &fSet)
		fSets.Fields[name] = &fSet

	}

	fSets.inited = true

	return fSets
}

func initCommonField() {
	RegisterFieldCreater("text", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="text" value="%v" class="form-control"%s%s>`,
				fSet.Id, fSet.Name, fSet.Value, fSet.Placeholder, fSet.Attrs))
		}
	})

	RegisterFieldCreater("textarea", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<textarea id="%s", name="%s" rows="5" class="form-control"%s%s>%v</textarea>`,
				fSet.Id, fSet.Name, fSet.Placeholder, fSet.Attrs, fSet.Value))
		}
	})

	RegisterFieldCreater("password", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="password" value="%v" class="form-control"%s%s`,
				fSet.Id, fSet.Name, fSet.Value, fSet.Placeholder, fSet.Attrs))
		}
	})

	RegisterFieldCreater("hidden", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			return template.HTML(fmt.Sprintf(
				`<input id="%s" name="%s" type="hidden" value="%v"%s>`, fSet.Id, fSet.Name, fSet.Value, fSet.Attrs))
		}
	})

	// TODO datetime 处理

	RegisterFieldCreater("checkbox", func(fSet *FieldSet) {
		fSet.Field = func() template.HTML {
			var checked string
			if b, ok := fSet.Value.(bool); ok && b {
				checked = "checked"
			}
			return template.HTML(fmt.Sprintf(
				`<label for="%s" class="checkbox">%s<input id="%s" name="%s" type="checkbox" %s></label>`,
				fSet.Id, fSet.LabelText, fSet.Id, fSet.Name, checked))
		}
	})

	// TODO select 表单
}

// assert an object must be a struct pointer
func panicAssertStructPtr(val reflect.Value) {
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		return
	}
	panic(fmt.Errorf("%s must be a struct pointer", val.Type().Name()))
}
