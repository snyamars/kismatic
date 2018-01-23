package plan

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/apprenda/kismatic/pkg/install"
	yaml "gopkg.in/yaml.v2"
)

const templateFilename = "cluster.yaml.tpl"

// The ProviderTemplatePlanner looks for plan file templates in the KET
// infrastructure provider directory.
type ProviderTemplatePlanner struct {
	ProvidersDir string
}

// GetPlanTemplate returns a template plan file for the given provider
func (ptp ProviderTemplatePlanner) GetPlanTemplate(provider string) (*install.Plan, error) {
	tpl := filepath.Join(ptp.ProvidersDir, provider, templateFilename)
	b, err := ioutil.ReadFile(tpl)
	if err != nil {
		return nil, err
	}
	var p install.Plan
	if err := yaml.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("error reading plan from template file %q: %v", tpl, err)
	}
	return &p, nil
}
