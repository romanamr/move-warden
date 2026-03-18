package config

import (
	"encoding/json"
	"os"
	"testing"
)

func readFixture(t *testing.T, fileName string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + fileName)
	if err != nil {
		t.Fatalf("no se pudo leer fixture %s: %v", fileName, err)
	}
	return data
}

func TestFilterRuleExtensionAllowed(t *testing.T) {
	rule := FilterRuleExtension{
		Extensions: []string{".jpg", ".png"},
	}

	if rule.Allowed("file.txt") {
		t.Errorf("Expected file.txt to be not allowed")
	}

	if rule.Allowed("document.md") {
		t.Errorf("Expected document.md to be not allowed")
	}

	if !rule.Allowed("image.jpg") {
		t.Errorf("Expected image.jpg to be allowed")
	}

	if !rule.Allowed("photo.png") {
		t.Errorf("Expected photo.png to be allowed")
	}
}

func TestTransformationRulePathChangeApply(t *testing.T) {
	rule := TransformationRulePathChange{From: "old", To: "new"}
	if got := rule.Apply("/tmp/old/file.jpg"); got != "/tmp/new/file.jpg" {
		t.Fatalf("resultado inesperado: %s", got)
	}
	if got := rule.Apply("/tmp/file.jpg"); got != "/tmp/file.jpg" {
		t.Fatalf("no deberia cambiar cuando no contiene From: %s", got)
	}
}

func TestTransformationRuleExtensionApply(t *testing.T) {
	rule := TransformationRuleExtension{Extensions: []ExtensionDuo{{From: ".jpg", To: ".png"}}}
	if got := rule.Apply("image.jpg"); got != "image.png" {
		t.Fatalf("resultado inesperado: %s", got)
	}
	if got := rule.Apply("image.gif"); got != "image.gif" {
		t.Fatalf("no deberia cambiar extension no mapeada: %s", got)
	}
}

func TestFilterRuleRegexAllowed(t *testing.T) {
	rule := FilterRuleRegex{Pattern: `.*\.jpg$`}
	if !rule.Allowed("image.jpg") {
		t.Fatal("deberia permitir image.jpg")
	}
	if rule.Allowed("image.png") {
		t.Fatal("no deberia permitir image.png")
	}
}

func TestUnmarshallTransformationRule(t *testing.T) {
	data := readFixture(t, "transformation_extension.json")
	rule, err := unmarshallTransformationRule(data)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if _, ok := rule.(*TransformationRuleExtension); !ok {
		t.Fatalf("tipo inesperado para transformation rule")
	}

	dataUnknown := readFixture(t, "transformation_unknown.json")
	unknown, err := unmarshallTransformationRule(dataUnknown)
	if err != nil {
		t.Fatalf("error inesperado para type desconocido: %v", err)
	}
	if unknown != nil {
		t.Fatal("se esperaba nil para type desconocido")
	}
}

func TestUnmarshallFilterRule(t *testing.T) {
	data := readFixture(t, "filter_regex.json")
	rule, err := unmarshallFilterRule(data)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if _, ok := rule.(*FilterRuleRegex); !ok {
		t.Fatalf("tipo inesperado para filter rule")
	}

	dataUnknown := readFixture(t, "filter_unknown.json")
	unknown, err := unmarshallFilterRule(dataUnknown)
	if err != nil {
		t.Fatalf("error inesperado para type desconocido: %v", err)
	}
	if unknown != nil {
		t.Fatal("se esperaba nil para type desconocido")
	}
}

func TestUnmarshallerFunctions(t *testing.T) {
	extData := readFixture(t, "transformation_extension.json")
	if _, err := extensionTransformation(extData); err != nil {
		t.Fatalf("extensionTransformation devolvio error: %v", err)
	}

	pathData := readFixture(t, "transformation_path_change.json")
	if _, err := pathChangeTransformation(pathData); err != nil {
		t.Fatalf("pathChangeTransformation devolvio error: %v", err)
	}

	regexData := readFixture(t, "filter_regex.json")
	if _, err := regexFilterTransformation(regexData); err != nil {
		t.Fatalf("regexFilterTransformation devolvio error: %v", err)
	}

	filterExtData := readFixture(t, "filter_extension.json")
	if _, err := extensionFilterTransformation(filterExtData); err != nil {
		t.Fatalf("extensionFilterTransformation devolvio error: %v", err)
	}
}

func TestMovementRunUnmarshalJSON(t *testing.T) {
	data := readFixture(t, "movement_run_valid.json")
	var run MovementRun
	if err := json.Unmarshal(data, &run); err != nil {
		t.Fatalf("json.Unmarshal devolvio error: %v", err)
	}
	if run.Source != "origin" {
		t.Fatalf("source inesperado: %s", run.Source)
	}
	if len(run.ChangeKeyMap) != 1 || run.ChangeKeyMap[0].Key != "filename" {
		t.Fatalf("change_key_map inesperado: %+v", run.ChangeKeyMap)
	}
	if len(run.TransformationRules) != 2 {
		t.Fatalf("transformation rules inesperadas: %d", len(run.TransformationRules))
	}
	if len(run.FilterRules) != 2 {
		t.Fatalf("filter rules inesperadas: %d", len(run.FilterRules))
	}
	if _, ok := run.TransformationRules[0].(*TransformationRuleExtension); !ok {
		t.Fatal("transformation_rules[0] deberia ser TransformationRuleExtension")
	}
	if _, ok := run.TransformationRules[1].(*TransformationRulePathChange); !ok {
		t.Fatal("transformation_rules[1] deberia ser TransformationRulePathChange")
	}
	if _, ok := run.FilterRules[0].(*FilterRuleRegex); !ok {
		t.Fatal("filter_rules[0] deberia ser FilterRuleRegex")
	}
	if _, ok := run.FilterRules[1].(*FilterRuleExtension); !ok {
		t.Fatal("filter_rules[1] deberia ser FilterRuleExtension")
	}
}

func TestMovementConfigurationUnmarshalJSON(t *testing.T) {
	data := readFixture(t, "movement_config_valid.json")
	var cfg MovementConfiguration
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("json.Unmarshal devolvio error: %v", err)
	}
	if !cfg.DryRun {
		t.Fatal("dry_run deberia ser true")
	}
	if len(cfg.Movements) != 1 {
		t.Fatalf("movements inesperado: %d", len(cfg.Movements))
	}
}

func TestApplyInsertions(t *testing.T) {
	run := MovementRun{}
	got := run.applyInsertions("/dest/{filename}.{ext}", map[string]string{"filename": "image", "ext": "png"})
	if got != "/dest/image.png" {
		t.Fatalf("resultado inesperado: %s", got)
	}
}

func TestApplyTransformations(t *testing.T) {
	run := MovementRun{
		TransformationRules: []TransformationRule{
			&TransformationRulePathChange{From: "origin", To: "processed"},
			&TransformationRuleExtension{Extensions: []ExtensionDuo{{From: ".jpg", To: ".png"}}},
		},
	}
	if got := run.ApplyTransformations("origin/image.jpg"); got != "processed/image.png" {
		t.Fatalf("resultado inesperado: %s", got)
	}
}
