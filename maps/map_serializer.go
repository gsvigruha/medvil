package maps

import (
	"reflect"
	"bytes"
)

func Serialize(o interface{}) {
	var writer bytes.Buffer
	SerializeObject(o, writer)

}

func SerializeObject(o interface{}, writer bytes.Buffer) {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		fmt.Println(strings.Repeat("\t", depth+1), "Contained type:", t.Key())
		SerializeObject(t.Elem(), depth+1)
	case reflect.Map:
	case reflect.Ptr:
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)

			fmt.Println(strings.Repeat("\t", depth+1), "Field", i+1, "name is", f.Name, "type is", f.Type.Name(), "and kind is", f.Type.Kind())

		}
	}
}