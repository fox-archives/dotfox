package main

import (
	"encoding/json"
	"io/ioutil"
)

type PackageJson struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Type        string `json:"typye"`
	Description string `json:"description"`
	Repository  struct {
		type: string
		url string
	} `json:"repository"`
	Bugs struct {
		email string
		url string
	} `json:"bugs"`
	Contributors struct {
		name string
		email string
		url string
	} `json:"contributors"`
	Keywords []string          `json:"keywords"`
	Homepage string            `json:"homepage"`
	License  string            `json:"license"`
	Files    []string          `json:"files"`
	Main     string            `json:"main"`
	Scripts  map[string]string `json:"scripts"`
	Os       []string          `json:"os"`
	Cpu      []string          `json:"cpu"`
	Private  bool              `json:"private"`
}

func main() {
	var packageJson PackageJson
	packageJsonContent, err := ioutil.ReadFile("package.json")
	if err != nil {
		panic(packageJsonContent)
	}

	json.Unmarshal(packageJsonContent, &packageJson)

	packageJson.
		fmt.Println(packageJson.Name)
}
