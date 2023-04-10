package jumpstartsolutions

import (
	"bytes"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/cli/bpmetadata"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	soyDiagramDescriptionMsg = "{msg desc=\"Step $COUNT of $SOLUTION_NAME diagram description\"}\n$SOLUTION_DIAGRAM_DESCRIPTION\n{/msg}\n"
	soyLineSeparator          = "{\\n}"
)

func generateSolutionId(solutionName string) string {
	return strings.ReplaceAll(strings.ToLower(solutionName), " ", "_")
}

func createDiagramDescription(steps []string, solutionName string) string {
	var buffer bytes.Buffer
	for iteration, step := range steps {
		if iteration > 0 {
			buffer.WriteString(soyLineSeparator)
		}
		buffer.WriteString(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(soyDiagramDescriptionMsg, "$SOLUTION_DIAGRAM_DESCRIPTION", step), "$SOLUTION_NAME", solutionName), "$COUNT", strconv.Itoa(iteration+1)))
	}
	return buffer.String()
}

func generateSoy(bpObj *bpmetadata.BlueprintMetadata) error {
	solutionName := bpObj.Spec.BlueprintInfo.Title
	solutionId := generateSolutionId(solutionName)
	solutionTitle := bpObj.Spec.BlueprintInfo.Title
	solutionSummary := bpObj.Spec.BlueprintInfo.Description.Tagline
	solutionDescription := bpObj.Spec.BlueprintInfo.Description.Detailed
	solutionDiagramSteps := bpObj.Spec.BlueprintContent.Architecture.Description
	solutionDiagramDescription := createDiagramDescription(solutionDiagramSteps, solutionName)

	replacer := strings.NewReplacer("$SOLUTION_ID", solutionId, "$SOLUTION_NAME", solutionName, "$SOLUTION_TITLE", solutionTitle, "$SOLUTION_SUMMARY", solutionSummary, "$SOLUTION_DESCRIPTION", solutionDescription, "$DIAGRAM_DESCRIPTION", solutionDiagramDescription)

	input, err := ioutil.ReadFile("soy_template.soy")
	if err != nil {
		return err
	}
	output := replacer.Replace(string(input))
	fileName := solutionId + ".soy"
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(currentDir, fileName), []byte(output), 0644); err != nil {
		return err
	}
	return nil
}