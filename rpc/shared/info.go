package shared

import (
	"net/http"
	"reflect"
)

// "InfoStruct" generates the specific APIReference structure dynamically.
func InfoStruct(r *http.Request, method string, model interface{}, returns string) APIReference {
	return APIReference{r.Host + r.URL.String(), method,
		ApiRef(model), returns,
		example(r.Host+r.URL.String(), ApiRef(model))}
}

// "example" creates the APIReference example string shown to the devs.
func example(url string, data string) string {
	return "curl --data " + data + " " + url
}

func ApiRef(i interface{}) string {
	data := "{"
	p := apiRef(data, i)
	return p[0:len(p)-2] + "}"
}

func apiRef(data string, i interface{}) string {
	// get the value of the structure
	v := reflect.Indirect(reflect.ValueOf(i))
	// for each field of the structure
	for j := 0; j < v.NumField(); j++ {
		// get field
		f := v.Field(j)
		// name of the field
		n := v.Type().Field(j).Name
		// type of the field
		t := f.Type()
		// switch on the field type
		switch t.Kind() {
		case reflect.Struct:
			// if a nested struct
			data += "\"" + n + "\":" + ApiRef(f) + ", "
		case reflect.Slice:
			// if a slice
			// get the underlying type of the slice
			e := t.Elem()
			if e.Kind() == reflect.Struct { // slice of structs
				// convert to slice object
				s := reflect.Zero(t).Interface()
				// convert to single element of slice object
				sa := reflect.Zero(reflect.TypeOf(s).Elem()).Interface()
				// recursively run the struct builder
				data += "\"" + n + "\":[" + ApiRef(sa) + "], "
				continue
			}
			if e.Kind() == reflect.Slice { // slice of slices
				if e.Elem().Kind() == reflect.Struct {
					// inner elem is struct
					// convert slice to object
					s := reflect.Zero(e).Interface()
					// convert single element of slice to object
					sa := reflect.Zero(reflect.TypeOf(s).Elem()).Interface()
					data += "\"" + n + "\":[[" + ApiRef(sa) + "]], "
					continue
				}
				// inner elm is not a struct
				data += "\"" + n + "\":[[" + e.String() + "1, " + e.String() + "2], [" + e.String() + "1, " + e.String() + "2]], "
				continue
			}
			data += "\"" + n + "\":[" + t.String() + "1, " + t.String() + "2], "
		default:
			data += "\"" + n + "\":\"" + t.String() + "\", "
		}
	}
	return data
}
