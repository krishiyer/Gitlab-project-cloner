package main

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type YamlConfig struct {
	Application struct {
		Url    string `yaml:"url"`
		Token  string `yaml:"token"`
		GrpUrl string `yaml:"grpUrl"`
		Dest   string `yaml:"destDir"`
	}
}

type GroupResponse struct {
	Projects []Project `json:"projects,omitempty"`
}

type Project struct {
	Id              int    `json:"id,omitempty"`
	Description     string `json:"description,omitempty"`
	Ssh_url_to_repo string `json:"ssh_url_to_repo,omitempty"`
}

// Func main should be as small as possible and do as little as possible by convention
func main() {
	// Generate our config based on the config supplied
	// by the user in the flags

	var yamlFile = "resources/application.yml"
	yamlConfig := readYaml(yamlFile)
	var Projects []Project = getProjects(&yamlConfig)
	destDir := yamlConfig.Application.Dest;

	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		os.Mkdir(destDir, os.ModeDir)
		os.Chdir(destDir)
	}

	for _, p := range Projects {
		cmd := exec.Command("git", "clone", p.Ssh_url_to_repo);
		cmd.Run()
		cmd.Wait()
	}

}

func readYaml(yamlFile string) YamlConfig {
	yamlConfig := YamlConfig{}
	yamlData, err := ioutil.ReadFile(yamlFile)

	err = yaml.Unmarshal([]byte(yamlData), &yamlConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return yamlConfig;

}

func getProjects(config *YamlConfig) []Project {

	url := config.Application.Url + config.Application.GrpUrl

	httpClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("PRIVATE-TOKEN", config.Application.Token)
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	groupResponse := GroupResponse{}
	jsonErr := json.Unmarshal(body, &groupResponse)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	// close response body
	resp.Body.Close()

	return groupResponse.Projects;

}
