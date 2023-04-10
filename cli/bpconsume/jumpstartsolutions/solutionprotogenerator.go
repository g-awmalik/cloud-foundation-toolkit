package jumpstartsolutions

import (
	"fmt"
	"strings"

	gen_protos "github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpconsume/jumpstartsolutions/gen-protos"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
)

// generateSolutionProto creates the Solution object from the BlueprintMetadata
// object.
func generateSolutionProto(bpObj *bpmetadata.BlueprintMetadata) (*gen_protos.Solution, error) {
	solution := &gen_protos.Solution{}

	addGitSource(solution, bpObj)
	addDeploymentTimeEstimate(solution, bpObj)
	addCostEstimate(solution, bpObj)

	solution.DeployData = &gen_protos.DeployData{}
	err := addRoles(solution, bpObj)
	if err != nil {
		return nil, err
	}

	addApis(solution, bpObj)
	addVariables(solution, bpObj)
	addOutputs(solution, bpObj)

	return solution, nil
}

// addGitSource adds the solution's git source to the solution object from the
// BlueprintMetadata object.
func addGitSource(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if solution.GitSource == nil {
		solution.GitSource = &gen_protos.GitSource{}
	}
	if bpObj.Spec.Source != nil {
		solution.GitSource.Repo = strings.TrimSuffix(bpObj.Spec.Source.Repo, ".git")
	}
}

// addDeploymentTimeEstimate adds the deployment time for the solution to the
// solution object from the BlueprintMetadata object.
func addDeploymentTimeEstimate(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if bpObj.Spec.DeploymentDuration.ConfigurationSecs > 0 && bpObj.Spec.DeploymentDuration.DeploymentSecs > 0 {
		solution.DeploymentEstimate = &gen_protos.DeploymentEstimate{
			ConfigurationMinutes: int32(bpObj.Spec.DeploymentDuration.ConfigurationSecs / 60),
			DeploymentMinutes:    int32(bpObj.Spec.DeploymentDuration.DeploymentSecs / 60),
		}
	}
}

// addCostEstimate adds the cost estimate for the solution to the solution
// object from the BlueprintMetadata object.
func addCostEstimate(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if bpObj.Spec.CostEstimate.URL != "" {
		solution.CostEstimateLink = bpObj.Spec.CostEstimate.URL
	}
}

// addRoles adds the roles required by the service account deploying the
// solution to the solution object from the BlueprintMetdata object.
func addRoles(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) error {
	if len(bpObj.Spec.Roles) == 0 {
		return nil
	}
	projectRoleCount := 0
	for _, bpRoles := range bpObj.Spec.Roles {
		if bpRoles.Level == "Project" {
			projectRoleCount += 1
		}
	}
	if projectRoleCount > 1 {
		return fmt.Errorf("more than one set of project level roles present in OSS solution metadata")
	}
	for _, bpRoles := range bpObj.Spec.Roles {
		if bpRoles.Level == "Project" {
			solution.DeployData.Roles = make([]string, len(bpRoles.Roles))
			copy(solution.DeployData.Roles, bpRoles.Roles)
		}
	}
	return nil
}

// addApis adds the APIs required for deploying the solution to the solution
// object from the BlueprintMetadata object.
func addApis(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if len(bpObj.Spec.Services) == 0 {
		return
	}
	solution.DeployData.Apis = make([]string, len(bpObj.Spec.Services))
	copy(solution.DeployData.Apis, bpObj.Spec.Services)
}

// addVariables adds the terraform input variables to the solution object from
// the BlueprintMetadata object.
func addVariables(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if len(bpObj.Spec.BlueprintInterface.Variables) == 0 {
		return
	}
	solution.DeployData.InputSections = []*gen_protos.Section{}
	for _, variable := range bpObj.Spec.BlueprintInterface.Variables {
		property := &gen_protos.Property{
			Name:       variable.Name,
			IsRequired: variable.Required,
		}
		switch variable.VarType {
		case "string":
			property.Type = gen_protos.Property_STRING
			if variable.Default != nil {
				property.DefaultValue = fmt.Sprintf("%v", variable.Default)
			}
		case "bool":
			property.Type = gen_protos.Property_BOOLEAN
			if variable.Default != nil {
				property.DefaultValue = fmt.Sprintf("%v", variable.Default)
			}
		case "list":
			property.Type = gen_protos.Property_ARRAY
		case "number":
			// Note: tf metadata uses "number" type for both "integer" and "number" type.
			// Hence, this might require manual update of textproto file.
			property.Type = gen_protos.Property_NUMBER
			if variable.Default != nil {
				property.DefaultValue = fmt.Sprintf("%v", variable.Default)
			}
		}
		solution.DeployData.InputSections = append(solution.DeployData.InputSections, &gen_protos.Section{
			Properties: []*gen_protos.Property{property},
		})
	}
}

// addOutputs adds the terraform outputs to the solution object from the
// BlueprintMetadata object.
func addOutputs(solution *gen_protos.Solution, bpObj *bpmetadata.BlueprintMetadata) {
	if len(bpObj.Spec.Outputs) == 0 {
		return
	}
	solution.DeployData.Links = []*gen_protos.DeploymentLink{}
	for _, link := range bpObj.Spec.Outputs {
		solution.DeployData.Links = append(solution.DeployData.Links, &gen_protos.DeploymentLink{
			OutputName: link.Name,
		})
	}
}
