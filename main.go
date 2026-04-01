package main

import (
	"flag"
	"fmt"
	"movewarden/internal/config"
	"movewarden/internal/engine"
	"os"
)

func main() {
	// Lectura de parametros de entrada dry run y rules
	dryRun := flag.Bool("dry-run", false, "Dry run")
	rules := flag.String("rules", "", "Rules file")
	// Para TUI: coleccionar los movimientos para luego mostrarlos en la interfaz
	collectMove := flag.Bool("interactive", false, "Colectar los movimientos para luego mostrarlos en la interfaz")
	// Para generar un ejemplo de rules.json con un caso realista de movimientos, filtros y transformaciones
	createExampleRules := flag.Bool("create-example-rules", false, "Genera example_rules.json con un caso realista de movimientos, filtros y transformaciones")
	flag.Parse()

	if *createExampleRules {
		createExampleRulesFile()
		return
	}

	fmt.Println("Dry run:", *dryRun)
	fmt.Println("Rules:", *rules)

	if *rules == "" {
		fmt.Println("Error: rules file is required")
		return
	}

	rulesData, err := os.ReadFile(*rules)
	if err != nil {
		fmt.Println("Error: failed to read rules file", err)
		return
	}

	// Ejecucion del motor de movimientos
	cfg := config.MovementConfiguration{}
	// Tiene su propio unmarshaller para manejar los movimientos y sus reglas de transformación y filtrado
	err = cfg.UnmarshalJSON(rulesData)
	if err != nil {
		fmt.Println("Error: failed to unmarshal rules data", err)
		return
	}

	fmt.Printf("Configuración de movimientos cargada correctamente: %d movimiento(s) encontrado(s) dryrun: %v\n, delete_empty_directories: %v\n", len(cfg.Movements), *dryRun, cfg.DeleteEmptyDirectories)
	for i, m := range cfg.Movements {
		fmt.Printf("  [%d] %s\n", i, m.Source)
	}

	appRunConfig := config.AppRunConfig{
		DryRun:   *dryRun,
		FilePath: *rules,
	}

	// Si es dry run tenemos que usar la funcion dryRunMove, si no, usamos la funcion realMove
	movePlans := []engine.MovePlan{}
	var moveFunc engine.MoveFunc
	if *dryRun {
		moveFunc = engine.ExecuteDryRunMove
	} else if *collectMove {
		moveFunc = engine.ExecuteCollectMove(&movePlans)
	} else {
		moveFunc = engine.ExecuteRealMove
	}
	err = engine.Run(cfg, appRunConfig, moveFunc)
	if err != nil {
		fmt.Println("Error: failed to run engine", err)
		os.Exit(1)
	}
	if *collectMove {
		fmt.Println("Move plans:", movePlans)
		os.Exit(0)
	}
	os.Exit(0)
}

func createExampleRulesFile() {
	err := engine.CreateExampleRulesFile()
	if err != nil {
		fmt.Println("Error: failed to create example rules file", err)
		return
	}
	fmt.Println("Archivo example_rules.json creado correctamente.")
	fmt.Println("Incluye un ejemplo base con dos movimientos (photos/docs), reglas de transformacion y filtros.")
	fmt.Println("Puedes usarlo como plantilla y ejecutarlo con: go run . --dry-run --rules example_rules.json")
}
