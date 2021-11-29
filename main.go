package main

import (
	"bufio"
	"fmt"
	"github.com/MichaelTJones/walk"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
)

var homeDir = "./"

func init() {
	if dir, err := os.UserHomeDir(); err == nil {
		homeDir = dir
	}
}

func main() {
	workspacePath := homeDir + "/workspace"
	if len(os.Args) > 1 {
		workspacePath = os.Args[1]
	}

	projectChan := make(chan string, runtime.NumCPU() * 2)
	wg := &sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go processProject(projectChan, wg)
	}

	walk.Walk(workspacePath,
		func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() && strings.HasSuffix(walkPath, "/.idea") {
				projectChan <- walkPath
			}

			return nil
		})
	close(projectChan)

	wg.Wait()
}

func processProject(projectChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		projectPath, ok := <- projectChan
		if !ok {
			return
		}

		projectId := extractProjectId(projectPath)
		if projectId != "" {
			projectType := determineProjectType(projectId)
			if projectType != "" {
				fmt.Println(projectType + ": " + path.Dir(projectPath))
			}
		}
	}
}

var projectIdExp = regexp.MustCompile(`component name="ProjectId" id="([^"]*)"`)
func extractProjectId(projectPath string) string {
	xml, err := os.Open(projectPath + "/workspace.xml")
	if err != nil {
		return ""
	}
	defer xml.Close()


	scanner := bufio.NewScanner(xml)

	for scanner.Scan() {
		found := projectIdExp.FindAllStringSubmatch(scanner.Text(), 1)
		if len(found) > 0 {
			return found[0][1]
		}
	}

	return ""
}

var projectTypeExp = regexp.MustCompile(`/JetBrains/([^/0-9.]*)[0-9.]*/`)
func determineProjectType(projectId string) (result string) {
	targetSuffix := "/" + projectId + ".xml"

	filepath.Walk(homeDir + "/.config/JetBrains", func(walkPath string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(walkPath, targetSuffix) {
			found := projectTypeExp.FindAllStringSubmatch(walkPath, 1)
			if len(found) > 0 {
				result = found[0][1]
				return io.EOF
			}
		}
		return nil
	})

	return
}