package engine

import (
	"log"
	"movewarden/internal/config"
	"os"
)

/*
	Necesitamos un tipo para guardar los fragmentos de ruta que despues pasaremos com fragment.position Cuando se pida cambiar los {}
	por ejemplo en el caso de un movimiento con source /tmp/origin/file.jpg y destination /tmp/dest/{filename}.{ext} tendriamos los siguientes fragmentos:
[
	{Value: "tmp", Position: 0},
	{Value: "origin", Position: 1},
	{Value: "file.jpg", Position: 2}
]

Luego, al aplicar las reglas de transformación y filtrado, iríamos actualizando el valor de cada fragmento según las reglas aplicadas, y al final reconstruiríamos la ruta destino reemplazando los {} por los valores actualizados de los fragmentos.
*/

func Run(cfg config.MovementConfiguration, runconfig config.AppRunConfig, moveFunc MoveFunc) error {
	// Aquí iría la lógica para ejecutar las reglas de transformación y filtrado
	for _, movement := range cfg.Movements {
		err := runMovement(movement, moveFunc)
		if err != nil {
			return err
		}
		err = removeEmptyDirs(movement.Source)
		if err != nil {
			return err
		}
	}
	return nil
}

func runMovement(movement config.MovementRun, moveFunc MoveFunc) error {
	// Aquí iría la lógica para ejecutar las reglas de transformación y filtrado de cada movimiento
	// Aqui vemos si es fichero o directorio y aplicamos las reglas correspondientes
	fileInfo, err := os.Stat(movement.Source)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return runDirectory(movement, moveFunc)
	}
	return runFile(movement, moveFunc)
}

func runDirectory(movement config.MovementRun, moveFunc MoveFunc) error {
	if movement.Recursive {
		return executeToFiles(movement.Source, movement, executeDirectory, executeFile)
	}
	destination := movement.Process(movement.Source, getMappingVariables(movement.Source))
	return moveFunc(movement.Source, destination)
}

func moveDirectory(source string, movement config.MovementRun, moveFunc MoveFunc) error {
	log.Printf("Moviendo directorio %s con reglas de movimiento: %+v", source, movement)
	destination := movement.Process(source, getMappingVariables(source))
	if movement.Recursive {
		log.Printf("El movimiento es recursivo, se moverán también los ficheros y subdirectorios dentro de %s", source)
	}
	return moveFunc(source, destination)
}

func executeDirectory(str string, movement config.MovementRun) error {
	// Aquí iría la lógica para ejecutar las reglas de transformación y filtrado de cada directorio, pero por ahora solo imprimimos lo que haríamos
	log.Printf("Procesando directorio: %s", str)
	configMovement := movement.Clone()
	configMovement.Source = str
	error := runDirectory(configMovement, ExecuteRealMove)
	if error != nil {
		return error
	}
	return nil
}

func executeFile(str string, movement config.MovementRun) error {
	configMovement := movement
	configMovement.Source = str
	// Aquí iría la lógica para ejecutar las reglas de transformación y filtrado de cada fichero, pero por ahora solo imprimimos lo que haríamos
	error := runFile(configMovement, ExecuteRealMove)
	if error != nil {
		return error
	}
	return nil
}

func runFile(movement config.MovementRun, moveFunc MoveFunc) error {
	// Primero obtenemos el mapping de variables para la ruta origen
	log.Printf("Procesando movimiento: %s", movement.Source)
	mapping := getMappingVariables(movement.Source)
	// Luego aplicamos las reglas de transformación y filtrado en orden
	log.Printf("Mapping de variables: %+v", mapping)
	destination := movement.Process(movement.Source, mapping)

	// Le solicitamos al sistema que mueva el fichero de origen a destino, pero si es dry_run entonces solo imprimimos lo que haríamos sin hacer nada realmente
	return moveFunc(movement.Source, destination)
}

// Ejecutar de verdad
func ExecuteRealMove(src, dst string) error {
	if err := createDestDirIfNotExist(dst); err != nil {
		return err
	}
	return os.Rename(src, dst)
}

// Dry run
func ExecuteDryRunMove(src, dst string) error {
	log.Printf("Dry run: %s → %s", src, dst)
	return nil
}

// Colectar para TUI: closure que captura el slice.
// Se usa en tests mientras se integra en el flujo de TUI.
func ExecuteCollectMove(plans *[]MovePlan) MoveFunc {
	return func(src, dst string) error {
		*plans = append(*plans, MovePlan{Source: src, Destination: dst})
		return nil
	}
}

type MoveFunc func(src, dst string) error

type MovePlan struct {
	Source      string
	Destination string
}
