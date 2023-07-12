package maps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"medvil/model"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/materials"
	"medvil/model/terrain"
	"medvil/model/vehicles"
	"os"
	"reflect"
	"strconv"
)

func Deserialize(file string) interface{} {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	var jsonData map[string]json.RawMessage
	if err := json.Unmarshal(byteValue, &jsonData); err != nil {
		fmt.Println(err)
	}
	var objects map[string]reflect.Value = make(map[string]reflect.Value)
	log.Printf("Objects to load %d", len(jsonData))
	result := DeserializeObject(jsonData["0"], reflect.TypeOf(model.Map{}), jsonData, objects, nil)
	return result.Addr().Interface()
}

func DeserializeObject(m json.RawMessage, t reflect.Type, jsonData map[string]json.RawMessage, objects map[string]reflect.Value, objKey *string) reflect.Value {
	fmt.Println(t.Kind(), t.Name())
	switch t.Kind() {
	case reflect.Slice:
		var mData []json.RawMessage
		if err := json.Unmarshal(m, &mData); err != nil {
			fmt.Println("Error when loading "+t.Name()+": ", err)
		}
		v := reflect.MakeSlice(t, 0, 0)
		for i := 0; i < len(mData); i++ {
			x := DeserializeObject(mData[i], t.Elem(), jsonData, objects, nil)
			v = reflect.Append(v, x)
		}
		return v
	case reflect.Array:
		var mData []json.RawMessage
		if err := json.Unmarshal(m, &mData); err != nil {
			fmt.Println("Error when loading "+t.Name()+": ", err)
		}
		v := reflect.New(t).Elem()
		for i := 0; i < len(mData); i++ {
			v.Index(i).Set(DeserializeObject(mData[i], t.Elem(), jsonData, objects, nil))
		}
		return v
	case reflect.Map:
		var mData map[string]json.RawMessage
		if err := json.Unmarshal(m, &mData); err != nil {
			fmt.Println("Error when loading "+t.Name()+"["+t.Key().Name()+"]"+t.Elem().Name()+": ", err)
		}
		v := reflect.MakeMap(t)
		for rk, rv := range mData {
			v.SetMapIndex(
				DeserializeObject([]byte(rk), t.Key(), jsonData, objects, nil),
				DeserializeObject(rv, t.Elem(), jsonData, objects, nil),
			)
		}
		return v
	case reflect.Ptr:
		var objKey string
		if err := json.Unmarshal(m, &objKey); err != nil {
			fmt.Println("Error when loading "+t.Name()+" "+string(m)+": ", err)
		}
		if objKey == "" {
			// Hacky way of testing for null
			return reflect.Zero(t)
		} else {
			if StaticType(t) {
				return LoadStaticType(t, objKey)
			} else {
				if _, ok := objects[objKey]; !ok {
					objects[objKey] = DeserializeObject(jsonData[objKey], t.Elem(), jsonData, objects, &objKey).Addr()
				}
				return objects[objKey]
			}
		}
	case reflect.Interface:
		var objKey string
		if err := json.Unmarshal(m, &objKey); err != nil {
			fmt.Println("Error when loading "+t.Name()+" "+string(m)+": ", err)
		}
		if objKey == "" {
			// Hacky way of testing for null
			return reflect.Zero(t)
		} else {
			if StaticType(t) {
				return LoadStaticType(t, objKey)
			} else {
				if _, ok := objects[objKey]; !ok {
					objects[objKey] = DeserializeObject(jsonData[objKey], t, jsonData, objects, &objKey)
				}
				return objects[objKey]
			}
		}
	case reflect.Struct:
		var mData map[string]json.RawMessage
		if err := json.Unmarshal(m, &mData); err != nil {
			fmt.Println("Error when loading "+t.Name()+": ", err)
		}
		v := reflect.New(t)
		if objKey != nil {
			objects[*objKey] = v
		}
		for i := 0; i < t.NumField(); i++ {
			sf := v.Elem().Field(i)
			if sf.Kind() == reflect.Func {
				fmt.Println(t.Field(i).Name)
			}
			if sf.CanInterface() {
				fv := DeserializeObject(mData[t.Field(i).Name], t.Field(i).Type, jsonData, objects, nil)
				sf.Set(fv)
			} else {
				if t.Field(i).Tag.Get("ser") != "false" {
					fmt.Println("Cannot serialize: " + t.Name() + "." + t.Field(i).Name)
				}
			}
		}
		return v.Elem()
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		i, err := strconv.ParseInt(string(m), 10, 64)
		if err != nil {
			fmt.Println("err")
		}
		return reflect.ValueOf(i).Convert(t)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(string(m), 10, 64)
		if err != nil {
			fmt.Println("err")
		}
		return reflect.ValueOf(i).Convert(t)
	case reflect.Bool:
		var s string
		if err := json.Unmarshal(m, &s); err != nil {
			fmt.Println("err")
		}
		i, err := strconv.ParseBool(s)
		if err != nil {
			fmt.Println("err")
		}
		return reflect.ValueOf(i)
	case reflect.String:
		var s string
		if err := json.Unmarshal(m, &s); err != nil {
			fmt.Println("err")
		}
		return reflect.ValueOf(s)
	case reflect.Float64:
		i, err := strconv.ParseFloat(string(m), 64)
		if err != nil {
			fmt.Println("err")
		}
		return reflect.ValueOf(i).Convert(t)
	}
	panic("Invalid type " + t.Name())
}

func LoadStaticType(t reflect.Type, key string) reflect.Value {
	switch t.Elem().Name() {
	case "Artifact":
		return reflect.ValueOf(artifacts.GetArtifact(key))
	case "Material":
		return reflect.ValueOf(materials.GetMaterial(key))
	case "PlantType":
		return reflect.ValueOf(terrain.GetPlantType(key))
	case "TerrainType":
		return reflect.ValueOf(terrain.GetTerrainType(key))
	case "VehicleType":
		return reflect.ValueOf(vehicles.GetVehicleType(key))
	case "EquipmentType":
		return reflect.ValueOf(economy.GetEquipmentType(key))
	case "Manufacture":
		return reflect.ValueOf(economy.GetManufacture(key))
	case "RoadType":
		return reflect.ValueOf(building.GetRoadType(key))
	case "AnimalType":
		return reflect.ValueOf(terrain.GetAnimalType(key))
	case "BuildingExtensionType":
		return reflect.ValueOf(building.GetBuildingExtensionType(key))
	case "VehicleConstruction":
		return reflect.ValueOf(economy.GetVehicleConstruction(key))
	}
	panic("Invalid type " + t.Elem().Name())
}
