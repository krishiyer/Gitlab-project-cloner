package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/yaml.v2"
)

type GroupResponse struct {
	projects []Project `json:"projects,omitempty"`
}

type Project struct {
	id              int `json:"id,omitempty"`
	description     string `json:"description,omitempty"`
	ssh_url_to_repo string `json:"ssh_url_to_repo,omitempty"`
}

type Application struct {
	url		string `yaml:"url"`
	token 	string `yaml:"private-token"`
}

func (c *Application) getConf() *Application {

	yamlFile, err := ioutil.ReadFile("application.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {
	var c Application
	c.getConf()

	fmt.Println(c)
	connect()
}

func connect() {

	groupResponse := GroupResponse{}

	url := "https://gitlab.lab.nbttech.com/api/v3/groups/itim";

	// make a sample HTTP GET request
	req, err := http.NewRequest("GET", url, bytes.NewBufferString(""))
	req.Header.Add("PRIVATE-TOKEN", "Rz3t3ULb8bd6Hezxdaut")

	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	res, err := client.Do(req)

	// check for response error
	if err != nil {
		log.Fatal(err)
	}

	// read all response body
	data, _ := ioutil.ReadAll(res.Body)

	reader := bytes.NewReader(data)
	decoder := json.NewDecoder(reader)

	err = decoder.Decode(data)

	if err != nil {
		fmt.Println("whoops:", err)
	}

	var projects interface{} = groupResponse.projects

	fmt.Println("Printing the projects")
	fmt.Println(projects)
	//for _, p := range projects {
	//	fmt.Println(p.ssh_url_to_repo)
	//}

	// close response body
	res.Body.Close()

	// print `data` as a string
	fmt.Printf("%s\n", data)

}
