package templates

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig/v3"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"

	"easyhelm/internal/config"
)

type Client struct {
	fs *embed.FS
}

func NewClient(fs *embed.FS) *Client {
	return &Client{
		fs: fs,
	}
}

func (c *Client) GenerateTemplates() error {
	useStaticTemplates := false

	tpls, err := getAllDynamicFilenames()
	if err != nil {
		return err
	}

	if len(tpls) == 0 {
		useStaticTemplates = true

		tpls, err = getAllFilenames(c.fs)
		if err != nil {
			return err
		}
	}

	var t *template.Template
	if useStaticTemplates {
		t, err = template.New("").
			Delims("{%", "%}").
			Funcs(funcMap()).
			ParseFS(c.fs,
				tpls...)
		if err != nil {
			return err
		}
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get current working directory: %w", err)
		}

		t, err = template.New("").
			Delims("{%", "%}").
			Funcs(funcMap()).
			ParseFS(os.DirFS(filepath.Join(wd, "input")), tpls...)
		if err != nil {
			return err
		}
	}

	for _, tplpath := range tpls {
		var (
			buf          bytes.Buffer
			templateName = path.Base(tplpath)
		)

		err = t.ExecuteTemplate(&buf, templateName, config.Unmarshal())
		if err != nil {
			return err
		}

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot get current working directory: %w", err)
		}

		err = os.MkdirAll(path.Dir(path.Join(cwd, "output", tplpath)), os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create directory: %w", err)
		}

		f, err := os.Create(path.Join("output", tplpath))
		if err != nil {
			return fmt.Errorf("cannot create file: %w", err)
		}

		zap.S().Infof("Generating template: %s...", path.Join(cwd, "output", tplpath))

		_, err = f.Write(buf.Bytes())
		if err != nil {
			return fmt.Errorf("cannot write file: %w", err)
		}

		_ = f.Close()
	}

	return nil
}

func getAllFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, "chart", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func getAllDynamicFilenames() (files []string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("cannot get current working directory: %w", err)
	}

	if err := fs.WalkDir(os.DirFS(wd), "input/chart", func(path string, d fs.DirEntry, err error) error {
		if d == nil {
			zap.L().Info("The `input/chart` directory is empty. Using static templates.")
			return nil
		}

		if d.IsDir() {
			return nil
		}

		stippedPath := filepath.Join(strings.Split(path, string(os.PathSeparator))[1:]...)
		files = append(files, stippedPath)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

func funcMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	extra := template.FuncMap{
		"toToml":        toTOML,
		"toYaml":        toYAML,
		"fromYaml":      fromYAML,
		"fromYamlArray": fromYAMLArray,
		"toJson":        toJSON,
		"fromJson":      fromJSON,
		"fromJsonArray": fromJSONArray,

		"include":  func(string, interface{}) string { return "not implemented" },
		"tpl":      func(string, interface{}) interface{} { return "not implemented" },
		"required": func(string, interface{}) (interface{}, error) { return "not implemented", nil },
		"lookup": func(string, string, string, string) (map[string]interface{}, error) {
			return map[string]interface{}{}, nil
		},
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

func fromYAML(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

func fromYAMLArray(str string) []interface{} {
	a := []interface{}{}

	if err := yaml.Unmarshal([]byte(str), &a); err != nil {
		a = []interface{}{err.Error()}
	}
	return a
}

func toTOML(v interface{}) string {
	b := bytes.NewBuffer(nil)
	e := toml.NewEncoder(b)
	err := e.Encode(v)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return string(data)
}

func fromJSON(str string) map[string]interface{} {
	m := make(map[string]interface{})

	if err := json.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

func fromJSONArray(str string) []interface{} {
	a := []interface{}{}

	if err := json.Unmarshal([]byte(str), &a); err != nil {
		a = []interface{}{err.Error()}
	}
	return a
}
