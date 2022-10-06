package yqlib

import (
	"fmt"
	"testing"

	"github.com/mikefarah/yq/v4/test"
)

var tomlScenarios = []formatScenario{
	{
		description:  "Parse toml: simple",
		input:        "[owner]\ncat = \"meow\"\ndog = 3\nfrog = true\n",
		expected:     "owner:\n  cat: meow\n  dog: 3\n  frog: true\n",
		scenarioType: "decode",
	},
	// {
	// 	description:  "Parse toml: array",
	// 	input:        "favouriteFoods = [\"pasta\", \"bananas\"]\n",
	// 	expected:     "favouriteFoods:\n  - pasta\n  - bananas\n",
	// 	scenarioType: "decode",
	// },
}

func testTomlScenario(t *testing.T, s formatScenario) {
	switch s.scenarioType {
	case "decode":
		test.AssertResultWithContext(t, s.expected, processFormatScenario(s, NewTomlDecoder(), NewYamlEncoder(2, false, true, true)), s.description)
	default:
		panic(fmt.Sprintf("unhandled scenario type %q", s.scenarioType))
	}
}

func TestTomlScenarios(t *testing.T) {
	for _, tt := range tomlScenarios {
		testTomlScenario(t, tt)
	}
	// genericScenarios := make([]interface{}, len(tomlScenarios))
	// for i, s := range tomlScenarios {
	// 	genericScenarios[i] = s
	// }
	// documentScenarios(t, "usage", "toml", genericScenarios, documentJSONScenario)
}
