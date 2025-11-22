package resume

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)


type Editor struct {
	scriptPath string
	templatePath string
	projectsJSON string
	outputDir string
	pythonPath string
}

func NewEditor(scriptPath, templatePath, projectsJSON, outputDir string) *Editor {
	//Try to find python in venv first
	venvPython := filepath.Join(filepath.Dir(scriptPath), "venv", "bin", "python")
	pythonPath := "python3"

	if _, err := os.Stat(venvPython); err == nil {
		pythonPath = venvPython
	}


	return &Editor{
		scriptPath: scriptPath,
		templatePath: templatePath,
		projectsJSON: projectsJSON,
		outputDir: outputDir,
		pythonPath: pythonPath,
	}
}

// TailorResume creates a tailored resume with only specified projects
func (e *Editor) TailorResume(projects []string, outputName string) (string, error) {
	//Ensure output directory exists
	if err := os.MkdirAll(e.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output dir: %w", err)
	}

	//Build output path
	outputPath := filepath.Join(e.outputDir, outputName)
	
	//Join projects with comma
	projectsArgs := strings.Join(projects, ",")


	//Run Python script
	cmd := exec.Command(e.pythonPath, e.scriptPath,
		"--template", e.templatePath,
		"--projects-json", e.projectsJSON,
		"--keep", projectsArgs,
		"--output", outputPath
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("resume editor failed: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("Resume editor outputL\n%s\n", string(output))


	return outputPath, nil
}

//SanitizeFilename removes/replaces characters unsuitable for filenames
func SanitizeFilename(s string) string {
	replacer := strings.NewReplacer(
		" ", "_",
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "-",
		"(", "",
		")", "",
		",", "",
	)

	return replacer.Replace(s)
}
