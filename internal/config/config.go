package config

import (
	"log"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

type AppRunConfig struct {
	DryRun   bool
	FilePath string
}

type MovementConfiguration struct {
	DryRun    bool          `json:"dry_run"`
	Movements []MovementRun `json:"movements"`
}

type MovementRun struct {
	Source              string               `json:"source"`
	Destination         string               `json:"destination"`
	TransformationRules []TransformationRule `json:"transformation_rules"`
	FilterRules         []FilterRule         `json:"filter_rules"`
}

type TransformationRule interface {
	Apply(source string) string
}

// TransformationRuleExtension es una regla que transforma el nombre del archivo por extension.
// Ejemplo:
//
//	{
//	  "type": "extension",
//	  "extensions": [
//	    { "from": "jpg", "to": "png" },
//	    { "from": "jpeg", "to": "png" }
//	  ]
//	}
//
// En este caso, si el archivo es "image.jpg", se transforma a "image.png".
// Si el archivo es "image.jpeg", se transforma a "image.png".
// Si el archivo no tiene una extension que coincida con las extensiones de la lista, se devuelve el archivo original.
type TransformationRuleExtension struct {
	Type       string         `json:"type"`
	Extensions []ExtensionDuo `json:"extensions"`
}

// TransformationRulePathChange es una regla que transforma el path del archivo.
// Ejemplo:
//
//	{
//	  "type": "path_change",
//	  "from": "old_path",
//	  "to": "new_path"
//	}
//
// En este caso, si el archivo es "old_path/image.jpg", se transforma a "new_path/image.png".
// Si el archivo no contiene el from, se devuelve el archivo original.
type TransformationRulePathChange struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
}

func (r *TransformationRulePathChange) Apply(source string) string {
	// Solo si contiene el from, se reemplaza por el to.
	if !strings.Contains(source, r.From) {
		return source
	}
	return strings.ReplaceAll(source, r.From, r.To)
}

type ExtensionDuo struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Primero buscamos la extension del archivo y luego comprobamos si en la lista tenemos una extension que coincida con la del archivo.
func (r *TransformationRuleExtension) Apply(source string) string {
	sourceExtension := strings.ToLower(filepath.Ext(source))
	for _, extension := range r.Extensions {
		if sourceExtension == extension.From {
			return strings.Replace(source, sourceExtension, extension.To, 1)
		}
	}
	return source
}

type FilterRule interface {
	Allowed(fullSourcePath string) bool
}

// FilterRuleRegex es una regla que filtra los archivos por regex.
// Ejemplo:
//
//	{
//	  "type": "regex",
//	  "pattern": ".*\\.jpg$"
//	}
type FilterRuleRegex struct {
	Type    string `json:"type"`
	Pattern string `json:"pattern"`
}

func (r *FilterRuleRegex) Allowed(fullSourcePath string) bool {
	re, err := regexp.Compile(r.Pattern)
	if err != nil {
		log.Fatalf("Error compiling regex: %v", err)
		return false
	}
	return re.MatchString(fullSourcePath)
}

// FilterRuleExtension es una regla que filtra los archivos por extension.
// Ejemplo:
//
//	{
//	  "type": "extension",
//	  "extensions": ["jpg", "png", "gif"]
//	}
type FilterRuleExtension struct {
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

func (r *FilterRuleExtension) Allowed(fullSourcePath string) bool {
	extension := strings.ToLower(filepath.Ext(fullSourcePath))
	return slices.Contains(r.Extensions, extension)
}
