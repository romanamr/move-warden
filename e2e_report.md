# E2E CLI Reporte

Ejecucion en cascada: dry-run y luego real-run si dry-run pasa.

## Resultados

### Caso 1: archivo .txt -> .md
- Estado final: `FAIL_DRY`

#### Dry-run
- Estado: `FAIL`
- Motivo: La app reporto error de ejecucion del motor.
- Comando: `go run . --dry-run --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case1_single_file.json`
- Exit code: `0`
- Stdout: `Dry run: true Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case1_single_file.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case1_single_file/single.txt false [] [0x130340a228a0] []}]} Error: failed to run engine open /home/keroveros/Proyectos/move-warden/e2e_generated/case1_single_file/single.txt: ...`
- Stderr: `2026/03/20 13:14:05 Procesando movimiento: /home/keroveros/Proyectos/move-warden/e2e_generated/case1_single_file/single.txt 2026/03/20 13:14:05 Mapping de variables: map[ext:txt filename:single fragment_0: fragment_1:home fragment_2:keroveros fragment_3:Proyectos fragment_4:move-warden fragment_5:e2e_generated fragment_6:case1_single_file fragme...`

#### Real-run
- Estado: `SKIPPED`
- Motivo: Real-run omitido porque dry-run no paso.
- Comando: `go run . --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case1_single_file.json`
- Exit code: `0`
- Stdout: `(sin salida)`
- Stderr: `(sin salida)`

### Caso 2: carpeta con 2 .txt -> .md
- Estado final: `FAIL_REAL`

#### Dry-run
- Estado: `PASS`
- Motivo: Dry-run ejecutado correctamente.
- Comando: `go run . --dry-run --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case2_folder_two_files.json`
- Exit code: `0`
- Stdout: `Dry run: true Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case2_folder_two_files.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case2_folder_two_files/docs false [] [0xeca17ff4c90] [0xeca17ff5080]}]}`
- Stderr: `2026/03/20 13:14:05 Dry run: /home/keroveros/Proyectos/move-warden/e2e_generated/case2_folder_two_files/docs → /home/keroveros/Proyectos/move-warden/e2e_generated/case2_folder_two_files/docs`

#### Real-run
- Estado: `FAIL`
- Motivo: La app reporto error de ejecucion del motor.
- Comando: `go run . --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case2_folder_two_files.json`
- Exit code: `0`
- Stdout: `Dry run: false Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case2_folder_two_files.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case2_folder_two_files/docs false [] [0x3ea6523aacc0] [0x3ea6523ab0b0]}]} Error: failed to run engine rename /home/keroveros/Proyectos/move-warden/e2e_generated/case2_fol...`
- Stderr: `(sin salida)`

### Caso 3: reorganizacion recursiva por tipo (JPG/TIFF/PDF)
- Estado final: `FAIL_REAL`

#### Dry-run
- Estado: `PASS`
- Motivo: Dry-run ejecutado correctamente.
- Comando: `go run . --dry-run --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case3_recursive_by_type.json`
- Exit code: `0`
- Stdout: `Dry run: true Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case3_recursive_by_type.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case3_recursive_by_type/Algo false [{algo2_name Algo2}] [0x2d728a1e2cc0 0x2d728a1e2e40] []}]}`
- Stderr: `2026/03/20 13:14:05 Dry run: /home/keroveros/Proyectos/move-warden/e2e_generated/case3_recursive_by_type/Algo → /home/keroveros/Proyectos/move-warden/e2e_generated/case3_recursive_by_type/Algo`

#### Real-run
- Estado: `FAIL`
- Motivo: La app reporto error de ejecucion del motor.
- Comando: `go run . --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case3_recursive_by_type.json`
- Exit code: `0`
- Stdout: `Dry run: false Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case3_recursive_by_type.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case3_recursive_by_type/Algo false [{algo2_name Algo2}] [0x2648d42aacf0 0x2648d42aae70] []}]} Error: failed to run engine rename /home/keroveros/Proyectos/move-warden/e2...`
- Stderr: `(sin salida)`
