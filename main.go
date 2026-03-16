package main

import (
	"flag"
	"fmt"
)

func main() {
	// Lectura de parametros de entrada dry run y rules
	dryRun := flag.Bool("dry-run", false, "Dry run")
	rules := flag.String("rules", "", "Rules file")
	flag.Parse()

	fmt.Println("Dry run:", *dryRun)
	fmt.Println("Rules:", *rules)
}
