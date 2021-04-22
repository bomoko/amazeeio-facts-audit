package gatherers

import (
	"encoding/json"
	"fmt"
	"github.com/bomoko/lagoon-facts/utils"
	"log"
)

type applicationGatherer struct {
	GatheredFacts []GatheredFact
}

type composerShowOutput struct {
	Name string `json:"name,omitempty"`
	Versions []string `json:"versions,omitempty"`
	Description string `json:"description,omitempty"`
}

func (p *applicationGatherer) AppliesToEnvironment() bool {

	a := []string{"drupal/core", "laravel/framework"}

	applies := false

	for _, name := range a {
		err, stdOut, _ := utils.Shellout(fmt.Sprintf("composer show -i --format=json %v 2> /dev/null", name))

		var result composerShowOutput
		if stdOut != "" {
			if err = json.Unmarshal([]byte(stdOut), &result); err != nil {
				log.Printf("Application gather cannot be applied: %s", err.Error())
				return false
			}
		}

		if name == result.Name {
			log.Printf("Found %s:%s", name, result.Versions[0])

			p.GatheredFacts = append(p.GatheredFacts, GatheredFact{
				Name:        result.Name,
				Value:       result.Versions[0],
				Source:      "application_via_composer",
				Description: result.Description,
				Category:  "Application",
			})
		}

		applies = true
	}

	return applies
}

func (p *applicationGatherer) GatherFacts() ([]GatheredFact, error) {
	return p.GatheredFacts, nil
}

func init() {
	RegisterGatherer("Application gatherer", &applicationGatherer{})
}