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
	soy_diagram_description_msg = "{msg desc=\"Step $COUNT of $SOLUTION_NAME diagram description\"}\n$SOLUTION_DIAGRAM_DESCRIPTION\n{/msg}\n"
	soy_line_separator          = "{\\n}"
)

func generate_solution_id(solution_name string) string {
	return strings.ReplaceAll(strings.ToLower(solution_name), " ", "_")
}

func create_diagram_description(steps []string, solution_name string) string {
	var buffer bytes.Buffer
	for i, step := range steps {
		if i > 0 {
			buffer.WriteString(soy_line_separator)
		}
		buffer.WriteString(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(soy_diagram_description_msg, "$SOLUTION_DIAGRAM_DESCRIPTION", step), "$SOLUTION_NAME", solution_name), "$COUNT", strconv.Itoa(i+1)))
	}
	return buffer.String()
}

func generate_soy(bpObj *bpmetadata.BlueprintMetadata) error {
	solution_name := bpObj.Spec.BlueprintInfo.Title
	solution_id := generate_solution_id(solution_name)
	solution_title := bpObj.Spec.BlueprintInfo.Title
	solution_summary := bpObj.Spec.BlueprintInfo.Description.Tagline
	solution_description := bpObj.Spec.BlueprintInfo.Description.Detailed
	solution_diagram_steps := bpObj.Spec.BlueprintContent.Architecture.Description
	solution_diagram_description := create_diagram_description(solution_diagram_steps, solution_name)

	replacer := strings.NewReplacer("$SOLUTION_ID", solution_id, "$SOLUTION_NAME", solution_name, "$SOLUTION_TITLE", solution_title, "$SOLUTION_SUMMARY", solution_summary, "$SOLUTION_DESCRIPTION", solution_description, "$DIAGRAM_DESCRIPTION", solution_diagram_description)

	input, err := ioutil.ReadFile("soy_template.soy")
	if err != nil {
		return err
	}
	output := replacer.Replace(string(input))
	fileName := solution_id + ".soy"
	curren_dir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(curren_dir, fileName), []byte(output), 0644); err != nil {
		return err
	}
	return nil
}
