package engine

import (
	"errors"
	"movewarden/internal/config"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestGetMappingVariables_ExtraeFilenameExtYFragmentos(t *testing.T) {
	path := "/tmp/origin/file.jpg"
	mapping := getMappingVariables(path)
	expected := map[string]string{
		"filename":      "file",
		"ext":           "jpg",
		"parent_dir":    "origin",
		"fragment_0":    "",
		"fragment_1":    "tmp",
		"fragment_2":    "origin",
		"fragment_3":    "file.jpg",
		"fragment_init": "",
		"fragment_last": "file.jpg",
	}
	for key, expectedValue := range expected {
		if value, ok := mapping[key]; !ok || value != expectedValue {
			t.Fatalf("para key %s se esperaba %s pero se obtuvo %s", key, expectedValue, value)
		}
	}
}

func TestCreateDestDirIfNotExist_CreaRutaCompletaSiNoExiste(t *testing.T) {
	base := t.TempDir()
	destination := filepath.Join(base, "a", "b", "c", "file.txt")
	dir := filepath.Dir(destination)

	if err := createDestDirIfNotExist(destination); err != nil {
		t.Fatalf("no se esperaba error creando ruta destino: %v", err)
	}

	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		t.Fatalf("se esperaba directorio creado en %s", dir)
	}
}

func TestCreateDestDirIfNotExist_NoFallaSiYaExiste(t *testing.T) {
	base := t.TempDir()
	destination := filepath.Join(base, "file.txt")

	if err := createDestDirIfNotExist(destination); err != nil {
		t.Fatalf("no se esperaba error si el directorio ya existe: %v", err)
	}
}

func TestRunFile_DryRun_NoMueveFichero(t *testing.T) {
	base := t.TempDir()
	source := filepath.Join(base, "origen.txt")
	if err := os.WriteFile(source, []byte("contenido"), 0644); err != nil {
		t.Fatalf("error preparando fichero origen: %v", err)
	}

	movement := config.MovementRun{
		Source: source,
		TransformationRules: []config.TransformationRule{
			&config.TransformationRuleRegex{
				Pattern:     "origen\\.txt",
				Replacement: "destino.txt",
			},
		},
	}

	if err := runFile(movement, ExecuteDryRunMove); err != nil {
		t.Fatalf("no se esperaba error en dry_run: %v", err)
	}

	if _, err := os.Stat(source); err != nil {
		t.Fatalf("el fichero origen debería seguir existiendo en dry_run: %v", err)
	}

	destination := filepath.Join(base, "destino.txt")
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatalf("no debería existir destino en dry_run: %s", destination)
	}
}

func TestRunFile_MueveFicheroCuandoNoEsDryRun(t *testing.T) {
	base := t.TempDir()
	source := filepath.Join(base, "fichero_origen.txt")
	destination := filepath.Join(base, "fichero_destino.txt")

	if err := os.WriteFile(source, []byte("contenido de prueba"), 0644); err != nil {
		t.Fatalf("error preparando fichero origen: %v", err)
	}

	movement := config.MovementRun{
		Source: source,
		TransformationRules: []config.TransformationRule{
			&config.TransformationRuleRegex{Pattern: "fichero_origen\\.txt", Replacement: "fichero_destino.txt"},
		},
	}

	if err := runFile(movement, ExecuteRealMove); err != nil {
		t.Fatalf("no se esperaba error moviendo fichero: %v", err)
	}

	if _, err := os.Stat(destination); os.IsNotExist(err) {
		t.Fatalf("se esperaba que el fichero destino existiera: %s", destination)
	}

	if _, err := os.Stat(source); !os.IsNotExist(err) {
		t.Fatalf("se esperaba que el fichero origen fuera movido y no existiera: %s", source)
	}
}

func TestRunMovement_ErrorSiSourceNoExiste(t *testing.T) {
	movement := config.MovementRun{Source: filepath.Join(t.TempDir(), "no-existe.txt")}

	err := runMovement(movement, ExecuteRealMove)
	if err == nil {
		t.Fatal("se esperaba error cuando source no existe")
	}
}

func TestRunMovement_RutaFichero_UsaRunFile(t *testing.T) {
	base := t.TempDir()
	source := filepath.Join(base, "input.txt")
	destination := filepath.Join(base, "output.txt")

	if err := os.WriteFile(source, []byte("contenido"), 0644); err != nil {
		t.Fatalf("error preparando fichero origen: %v", err)
	}

	movement := config.MovementRun{
		Source: source,
		TransformationRules: []config.TransformationRule{
			&config.TransformationRuleRegex{Pattern: "input\\.txt", Replacement: "output.txt"},
		},
	}

	if err := runMovement(movement, ExecuteRealMove); err != nil {
		t.Fatalf("runMovement no debería fallar para fichero válido: %v", err)
	}

	if _, err := os.Stat(destination); os.IsNotExist(err) {
		t.Fatalf("se esperaba fichero destino en %s", destination)
	}
}

func TestMoveDirectory_DryRun_NoMueveDirectorio(t *testing.T) {
	base := t.TempDir()
	source := filepath.Join(base, "origen")
	if err := os.MkdirAll(source, 0755); err != nil {
		t.Fatalf("error creando directorio origen: %v", err)
	}

	movement := config.MovementRun{
		Source: source,
		TransformationRules: []config.TransformationRule{
			&config.TransformationRuleRegex{Pattern: "origen$", Replacement: "destino"},
		},
	}

	if err := moveDirectory(source, movement, ExecuteDryRunMove); err != nil {
		t.Fatalf("no se esperaba error en dry_run: %v", err)
	}

	if _, err := os.Stat(source); err != nil {
		t.Fatalf("el directorio origen debería seguir existiendo en dry_run: %v", err)
	}

	destination := filepath.Join(base, "destino")
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		t.Fatalf("el directorio destino no debería existir en dry_run")
	}
}

func TestExecuteToFiles_RecorreArbolYEjecutaCallbacks(t *testing.T) {
	base := t.TempDir()

	dirA := filepath.Join(base, "a")
	dirB := filepath.Join(dirA, "b")
	if err := os.MkdirAll(dirB, 0755); err != nil {
		t.Fatalf("error creando arbol de directorios: %v", err)
	}

	file1 := filepath.Join(base, "root.txt")
	file2 := filepath.Join(dirA, "a.txt")
	file3 := filepath.Join(dirB, "b.txt")

	for _, f := range []string{file1, file2, file3} {
		if err := os.WriteFile(f, []byte("x"), 0644); err != nil {
			t.Fatalf("error creando fichero %s: %v", f, err)
		}
	}

	dirsVisited := 0
	filesVisited := 0

	err := executeToFiles(
		base,
		config.MovementRun{},
		func(_ string, _ string) error {
			dirsVisited++
			return nil
		},
		func(_ string, _ string) error {
			filesVisited++
			return nil
		},
	)

	if err != nil {
		t.Fatalf("no se esperaba error recorriendo ficheros: %v", err)
	}

	if filesVisited != 3 {
		t.Fatalf("se esperaban 3 ficheros visitados y se obtuvieron %d", filesVisited)
	}

	if dirsVisited != 2 {
		t.Fatalf("se esperaban 2 directorios visitados y se obtuvieron %d", dirsVisited)
	}
}

func TestExecuteToFiles_PropagaErrorDeCallback(t *testing.T) {
	base := t.TempDir()
	file := filepath.Join(base, "x.txt")
	if err := os.WriteFile(file, []byte("x"), 0644); err != nil {
		t.Fatalf("error creando fichero temporal: %v", err)
	}

	expectedErr := errors.New("fallo callback")

	err := executeToFiles(
		base,
		config.MovementRun{},
		nil,
		func(_ string, _ string) error {
			return expectedErr
		},
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("se esperaba propagación de error callback")
	}
}

func TestCollectMove_AgregaPlan(t *testing.T) {
	plans := []MovePlan{}
	moveFn := ExecuteCollectMove(&plans)

	err := moveFn("source.txt", "dest.txt")
	if err != nil {
		t.Fatalf("no se esperaba error al recolectar movimiento: %v", err)
	}

	if len(plans) != 1 {
		t.Fatalf("se esperaba 1 plan y se obtuvo %d", len(plans))
	}

	if plans[0].Source != "source.txt" || plans[0].Destination != "dest.txt" {
		t.Fatalf("plan inesperado: %+v", plans[0])
	}
}

func TestRun_FileSource_DeleteEmptyDirectoriesTrue_EliminaDirectorioPadre(t *testing.T) {
	base := t.TempDir()
	originDir := filepath.Join(base, "origin")
	destDir := filepath.Join(base, "dest")
	if err := os.MkdirAll(originDir, 0755); err != nil {
		t.Fatalf("error creando originDir: %v", err)
	}
	source := filepath.Join(originDir, "file.txt")
	if err := os.WriteFile(source, []byte("x"), 0644); err != nil {
		t.Fatalf("error creando source: %v", err)
	}

	cfg := config.MovementConfiguration{
		DeleteEmptyDirectories: true,
		Movements: []config.MovementRun{
			{
				Source: source,
				TransformationRules: []config.TransformationRule{
					&config.TransformationRulePathChange{Type: "path_change", From: originDir, To: destDir},
					&config.TransformationRuleExtension{
						Type:       "extension",
						Extensions: []config.ExtensionDuo{{From: ".txt", To: ".md"}},
					},
				},
			},
		},
	}

	if err := Run(cfg, config.AppRunConfig{}, ExecuteRealMove); err != nil {
		t.Fatalf("Run no deberia fallar: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "file.md")); err != nil {
		t.Fatalf("se esperaba fichero movido en destino: %v", err)
	}
	if _, err := os.Stat(originDir); !os.IsNotExist(err) {
		t.Fatalf("se esperaba originDir eliminado por estar vacio")
	}
}

func TestRun_FileSource_DeleteEmptyDirectoriesFalse_ConservaDirectorioPadre(t *testing.T) {
	base := t.TempDir()
	originDir := filepath.Join(base, "origin")
	destDir := filepath.Join(base, "dest")
	if err := os.MkdirAll(originDir, 0755); err != nil {
		t.Fatalf("error creando originDir: %v", err)
	}
	source := filepath.Join(originDir, "file.txt")
	if err := os.WriteFile(source, []byte("x"), 0644); err != nil {
		t.Fatalf("error creando source: %v", err)
	}

	cfg := config.MovementConfiguration{
		DeleteEmptyDirectories: false,
		Movements: []config.MovementRun{
			{
				Source: source,
				TransformationRules: []config.TransformationRule{
					&config.TransformationRulePathChange{Type: "path_change", From: originDir, To: destDir},
					&config.TransformationRuleExtension{
						Type:       "extension",
						Extensions: []config.ExtensionDuo{{From: ".txt", To: ".md"}},
					},
				},
			},
		},
	}

	if err := Run(cfg, config.AppRunConfig{}, ExecuteRealMove); err != nil {
		t.Fatalf("Run no deberia fallar: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destDir, "file.md")); err != nil {
		t.Fatalf("se esperaba fichero movido en destino: %v", err)
	}
	if info, err := os.Stat(originDir); err != nil || !info.IsDir() {
		t.Fatalf("se esperaba originDir conservado cuando delete_empty_directories=false")
	}
}

func TestRunDirectory_Recursive_Caso2_CambiaExtensionDeArchivos(t *testing.T) {
	base := t.TempDir()
	docs := filepath.Join(base, "docs")
	if err := os.MkdirAll(docs, 0755); err != nil {
		t.Fatalf("error creando docs: %v", err)
	}
	for _, name := range []string{"a.txt", "b.txt"} {
		if err := os.WriteFile(filepath.Join(docs, name), []byte("x"), 0644); err != nil {
			t.Fatalf("error creando %s: %v", name, err)
		}
	}
	if err := os.WriteFile(filepath.Join(docs, "skip.pdf"), []byte("x"), 0644); err != nil {
		t.Fatalf("error creando skip.pdf: %v", err)
	}

	plans := []MovePlan{}
	movement := config.MovementRun{
		Source:    docs,
		Recursive: true,
		TransformationRules: []config.TransformationRule{
			&config.TransformationRuleExtension{
				Type:       "extension",
				Extensions: []config.ExtensionDuo{{From: ".txt", To: ".md"}},
			},
		},
		FilterRules: []config.FilterRule{
			&config.FilterRuleExtension{Type: "extension", Extensions: []string{".txt"}},
		},
	}

	if err := runDirectory(movement, ExecuteCollectMove(&plans)); err != nil {
		t.Fatalf("runDirectory recursivo no deberia fallar: %v", err)
	}

	got := []string{}
	for _, p := range plans {
		got = append(got, p.Destination)
	}
	slices.Sort(got)

	want := []string{
		filepath.Join(docs, "a.md"),
		filepath.Join(docs, "b.md"),
	}
	slices.Sort(want)

	if len(got) != len(want) {
		t.Fatalf("cantidad de planes inesperada: got=%d want=%d; plans=%+v", len(got), len(want), plans)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("destino[%d] inesperado: got=%s want=%s", i, got[i], want[i])
		}
	}
}

func TestRunDirectory_Recursive_Caso3_ReordenaPorTipo(t *testing.T) {
	base := t.TempDir()
	algoRoot := filepath.Join(base, "Algo")
	for _, branch := range []string{"algo", "algo2"} {
		for _, typ := range []string{"JPG", "TIFF", "PDF"} {
			dir := filepath.Join(algoRoot, branch, typ)
			if err := os.MkdirAll(dir, 0755); err != nil {
				t.Fatalf("error creando %s: %v", dir, err)
			}
			if err := os.WriteFile(filepath.Join(dir, "placeholder.txt"), []byte("x"), 0644); err != nil {
				t.Fatalf("error creando placeholder en %s: %v", dir, err)
			}
		}
	}

	plans := []MovePlan{}
	movement := config.MovementRun{
		Source:    algoRoot,
		Recursive: true,
		ChangeKeyMap: []config.ChangeKey{
			{Key: "algo2_name", Value: "Algo2"},
		},
		TransformationRules: []config.TransformationRule{
			&config.TransformationRuleRegex{
				Type:        "regex",
				Pattern:     `(.*)/Algo/algo/(JPG|TIFF|PDF)(/.*)?$`,
				Replacement: `${1}/Algo/${2}/Algo${3}`,
			},
			&config.TransformationRuleRegex{
				Type:        "regex",
				Pattern:     `(.*)/Algo/algo2/(JPG|TIFF|PDF)(/.*)?$`,
				Replacement: `${1}/Algo/${2}/{algo2_name}${3}`,
			},
		},
	}

	if err := runDirectory(movement, ExecuteCollectMove(&plans)); err != nil {
		t.Fatalf("runDirectory recursivo no deberia fallar: %v", err)
	}

	expected := map[string]bool{
		filepath.Join(algoRoot, "JPG", "Algo", "placeholder.txt"):   false,
		filepath.Join(algoRoot, "JPG", "Algo2", "placeholder.txt"):  false,
		filepath.Join(algoRoot, "TIFF", "Algo", "placeholder.txt"):  false,
		filepath.Join(algoRoot, "TIFF", "Algo2", "placeholder.txt"): false,
		filepath.Join(algoRoot, "PDF", "Algo", "placeholder.txt"):   false,
		filepath.Join(algoRoot, "PDF", "Algo2", "placeholder.txt"):  false,
	}

	for _, p := range plans {
		if _, ok := expected[p.Destination]; ok {
			expected[p.Destination] = true
		}
	}
	for dst, seen := range expected {
		if !seen {
			t.Fatalf("no se encontro destino esperado en planes: %s; plans=%+v", dst, plans)
		}
	}
}
