package gogen

import (
	"fmt"
	goformat "go/format"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	goctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
)

const goModeIdentifier = "go.mod"

func getParentPackage(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	var rootPath string
	var tempPath = absDir
	var hasGoMod = false
	for {
		tempPath = filepath.Dir(tempPath)
		if goctlutil.FileExists(filepath.Join(tempPath, goModeIdentifier)) {
			tempPath = filepath.Dir(tempPath)
			rootPath = absDir[len(tempPath)+1:]
			hasGoMod = true
			break
		}
		if tempPath == string(filepath.Separator) {
			break
		}
	}
	if !hasGoMod {
		gopath := os.Getenv("GOPATH")
		parent := path.Join(gopath, "src")
		pos := strings.Index(absDir, parent)
		if pos < 0 {
			message := fmt.Sprintf("%s not in gomod project path, or not in GOPATH of %s directory", absDir, gopath)
			println(message)
			tempPath = filepath.Dir(absDir)
			rootPath = absDir[len(tempPath)+1:]
		} else {
			rootPath = absDir[len(parent)+1:]
		}
	}
	return rootPath, nil
}

func writeIndent(writer io.Writer, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Fprint(writer, "\t")
	}
}

func writeProperty(writer io.Writer, name, tp, tag, comment string, indent int) error {
	writeIndent(writer, indent)
	var err error
	if len(comment) > 0 {
		comment = strings.TrimPrefix(comment, "//")
		comment = "//" + comment
		_, err = fmt.Fprintf(writer, "%s %s %s %s\n", strings.Title(name), tp, tag, comment)
	} else {
		_, err = fmt.Fprintf(writer, "%s %s %s\n", strings.Title(name), tp, tag)
	}
	return err
}

func getAuths(api *spec.ApiSpec) []string {
	var authNames = collection.NewSet()
	for _, g := range api.Service.Groups {
		if value, ok := util.GetAnnotationValue(g.Annotations, "server", "jwt"); ok {
			authNames.Add(value)
		}
		if value, ok := util.GetAnnotationValue(g.Annotations, "server", "signature"); ok {
			authNames.Add(value)
		}
	}
	return authNames.KeysStr()
}

func formatCode(code string) string {
	ret, err := goformat.Source([]byte(code))
	if err != nil {
		return code
	}
	return string(ret)
}
