# E2E CLI Reporte

Ejecucion en cascada: dry-run y luego real-run si dry-run pasa.

## Resultados

### Caso 1: archivo .txt -> .md
- Estado final: `PASS`

#### Dry-run
- Estado: `PASS`
- Motivo: Dry-run ejecutado correctamente.
- Comando: `go run . --dry-run --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case1_single_file.json`
- Exit code: `0`
- Stdout: `Dry run: true Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case1_single_file.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case1_single_file/single.txt false [] [0x3d23d1274c90] []}]}`
- Stderr: `2026/03/20 15:10:03 Procesando movimiento: /home/keroveros/Proyectos/move-warden/e2e_generated/case1_single_file/single.txt 2026/03/20 15:10:03 Mapping de variables: map[ext:txt filename:single fragment_0: fragment_1:home fragment_2:keroveros fragment_3:Proyectos fragment_4:move-warden fragment_5:e2e_generated fragment_6:case1_single_file fragme...`

#### Real-run
- Estado: `PASS`
- Motivo: Real-run ejecutado correctamente.
- Comando: `go run . --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case1_single_file.json`
- Exit code: `0`
- Stdout: `Dry run: false Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case1_single_file.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case1_single_file/single.txt false [] [0x35ba0c49acf0] []}]}`
- Stderr: `2026/03/20 15:10:03 Procesando movimiento: /home/keroveros/Proyectos/move-warden/e2e_generated/case1_single_file/single.txt 2026/03/20 15:10:03 Mapping de variables: map[ext:txt filename:single fragment_0: fragment_1:home fragment_2:keroveros fragment_3:Proyectos fragment_4:move-warden fragment_5:e2e_generated fragment_6:case1_single_file fragme...`

### Caso 2: carpeta con 2 .txt -> .md
- Estado final: `PASS`

#### Dry-run
- Estado: `PASS`
- Motivo: Dry-run ejecutado correctamente.
- Comando: `go run . --dry-run --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case2_folder_two_files.json`
- Exit code: `0`
- Stdout: `Dry run: true Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case2_folder_two_files.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case2_folder_two_files/docs true [] [0x3e9fc3ea08a0] [0x3e9fc3ea0c90]}]}`
- Stderr: `2026/03/20 15:10:03 Procesando movimiento: /home/keroveros/Proyectos/move-warden/e2e_generated/case2_folder_two_files/docs/a.txt 2026/03/20 15:10:03 Mapping de variables: map[ext:txt filename:a fragment_0: fragment_1:home fragment_2:keroveros fragment_3:Proyectos fragment_4:move-warden fragment_5:e2e_generated fragment_6:case2_folder_two_files f...`

#### Real-run
- Estado: `PASS`
- Motivo: Real-run ejecutado correctamente.
- Comando: `go run . --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case2_folder_two_files.json`
- Exit code: `0`
- Stdout: `Dry run: false Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case2_folder_two_files.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case2_folder_two_files/docs true [] [0x165d40186cf0] [0x165d401870e0]}]}`
- Stderr: `2026/03/20 15:10:03 Procesando movimiento: /home/keroveros/Proyectos/move-warden/e2e_generated/case2_folder_two_files/docs/a.txt 2026/03/20 15:10:03 Mapping de variables: map[ext:txt filename:a fragment_0: fragment_1:home fragment_2:keroveros fragment_3:Proyectos fragment_4:move-warden fragment_5:e2e_generated fragment_6:case2_folder_two_files f...`

### Caso 3: reorganizacion recursiva por tipo (JPG/TIFF/PDF)
- Estado final: `PASS`

#### Dry-run
- Estado: `PASS`
- Motivo: Dry-run ejecutado correctamente.
- Comando: `go run . --dry-run --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case3_recursive_by_type.json`
- Exit code: `0`
- Stdout: `Dry run: true Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case3_recursive_by_type.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case3_recursive_by_type/Algo true [{algo2_name Algo2}] [0xb04299b2cc0 0xb04299b2e40] []}]}`
- Stderr: `2026/03/20 15:10:03 Procesando movimiento: /home/keroveros/Proyectos/move-warden/e2e_generated/case3_recursive_by_type/Algo/algo/JPG/placeholder.txt 2026/03/20 15:10:03 Mapping de variables: map[ext:txt filename:placeholder fragment_0: fragment_1:home fragment_10:placeholder.txt fragment_2:keroveros fragment_3:Proyectos fragment_4:move-warden fr...`

#### Real-run
- Estado: `PASS`
- Motivo: Real-run ejecutado correctamente.
- Comando: `go run . --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case3_recursive_by_type.json`
- Exit code: `0`
- Stdout: `Dry run: false Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case3_recursive_by_type.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case3_recursive_by_type/Algo true [{algo2_name Algo2}] [0xd209d96ed20 0xd209d96eea0] []}]}`
- Stderr: `2026/03/20 15:10:03 Procesando movimiento: /home/keroveros/Proyectos/move-warden/e2e_generated/case3_recursive_by_type/Algo/algo/JPG/placeholder.txt 2026/03/20 15:10:03 Mapping de variables: map[ext:txt filename:placeholder fragment_0: fragment_1:home fragment_10:placeholder.txt fragment_2:keroveros fragment_3:Proyectos fragment_4:move-warden fr...`

### Caso 4: combinacion de caso 2 y caso 3
- Estado final: `PASS`

#### Dry-run
- Estado: `PASS`
- Motivo: Dry-run ejecutado correctamente.
- Comando: `go run . --dry-run --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case4_combo_case2_case3.json`
- Exit code: `0`
- Stdout: `Dry run: true Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case4_combo_case2_case3.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case4_combo/docs_block/docs true [] [0x24eb05f2acc0] [0x24eb05f2b0b0]} {/home/keroveros/Proyectos/move-warden/e2e_generated/case4_combo/algo_block/Algo true [{algo2_name ...`
- Stderr: `2026/03/20 15:10:03 Procesando movimiento: /home/keroveros/Proyectos/move-warden/e2e_generated/case4_combo/docs_block/docs/a.txt 2026/03/20 15:10:03 Mapping de variables: map[ext:txt filename:a fragment_0: fragment_1:home fragment_2:keroveros fragment_3:Proyectos fragment_4:move-warden fragment_5:e2e_generated fragment_6:case4_combo fragment_7:d...`

#### Real-run
- Estado: `PASS`
- Motivo: Real-run ejecutado correctamente.
- Comando: `go run . --rules /home/keroveros/Proyectos/move-warden/e2e_test_config/case4_combo_case2_case3.json`
- Exit code: `0`
- Stdout: `Dry run: false Rules: /home/keroveros/Proyectos/move-warden/e2e_test_config/case4_combo_case2_case3.json Rules: {true false [{/home/keroveros/Proyectos/move-warden/e2e_generated/case4_combo/docs_block/docs true [] [0xc603fcb6d20] [0xc603fcb7110]} {/home/keroveros/Proyectos/move-warden/e2e_generated/case4_combo/algo_block/Algo true [{algo2_name A...`
- Stderr: `2026/03/20 15:10:03 Procesando movimiento: /home/keroveros/Proyectos/move-warden/e2e_generated/case4_combo/docs_block/docs/a.txt 2026/03/20 15:10:03 Mapping de variables: map[ext:txt filename:a fragment_0: fragment_1:home fragment_2:keroveros fragment_3:Proyectos fragment_4:move-warden fragment_5:e2e_generated fragment_6:case4_combo fragment_7:d...`
