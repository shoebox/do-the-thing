package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	General        General                  `yaml:"general"`
	SigninConfig   map[string]SigninConfig  `yaml:"signinConfig"`
	ProductFlavors map[string]ProductFlavor `yaml:"productFlavors"`
}

type General struct {
	Project         string `yaml:"project"`
	Workspace       string `yaml:"workspace"`
	OutputDirectory string `yaml:"outputDirectory"`
	FileName        string `yaml:"name"`
}

type SigninConfig struct {
	Path string `yaml:"path"`
}

type ProductFlavor struct {
	BuildConfiguration string `yaml:"buildConfiguration"`
	Scheme             string `yaml:"scheme"`
	SigninConfig       string `yaml:"signinConfig"`
}

func Parse() (Config, error) {
	var res Config

	// config
	y, err := ioutil.ReadFile("project.yml")
	if err != nil {
		return res, err
	}

	err = yaml.Unmarshal(y, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}
