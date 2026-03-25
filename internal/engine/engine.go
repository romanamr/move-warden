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
		if cfg.DeleteEmptyDirectories {
			// Primero sacamos el root del movimiento
			cleanupRoot := getParentDirectory(movement.Source)
			// Luego eliminamos los directorios vacíos
			if cleanupRoot == "" {
				continue
			}
			defer func() {
				errc := removeEmptyDirs(cleanupRoot)
				if errc != nil {
					log.Printf("Error al eliminar directorios vacíos: %v", errc)
				}
			}()
		}
		err := runMovement(movement, moveFunc)
		if err != nil {
			return err
		}
		// Si esta marcado remove_empty_directories, eliminamos los directorios vacíos
		if !cfg.DeleteEmptyDirectories {
			continue
		}
		errc := removeEmptyDirs(movement.Source)
		if errc != nil {
			return errc
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
		// en caso de que sea un directorio volvera a llamar a runDirectory para cada fichero y subdirectorio asi que crearemos
		return executeToFiles(movement.Source, movement, executeRunMovement(movement, moveFunc), executeRunFile(movement, moveFunc))
	}
	return moveDirectory(movement.Source, movement, moveFunc)
}

func moveDirectory(source string, movement config.MovementRun, moveFunc MoveFunc) error {
	destination := processDestination(source, movement)
	// Si el destino es el mismo que el origen, no hacemos nada
	if movement.Source == destination {
		return nil
	}
	return moveFunc(source, destination)
}

func runFile(movement config.MovementRun, moveFunc MoveFunc) error {
	// Primero aplicamos filtros si los hay, si no es allowed, no hacemos nada
	if !movement.AllowedByFilters(movement.Source) {
		return nil
	}
	// Primero obtenemos el mapping de variables para la ruta origen
	mapping := getMappingVariables(movement.Source)
	// Luego aplicamos las reglas de transformación y filtrado en orden
	destination := movement.Process(movement.Source, mapping)
	if samePath(movement.Source, destination) {
		return nil
	}
	log.Printf("Procesando movimiento: %s", movement.Source)
	// Le solicitamos al sistema que mueva el fichero de origen a destino, pero si es dry_run entonces solo imprimimos lo que haríamos sin hacer nada realmente
	return moveFunc(movement.Source, destination)
}

// Ejecutar de verdad
func ExecuteRealMove(src, dst string) error {
	if err := createDestDirIfNotExist(dst); err != nil {
		return err
	}
	// Ahora hay que validar que el destino no exista, si existe, hay que renombrarlo o eliminarlo según la configuración, pero para simplificar, vamos a asumir que no existe y si existe, se sobreescribe.
	err := os.Rename(src, dst)
	if err != nil {
		//Validamos que el destino existe y si no existe es un error.
		if os.IsNotExist(err) {
			// Tenemos que crear el error de destino
			return os.ErrNotExist
		}
	}

	return err
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

// Funcion que devuelve un MoveFunc para la ejecucion recursiva de un directorio
func executeRunMovement(movement config.MovementRun, moveFunc MoveFunc) MoveFunc {
	return func(src, dst string) error {
		recMovement := movement.Clone()
		recMovement.Source = src
		return runMovement(recMovement, moveFunc)
	}
}

func executeRunFile(movement config.MovementRun, moveFunc MoveFunc) MoveFunc {
	return func(src, dst string) error {
		recMovement := movement.Clone()
		recMovement.Source = src
		return runFile(recMovement, moveFunc)
	}
}

type MoveFunc func(src, dst string) error

type MovePlan struct {
	Source      string
	Destination string
}
