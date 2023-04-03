package maps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"medvil/model"
	"os"
)

func SaveMap(m *model.Map, dir string) {
	os.MkdirAll(dir, os.ModePerm)
	content, err := json.MarshalIndent(map[string]interface{}{
		"SX": m.SX,
		"SY": m.SY,
	}, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(dir+"/meta.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}

	content, err = json.MarshalIndent(m.Fields, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(dir+"/fields.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}

	content, err = json.MarshalIndent(m.Countries, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(dir+"/society.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
