package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/untrustedmodders/go-plugify"
	"gopkg.in/yaml.v3"
)

func NewConfig() *Config {
	return &Config{}
}

func LoadConfig[T ConfigsType](folder, fileName string, defaultData T) (*T, error) {
	configPath := filepath.Join(plugify.ConfigsDir, folder, fileName)

	//fmt.Printf("Loading config from %s\n", configPath)
	err := createConfigFile(configPath, defaultData)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file %s: %w", configPath, err)
	}

	defer file.Close()

	var config T
	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	//fmt.Printf("New Config: %v\n", config)

	return &config, nil
}

func createConfigFile(filePath string, defaultData any) error {
	if _, err := os.Stat(filePath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	dir := filepath.Dir(filePath)
	if dir != "." && dir != string(filepath.Separator) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	var node yaml.Node
	err := node.Encode(defaultData)
	if err != nil {
		return err
	}

	addCommentsFromTags(&node, reflect.TypeOf(defaultData))

	data, err := yaml.Marshal(&node)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0666)
}

func addCommentsFromTags(node *yaml.Node, t reflect.Type) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if node.Kind != yaml.MappingNode {
		return
	}

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valueNode := node.Content[i+1]

		for j := range t.NumField() {
			field := t.Field(j)
			yamlTag := field.Tag.Get("yaml")

			if yamlTag != keyNode.Value {
				continue
			}

			if comment, ok := field.Tag.Lookup("head_comment"); ok {
				prefix := "\n# "

				if i == 0 {
					prefix = "# "
				}

				node.Content[i].HeadComment = prefix + comment
			}

			node.Content[i].LineComment = field.Tag.Get("line_comment")
			node.Content[i].FootComment = field.Tag.Get("foot_comment")

			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			if fieldType.Kind() == reflect.Struct && valueNode.Kind == yaml.MappingNode {
				addCommentsFromTags(valueNode, fieldType)
			}

			break
		}
	}
}
