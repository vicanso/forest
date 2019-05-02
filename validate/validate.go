package validate

import (
	"regexp"
	"strings"

	"github.com/asaskevich/govalidator"
	jsoniter "github.com/json-iterator/go"
	"github.com/vicanso/hes"
)

var (
	paramTagRegexMap = govalidator.ParamTagRegexMap
	paramTagMap      = govalidator.ParamTagMap
	customTypeTagMap = govalidator.CustomTypeTagMap
	json             = jsoniter.ConfigCompatibleWithStandardLibrary

	errCategory = "validate"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	AddRegex("xIntRange", "^xIntRange\\((\\d+)\\|(\\d+)\\)$", func(value string, params ...string) bool {
		return govalidator.InRangeInt(value, params[0], params[1])
	})

	AddRegex("xIntIn", `^xIntIn\((.*)\)$`, func(value string, params ...string) bool {
		if len(params) == 1 {
			rawParams := params[0]
			parsedParams := strings.Split(rawParams, "|")
			return govalidator.IsIn(value, parsedParams...)
		}
		return false
	})

	methods := []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"HEAD",
	}
	Add("xMethods", func(i interface{}, _ interface{}) bool {
		value, ok := i.(string)
		if !ok {
			return false
		}
		return govalidator.IsIn(value, methods...)
	})

	Add("xAccount", func(i interface{}, _ interface{}) bool {
		return checkASCIILength(i, 2, 20)
	})

	Add("xPassword", func(i interface{}, _ interface{}) bool {
		return checkASCIILength(i, 6, 64)
	})
}

func checkASCIILength(i interface{}, min, max int) bool {
	value, ok := i.(string)
	if !ok {
		return false
	}
	size := len(value)
	// ascii 而且 长度>=min <=20
	if !govalidator.IsASCII(value) || size < min || size > max {
		return false
	}
	return true
}

func wrapError(e error) error {
	err := hes.Wrap(e)
	err.Category = errCategory
	return err
}

// Do do validate
func Do(s interface{}, data interface{}) (err error) {
	// statusCode := http.StatusBadRequest
	if data != nil {
		switch data.(type) {
		case []byte:
			e := json.Unmarshal(data.([]byte), s)
			if e != nil {
				err = wrapError(e)
				return
			}
		default:
			buf, e := json.Marshal(data)
			if e != nil {
				err = wrapError(e)
				return
			}
			e = json.Unmarshal(buf, s)
			if e != nil {
				err = wrapError(e)
				return
			}
		}
	}
	_, e := govalidator.ValidateStruct(s)
	if e != nil {
		err = wrapError(e)
	}
	return
}

// AddRegex add a regexp validate
func AddRegex(name, reg string, fn govalidator.ParamValidator) {
	if paramTagMap[name] != nil {
		panic(name + ", reg:" + reg + " is duplicated")
	}
	paramTagRegexMap[name] = regexp.MustCompile(reg)
	paramTagMap[name] = fn
}

// Add add validate
func Add(name string, fn govalidator.CustomTypeValidator) {
	_, exists := customTypeTagMap.Get(name)
	if exists {
		panic(name + " is duplicated")
	}
	customTypeTagMap.Set(name, fn)
}
