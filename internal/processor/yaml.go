package processor

import "gopkg.in/yaml.v3"

func ProcessYAML(data []byte) ([]byte, error) {
	var parsedYAML interface{}
	if err := yaml.Unmarshal(data, &parsedYAML); err != nil {
		return nil, err
	}

	processedYAML := recursiveProcessString(parsedYAML)
	return yaml.Marshal(processedYAML)
}
