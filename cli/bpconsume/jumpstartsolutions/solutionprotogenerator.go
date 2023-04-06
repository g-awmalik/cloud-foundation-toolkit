package jumpstartsolutions

import (
	"fmt"
	"strings"

	gen_protos "github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpconsume/jumpstartsolutions/gen-protos"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
)

func generateSolutionProto(bpObj *bpmetadata.BlueprintMetadata) *gen_protos.Solution {
	solution := &gen_protos.Solution{}

	addGitSource(solution, bpObj)
	addDeploymentTimeEstimate(solution, bpObj)
	addCostEstimate(solution, bpObj)
	addRoles(solution, bpObj)
	addApis(solution, bpObj)
	addVariables(solution, bpObj)
	addOutputs(solution, bpObj)

	return solution
}

func addGitSource(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if solution.GitSource == nil {
		solution.GitSource = &gen_protos.GitSource{}
	}
	if bpObj.Spec.Source != nil {
		solution.GitSource.Repo = strings.TrimSuffix(bpObj.Spec.Source.Repo, ".git")
	}
}

func addDeploymentTimeEstimate(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if bpObj.Spec.DeploymentDuration.ConfigurationSecs > 0 && bpObj.Spec.DeploymentDuration.DeploymentSecs > 0 {
		solution.DeploymentEstimate = &gen_protos.DeploymentEstimate{
			ConfigurationMinutes: int32(bpObj.Spec.DeploymentDuration.ConfigurationSecs / 60),
			DeploymentMinutes:    int32(bpObj.Spec.DeploymentDuration.DeploymentSecs / 60),
		}
	}
}

func addCostEstimate(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if bpObj.Spec.CostEstimate.URL != "" {
		solution.CostEstimateLink = bpObj.Spec.CostEstimate.URL
	}
}

func addRoles(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if len(bpObj.Spec.Roles) == 0 {
		return
	}
	for _, bpRoles := range bpObj.Spec.Roles {
		if bpRoles.Level == "Project" {
			solution.Roles = make([]string, len(bpRoles.Roles))
			copy(solution.Roles, bpRoles.Roles)
		}
	}
}

func addApis(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if len(bpObj.Spec.Services) == 0 {
		return
	}
	solution.Apis = make([]string, len(bpObj.Spec.Services))
	copy(solution.Apis, bpObj.Spec.Services)
}

func addVariables(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	// TODO: add variable type
	if len(bpObj.Spec.BlueprintInterface.Variables) == 0 {
		return
	}
	solution.InputSections = []*gen_protos.Section{}
	for _, variable := range bpObj.Spec.BlueprintInterface.Variables {
		solution.InputSections = append(solution.InputSections, &gen_protos.Section{
			Properties: []*gen_protos.Property{{
				Name:         variable.Name,
				DefaultValue: fmt.Sprintf("%v", variable.Default),
				IsRequired:   variable.Required,
			}},
		})
	}
}

func addOutputs(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if len(bpObj.Spec.Outputs) == 0 {
		return
	}
	solution.Links = []*gen_protos.DeploymentLink{}
	for _, link := range bpObj.Spec.Outputs {
		solution.Links = append(solution.Links, &gen_protos.DeploymentLink{
			OutputName: link.Name,
		})
	}
}
