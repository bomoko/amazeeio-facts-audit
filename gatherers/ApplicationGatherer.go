package gatherers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/uselagoon/lagoon-facts-app/utils"
)

type applicationGatherer struct {
	GatheredFacts []GatheredFact
}

func (p *applicationGatherer) GetGathererCmdType() string {
	return GATHERER_TYPE_STATIC
}

type composerShowOutput struct {
	Name        string   `json:"name,omitempty"`
	Versions    []string `json:"versions,omitempty"`
	Description string   `json:"description,omitempty"`
}

func (p *applicationGatherer) AppliesToEnvironment() bool {
	applies := false

	a := []string{"drupal/core", "laravel/framework"}
	// Check composer show
	for _, name := range a {
		err, _, stdOut := utils.Shellout(fmt.Sprintf("composer show -i --format=json %v 2> /dev/null", name))

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
				Category:    Application,
				KeyFact: true,
			})
		}

		applies = true
	}

	// Check for wordpress
	wpErr, wpStdOut, _ := utils.Shellout("wp core version 2> /dev/null")
	if wpStdOut != "" {
		if wpErr != nil {
			log.Println(wpErr)
		}

		version := strings.TrimSuffix(wpStdOut, "\n")

		if version != "" {
			log.Printf("Found wordpress version %s", version)
			p.GatheredFacts = append(p.GatheredFacts, GatheredFact{
				Name:        "wordpress",
				Value:       version,
				Source:      "application_via_cli",
				Description: "The current version of Wordpress that is running",
				Category:    Framework,
				KeyFact: true,
			})

			applies = true
		}
	}

	return applies
}

func (p *applicationGatherer) GatherFacts() ([]GatheredFact, error) {
	return p.GatheredFacts, nil
}

func init() {
	RegisterGatherer("Application gatherer", &applicationGatherer{})
}
