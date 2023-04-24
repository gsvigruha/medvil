package maps

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

func Serialize(o interface{}, dir string) {
	var writer bytes.Buffer
	var objects map[uintptr]interface{} = make(map[uintptr]interface{})
	CollectObjects(o, objects)
	log.Printf("Objects to save %d", len(objects))
	writer.WriteString("{")
	var first bool = true
	for ptr, obj := range objects {
		if !first {
			writer.WriteString(", ")
		}
		writer.WriteString("\"" + fmt.Sprint(ptr) + "\"")
		writer.WriteString(": ")
		SerializeObject(obj, &writer)
		first = false
	}
	writer.WriteString("}")
	os.MkdirAll(dir, os.ModePerm)
	err := ioutil.WriteFile(dir+"/"+time.Now().Format(time.RFC3339), writer.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func staticType(t reflect.Type) bool {
	return t.Elem().Name() == "Artifact" ||
		t.Elem().Name() == "Material" ||
		t.Elem().Name() == "PlantType" ||
		t.Elem().Name() == "TerrainType" ||
		t.Elem().Name() == "VehicleType" ||
		t.Elem().Name() == "EquipmentType" ||
		t.Elem().Name() == "Manufacture" ||
		t.Elem().Name() == "RoadType" ||
		t.Elem().Name() == "AnimalType" ||
		t.Elem().Name() == "BuildingExtensionType" ||
		t.Elem().Name() == "VehicleConstruction"
}

func SerializeObject(o interface{}, writer *bytes.Buffer) {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		writer.WriteString("[")
		var first bool = true
		for i := 0; i < v.Len(); i++ {
			if !first {
				writer.WriteString(", ")
			}
			SerializeObject(v.Index(i).Interface(), writer)
			first = false
		}
		writer.WriteString("]")
	case reflect.Map:
		writer.WriteString("{")
		var first bool = true
		for _, key := range v.MapKeys() {
			if !first {
				writer.WriteString(", ")
			}
			SerializeObject(key.Interface(), writer)
			writer.WriteString(": ")
			SerializeObject(v.MapIndex(key).Interface(), writer)
			first = false
		}
		writer.WriteString("}")
	case reflect.Ptr:
		if !v.IsNil() {
			if staticType(t) {
				writer.WriteString("\"" + v.Elem().FieldByName("Name").String() + "\"")
			} else {
				writer.WriteString("\"" + fmt.Sprint(v.Pointer()) + "\"")
			}
		} else {
			writer.WriteString("null")
		}
	case reflect.Struct:
		writer.WriteString("{")
		writer.WriteString("\"$pkg\": \"" + t.PkgPath() + "\", ")
		writer.WriteString("\"$type\": \"" + t.Name() + "\"")
		for i := 0; i < t.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Func {
				fmt.Println(t.Field(i).Name)
			}
			if v.Field(i).CanInterface() {
				writer.WriteString(", ")
				writer.WriteString("\"" + t.Field(i).Name + "\": ")
				if v.Field(i).Kind() == reflect.Ptr || v.Field(i).Kind() == reflect.Interface {
					if !v.Field(i).IsNil() {
						SerializeObject(v.Field(i).Interface(), writer)
					} else {
						writer.WriteString("null")
					}
				} else {
					fv := v.Field(i).Interface()
					SerializeObject(fv, writer)
				}
			}
		}
		writer.WriteString("}")
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		writer.WriteString(strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		writer.WriteString(strconv.FormatUint(v.Uint(), 10))
	case reflect.Bool:
		writer.WriteString("\"" + strconv.FormatBool(v.Bool()) + "\"")
	case reflect.String:
		writer.WriteString("\"" + v.String() + "\"")
	case reflect.Float64:
		writer.WriteString(strconv.FormatFloat(v.Float(), 'E', -1, 32))
	default:
		fmt.Println(t.Kind(), o)
	}
}

func CollectObjects(o interface{}, objects map[uintptr]interface{}) {
	t := reflect.TypeOf(o)
	v := reflect.ValueOf(o)
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			CollectObjects(v.Index(i).Interface(), objects)
		}
	case reflect.Map:
		for _, key := range v.MapKeys() {
			CollectObjects(key.Interface(), objects)
			CollectObjects(v.MapIndex(key).Interface(), objects)
		}
	case reflect.Ptr:
		if !v.IsNil() && !staticType(t) {
			if _, ok := objects[v.Pointer()]; !ok {
				objects[v.Pointer()] = v.Elem().Interface()
				CollectObjects(v.Elem().Interface(), objects)
			}
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if v.Field(i).CanInterface() {
				if v.Field(i).Kind() == reflect.Ptr || v.Field(i).Kind() == reflect.Interface {
					if !v.Field(i).IsNil() {
						CollectObjects(v.Field(i).Interface(), objects)
					}
				} else {
					fv := v.Field(i).Interface()
					CollectObjects(fv, objects)
				}
			}
		}
	}
}
