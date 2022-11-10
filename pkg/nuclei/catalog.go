package nuclei

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"strings"

	"github.com/pkg/errors"
	"github.com/projectdiscovery/gologger"
)

//go:embed embed_fs
var templatesFS embed.FS

func NewCatalog() *Catalog {
	return &Catalog{}
}

type Catalog struct {
}

// 循环遍历文件夹中的文件内容
func templateFsList(fileDir string, templateMap map[string]string) {
	dirEntries, _ := templatesFS.ReadDir(fileDir)
	for _, de := range dirEntries {
		if de.IsDir() {
			templateFsList(fileDir+"/"+de.Name(), templateMap)
		} else {
			if strings.Contains(de.Name(), ".yaml") {
				templateMap[fileDir+"/"+de.Name()] = de.Name()
			}
		}
	}
}

// GetTemplatePath parses the specified input template path and returns a compiled
// list of finished absolute paths to the templates evaluating any glob patterns
// or folders provided as in.
func (c *Catalog) GetTemplatePath(target string) ([]string, error) {
	matches, err := fs.Glob(templatesFS, target)
	templateMap := make(map[string]string)
	for i := 0; i < len(matches); i++ {
		mt := matches[i]
		if strings.Contains(mt, ".yaml") {
			templateMap[mt] = mt
		} else {
			templateFsList(mt, templateMap)
		}
	}
	var newMatches []string
	if len(templateMap) > 0 {
		for key, _ := range templateMap {
			newMatches = append(newMatches, key)
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, "could not find glob matches")
	}

	if len(matches) == 0 {
		return nil, errors.Errorf("no templates found for path")
	}

	return newMatches, nil
}

func (c *Catalog) GetTemplatesPath(definitions []string) []string {
	// keeps track of processed dirs and files
	processed := make(map[string]bool)
	var allTemplates []string

	for _, t := range definitions {
		if strings.HasPrefix(t, "http") && (strings.HasSuffix(t, ".yaml") || strings.HasSuffix(t, ".yml")) {
			if _, ok := processed[t]; !ok {
				processed[t] = true
				allTemplates = append(allTemplates, t)
			}
		} else {
			paths, err := c.GetTemplatePath(t)
			if err != nil {
				gologger.Error().Msgf("Could not find template '%s': %s\n", t, err)
			}
			for _, path := range paths {
				if _, ok := processed[path]; !ok {
					processed[path] = true
					allTemplates = append(allTemplates, path)
				}
			}
		}
	}
	return allTemplates
}

// ResolvePath resolves the path to an absolute one in various ways.
//
// It checks if the filename is an absolute path, looks in the current directory
// or checking the nuclei templates directory. If a second path is given,
// it also tries to find paths relative to that second path.
func (c *Catalog) ResolvePath(templateName, second string) (string, error) {
	return templateName, nil
}

// OpenFile opens a file and returns an io.ReadCloser to the file.
// It is used to read template and payload files based on catalog responses.
func (c *Catalog) OpenFile(filename string) (io.ReadCloser, error) {
	if data, err := templatesFS.ReadFile(filename); err != nil {
		return nil, err
	} else {
		return io.NopCloser(bytes.NewReader(data)), nil
	}
}
