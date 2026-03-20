package config

import (
	"encoding/json"
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
	DryRun                 bool          `json:"dry_run" default:"false"`
	DeleteEmptyDirectories bool          `json:"delete_empty_directories" default:"false"`
	Movements              []MovementRun `json:"movements"`
}

type MovementConfigurationAlias struct {
	DryRun                 bool              `json:"dry_run"`
	DeleteEmptyDirectories bool              `json:"delete_empty_directories" default:"false"`
	Movements              []json.RawMessage `json:"movements"`
}

type MovementRun struct {
	Source              string               `json:"source"`
	Recursive           bool                 `json:"recursive" default:"false"`
	ChangeKeyMap        []ChangeKey          `json:"change_key_map"`
	TransformationRules []TransformationRule `json:"transformation_rules"`
	FilterRules         []FilterRule         `json:"filter_rules"`
}
type MovementRunAlias struct {
	Source              string            `json:"source"`
	Recursive           bool              `json:"recursive" -:"false"`
	ChangeKeyMap        []ChangeKey       `json:"change_key_map"`
	TransformationRules []json.RawMessage `json:"transformation_rules"`
	FilterRules         []json.RawMessage `json:"filter_rules"`
}

type ChangeKey struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

/*
Las transformaciones:

1. Se aplican en el mismo orden en que aparecen en la configuración.
2. Siempre se aplican sobre el path completo del archivo (no solo sobre el nombre).

Ejemplos:

  - Caso 1:
    Entrada: "old_path/document.txt"
    Reglas:
    a) path_change: "old_path" -> "new_path"
    b) regex: "(.*)\\.txt" -> "$1.md"
    Resultado final: "new_path/document.md"

  - Caso 2:
    Entrada: "image.jpg"
    Reglas:
    a) extension: ".jpg" -> ".png"
    b) regex: "(.*)\\.png" -> "$1.gif"
    Resultado final: "image.gif"

  - Caso 3:
    Entrada: "old_path/image.jpg"
    Reglas:
    a) extension: ".jpg" -> ".png"
    b) path_change: "old_path" -> "new_path"
    Resultado final: "new_path/image.png"
*/
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

// El pattern se aplica sobre el path completo del archivo, no solo sobre el nombre.
// Ejemplo:
//
//	{
//	  "type": "regex",
//	  "pattern": "(.*)\\.txt",
//	  "replacement": "$1.md"
//	}
//
// En este caso, si el archivo es "old_path/document.txt", se transforma a "old_path/document.md".
// Ejemplo de cambio de path completo:
//	{
//	  "type": "regex",
//	  "pattern": "old_path/(.*)\\.txt",
//	  "replacement": "new_path/$1.md"
//	}
// En este caso, si el archivo es "old_path/document.txt", se transforma a "new_path/document.md".

type TransformationRuleRegex struct {
	Type        string `json:"type"`
	Pattern     string `json:"pattern"`
	Replacement string `json:"replacement"`
}

func (r *TransformationRuleRegex) Apply(source string) string {
	re, err := regexp.Compile(r.Pattern)
	if err != nil {
		log.Fatalf("Error compiling regex: %v", err)
		return source
	}
	return re.ReplaceAllString(source, r.Replacement)
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

func (m *MovementRun) UnmarshalJSON(data []byte) error {
	var aux MovementRunAlias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	m.Source = aux.Source
	m.ChangeKeyMap = aux.ChangeKeyMap
	m.Recursive = aux.Recursive

	for _, tr := range aux.TransformationRules {
		rule, err := unmarshallTransformationRule(tr)
		if err != nil {
			return err
		}
		// Solo si no es null, se agrega a la lista de reglas de transformación.
		if rule != nil {
			m.TransformationRules = append(m.TransformationRules, rule)
		}
	}

	for _, fr := range aux.FilterRules {
		rule, err := unmarshallFilterRule(fr)
		if err != nil {
			return err
		}
		if rule != nil {
			m.FilterRules = append(m.FilterRules, rule)
		}
	}
	return nil
}

func unmarshallTransformationRule(data []byte) (TransformationRule, error) {
	var head struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, err
	}
	if unmarshaller, ok := transformationRuleUnmarshallers[head.Type]; ok {
		return unmarshaller(data)
	}
	return nil, nil
}

// Aqui tenemos un mapa con string y la funcion de unmarshall
var transformationRuleUnmarshallers = map[string]func(data []byte) (TransformationRule, error){
	"extension":   extensionTransformation,
	"path_change": pathChangeTransformation,
	"regex":       regexTransformation,
}

func extensionTransformation(data []byte) (TransformationRule, error) {
	var rule TransformationRuleExtension
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

func pathChangeTransformation(data []byte) (TransformationRule, error) {
	var rule TransformationRulePathChange
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

func regexTransformation(data []byte) (TransformationRule, error) {
	var rule TransformationRuleRegex
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

var filterRuleUnmarshallers = map[string]func(data []byte) (FilterRule, error){
	"regex":     regexFilterTransformation,
	"extension": extensionFilterTransformation,
}

func unmarshallFilterRule(data []byte) (FilterRule, error) {
	var head struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &head); err != nil {
		return nil, err
	}
	if unmarshaller, ok := filterRuleUnmarshallers[head.Type]; ok {
		return unmarshaller(data)
	}
	return nil, nil
}

func regexFilterTransformation(data []byte) (FilterRule, error) {
	var rule FilterRuleRegex
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

func extensionFilterTransformation(data []byte) (FilterRule, error) {
	var rule FilterRuleExtension
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// Aqui buscamos partes de la url objetivo despues de la transformacion buscamos el texto entre {}
// y lo cambiamos por lo que nos viene en el mapa de inserciones. Por ejemplo, si tenemos una url "new_path/{filename}.jpg" y un mapa de inserciones {"filename": "image"}, el resultado seria "new_path/image.jpg".
func (c *MovementRun) applyInsertions(source string, insertionMap map[string]string) string {
	result := source
	for key, value := range insertionMap {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	//Applicar internas si las hay
	for _, dupla := range c.ChangeKeyMap {
		placeholder := "{" + dupla.Key + "}"
		result = strings.ReplaceAll(result, placeholder, dupla.Value)
	}
	return result
}

func (mr *MovementRun) ApplyTransformations(source string) string {
	result := source
	for _, tr := range mr.TransformationRules {
		result = tr.Apply(result)
	}
	return result
}

func (mr *MovementRun) AllowedByFilters(fullSourcePath string) bool {
	for _, fr := range mr.FilterRules {
		if !fr.Allowed(fullSourcePath) {
			return false
		}
	}
	return true
}

func (mr *MovementRun) Process(source string, insertionMap map[string]string) string {
	transformed := mr.ApplyTransformations(source)
	withInsertions := mr.applyInsertions(transformed, insertionMap)
	return withInsertions
}

func (mc *MovementConfiguration) UnmarshalJSON(data []byte) error {
	var aux MovementConfigurationAlias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	mc.DryRun = aux.DryRun
	mc.DeleteEmptyDirectories = aux.DeleteEmptyDirectories
	for _, movement := range aux.Movements {
		var movementRun MovementRun
		if err := movementRun.UnmarshalJSON(movement); err != nil {
			return err
		}
		mc.Movements = append(mc.Movements, movementRun)
	}
	return nil
}

func (m *MovementRun) Clone() MovementRun {
	return MovementRun{
		Source:              m.Source,
		Recursive:           m.Recursive,
		ChangeKeyMap:        append([]ChangeKey(nil), m.ChangeKeyMap...),
		TransformationRules: append([]TransformationRule(nil), m.TransformationRules...),
		FilterRules:         append([]FilterRule(nil), m.FilterRules...),
	}
}
