package engine

import (
	"encoding/json"
	"io"
	"log"
	"movewarden/internal/config"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Fragment struct {
	Value    string
	Position int
}

func getMappingVariables(path string) map[string]string {
	fragments := loadFragment(path)
	mapping := make(map[string]string)
	for pos, fragment := range fragments {
		// ? Por ahora solo soportamos filename y ext, pero en el futuro podríamos soportar mas variables como parent_dir, etc.
		if strings.Contains(fragment.Value, ".") {
			parts := strings.SplitN(fragment.Value, ".", 2)
			mapping["filename"] = parts[0]
			mapping["ext"] = parts[1]
			if pos > 0 {
				mapping["parent_dir"] = fragments[pos-1].Value
			} else {
				mapping["parent_dir"] = ""
			}
		} else {
			mapping["filename"] = fragment.Value
			mapping["ext"] = ""
		}
		// Mapeamos fragment_{pos} a fragment.Value para poder usarlo en las reglas de transformación y filtrado
		mapping["fragment_"+strconv.Itoa(pos)] = fragment.Value
	}
	// Ahora agregamos tambien fragment_last y fragment_init
	if len(fragments) > 0 {
		mapping["fragment_last"] = fragments[len(fragments)-1].Value
		mapping["fragment_init"] = fragments[0].Value
	} else {
		mapping["fragment_last"] = ""
		mapping["fragment_init"] = ""
	}
	return mapping
}

func loadFragment(path string) []Fragment {
	// puede ser / o \\ dependiendo del sistema operativo, asi que mejor usar strings.Split con ambos
	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		// Si es ruta de red tenemos que ignorar los primeros \\
		path = strings.TrimPrefix(path, `\\`)
		parts = strings.Split(path, "\\")
	}
	fragments := make([]Fragment, len(parts))
	for i, part := range parts {
		fragments[i] = Fragment{Value: part, Position: i}
	}
	return fragments
}

// A partir de una ruta recorremos de forma recursiva y si el directorio esta vacio lo eliminamos, y asi con cada directorio padre hasta llegar a la raiz o a un directorio que no este vacio
func removeEmptyDirs(path string) error {
	entries, err := os.ReadDir(path)
	/// si path es un fichero y no un directorio, tomamos el directorio padre
	if isFile(path) {
		path = filepath.Dir(path)
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			child := filepath.Join(path, entry.Name())
			if err = removeEmptyDirs(child); err != nil {
				return err
			}

		}
	}
	if isEmptyDir(path) {
		log.Printf("El directorio %s está vacío, eliminando...", path)
		return os.Remove(path)
	}

	return nil
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func isEmptyDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()
	_, err = f.Readdirnames(1) // Intentamos leer un nombre de archivo del directorio
	return err == io.EOF       // Si obtenemos EOF significa que el directorio está vacío
}

func getFiles(path string) []os.FileInfo {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	files, err := f.Readdir(-1)
	if err != nil {
		return nil
	}
	return files
}

func createDestDirIfNotExist(destination string) error {
	dir := filepath.Dir(destination)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("El directorio %s no existe, creando...", dir)
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

func executeToFiles(path string, configMovement config.MovementRun,
	fnDir func(string, config.MovementRun) error,
	fnFile func(string, config.MovementRun) error) error {
	files := getFiles(path)
	for _, file := range files {
		if file.IsDir() {
			if err := executeToFiles(filepath.Join(path, file.Name()), configMovement, fnDir, fnFile); err != nil {
				return err
			}
			if fnDir != nil {
				if err := fnDir(filepath.Join(path, file.Name()), configMovement); err != nil {
					return err
				}
			}
		} else {
			if fnFile != nil {
				if err := fnFile(filepath.Join(path, file.Name()), configMovement); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CreateExampleRulesFile() error {
	movementConfiguration := config.MovementConfiguration{
		DryRun:                 true,
		DeleteEmptyDirectories: false,
		Movements: []config.MovementRun{
			{
				Source:    "workspace/inbox/photos",
				Recursive: true,
				ChangeKeyMap: []config.ChangeKey{
					{Key: "parent_dir", Value: "fotos"},
				},
				TransformationRules: []config.TransformationRule{
					&config.TransformationRulePathChange{
						Type: "path_change",
						From: "workspace/inbox/photos",
						To:   "workspace/library/{parent_dir}",
					},
					&config.TransformationRuleExtension{
						Type: "extension",
						Extensions: []config.ExtensionDuo{
							{From: ".jpeg", To: ".jpg"},
							{From: ".heic", To: ".jpg"},
						},
					},
					&config.TransformationRuleRegex{
						Type:        "regex",
						Pattern:     " ",
						Replacement: "_",
					},
				},
				FilterRules: []config.FilterRule{
					&config.FilterRuleExtension{
						Type:       "extension",
						Extensions: []string{".jpg", ".jpeg", ".png", ".heic"},
					},
				},
			},
			{
				Source:       "workspace/inbox/docs",
				Recursive:    true,
				ChangeKeyMap: []config.ChangeKey{},
				TransformationRules: []config.TransformationRule{
					&config.TransformationRulePathChange{
						Type: "path_change",
						From: "workspace/inbox/docs",
						To:   "workspace/library/docs",
					},
					&config.TransformationRuleExtension{
						Type: "extension",
						Extensions: []config.ExtensionDuo{
							{From: ".txt", To: ".md"},
						},
					},
				},
				FilterRules: []config.FilterRule{
					&config.FilterRuleExtension{
						Type:       "extension",
						Extensions: []string{".txt", ".md", ".pdf"},
					},
					&config.FilterRuleRegex{
						Type:    "regex",
						Pattern: ".*",
					},
				},
			},
		},
	}

	jsonData, err := json.MarshalIndent(movementConfiguration, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling movement configuration: %v", err)
	}
	return os.WriteFile("example_rules.json", jsonData, 0644)
}
