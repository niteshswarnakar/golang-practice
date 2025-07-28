package mysystemctl

import (
	"fmt"
	"os/exec"
	"strings"
)

type ServiceStatus struct {
	Name       string
	Active     string
	Status     string
	Enabled    string
	FullOutput string
}

func Run() {
	executableFilePath := "/home/iamnitesh/Desktop/work_space/PJNube/runnable-build-2/control-engine"
	output, err := hasAllDependencies(executableFilePath)
	if err != nil {
		fmt.Println("Error checking dependencies:", err)
		return
	}

	fmt.Println("Dependencies for ", executableFilePath)
	fmt.Println(output)
}

// CheckServiceStatus checks the status of a systemd service
func CheckServiceStatus(serviceName string) (*ServiceStatus, error) {
	cmd := exec.Command("systemctl", "status", serviceName)
	output, err := cmd.CombinedOutput()

	status := &ServiceStatus{
		Name:       serviceName,
		FullOutput: string(output),
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Active:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				status.Active = parts[1]
			}
			if len(parts) >= 3 {
				status.Status = strings.Join(parts[2:], " ")
			}
		}
	}

	return status, err
}

func IsServiceActive(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.Output()

	if err != nil {
		return false
	}

	return strings.TrimSpace(string(output)) == "active"
}

func StartService(serviceName string) error {
	cmd := exec.Command("sudo", "systemctl", "start", serviceName)
	return cmd.Run()
}

func hasAllDependencies(exeFilePath string) (bool, []string) {
	status := true
	cmd := exec.Command("ldd", exeFilePath)
	output, err := cmd.Output()
	if err != nil {
		return false, []string{err.Error()}
	}
	var missingDeps []string
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "not found") {
			status = false
			parts := strings.Fields(line)
			if len(parts) > 0 {
				missingDeps = append(missingDeps, parts[0])
			}
		}
	}
	return status, missingDeps
}
