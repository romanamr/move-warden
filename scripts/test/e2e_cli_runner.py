#!/usr/bin/env python3
from __future__ import annotations

import json
import shutil
import subprocess
from dataclasses import asdict, dataclass
from pathlib import Path
from typing import Any, Callable


@dataclass
class StepResult:
    status: str
    reason: str
    command: str
    returncode: int
    stdout: str
    stderr: str
    skipped: bool = False


@dataclass
class CaseResult:
    case_name: str
    dry_run: StepResult
    real_run: StepResult
    final_status: str


ROOT = Path(__file__).resolve().parents[2]
E2E_GENERATED = ROOT / "e2e_generated"
E2E_CONFIG = ROOT / "e2e_test_config"
REPORT_PATH = ROOT / "e2e_report.md"
REPORT_JSON_PATH = ROOT / "e2e_report.json"


def cleanup_rogue_backref_dirs() -> None:
    """Remove accidental dirs like '\\1', '\\2foo' created by bad regex replacements."""
    for child in ROOT.iterdir():
        name = child.name
        if not child.is_dir():
            continue
        if len(name) >= 2 and name[0] == "\\" and name[1].isdigit():
            shutil.rmtree(child, ignore_errors=True)


def clean_and_prepare_dirs() -> None:
    cleanup_rogue_backref_dirs()
    for target in (E2E_GENERATED, E2E_CONFIG):
        if target.exists():
            shutil.rmtree(target)
        target.mkdir(parents=True, exist_ok=True)


def write_json(path: Path, payload: dict[str, Any]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload, indent=2), encoding="utf-8")


def run_cli(config_path: Path, dry_run: bool) -> subprocess.CompletedProcess[str]:
    cmd = ["go", "run", "."]
    if dry_run:
        cmd.append("--dry-run")
    cmd += ["--rules", str(config_path)]
    return subprocess.run(cmd, cwd=ROOT, capture_output=True, text=True)


def detect_not_supported(result: subprocess.CompletedProcess[str]) -> tuple[bool, str]:
    merged = f"{result.stdout}\n{result.stderr}".lower()
    if result.returncode == 0:
        return False, ""
    if "failed to run engine" in merged:
        return True, "El CLI ejecuta pero no puede completar el motor para este escenario."
    return False, ""


def build_step_result(
    result: subprocess.CompletedProcess[str], command: str, pass_reason: str
) -> StepResult:
    merged = f"{result.stdout}\n{result.stderr}".lower()
    not_supported, reason = detect_not_supported(result)
    if not_supported:
        return StepResult(
            status="NOT_SUPPORTED",
            reason=reason,
            command=command,
            returncode=result.returncode,
            stdout=result.stdout.strip(),
            stderr=result.stderr.strip(),
        )
    if result.returncode != 0:
        return StepResult(
            status="FAIL",
            reason="El comando devolvio error no clasificado como incompatibilidad conocida.",
            command=command,
            returncode=result.returncode,
            stdout=result.stdout.strip(),
            stderr=result.stderr.strip(),
        )
    if "error: failed to run engine" in merged:
        return StepResult(
            status="FAIL",
            reason="La app reporto error de ejecucion del motor.",
            command=command,
            returncode=result.returncode,
            stdout=result.stdout.strip(),
            stderr=result.stderr.strip(),
        )
    return StepResult(
        status="PASS",
        reason=pass_reason,
        command=command,
        returncode=result.returncode,
        stdout=result.stdout.strip(),
        stderr=result.stderr.strip(),
    )


def build_case1() -> Path:
    base = E2E_GENERATED / "case1_single_file"
    base.mkdir(parents=True, exist_ok=True)
    source = base / "single.txt"
    source.write_text("hola", encoding="utf-8")
    config = {
        "dry_run": True,
        "delete_empty_directories": False,
        "movements": [
            {
                "source": str(source),
                "recursive": False,
                "change_key_map": [],
                "transformation_rules": [
                    {"type": "extension", "extensions": [{"from": ".txt", "to": ".md"}]}
                ],
                "filter_rules": [],
            }
        ],
    }
    path = E2E_CONFIG / "case1_single_file.json"
    write_json(path, config)
    return path


def build_case2() -> Path:
    base = E2E_GENERATED / "case2_folder_two_files"
    source_dir = base / "docs"
    source_dir.mkdir(parents=True, exist_ok=True)
    (source_dir / "a.txt").write_text("a", encoding="utf-8")
    (source_dir / "b.txt").write_text("b", encoding="utf-8")
    config = {
        "dry_run": True,
        "delete_empty_directories": False,
        "movements": [
            {
                "source": str(source_dir),
                "recursive": True,
                "change_key_map": [],
                "transformation_rules": [
                    {"type": "extension", "extensions": [{"from": ".txt", "to": ".md"}]}
                ],
                "filter_rules": [{"type": "extension", "extensions": [".txt"]}],
            }
        ],
    }
    path = E2E_CONFIG / "case2_folder_two_files.json"
    write_json(path, config)
    return path


def build_case3() -> Path:
    base = E2E_GENERATED / "case3_recursive_by_type" / "Algo"
    for branch in ("algo", "algo2"):
        for file_type in ("JPG", "TIFF", "PDF"):
            src = base / branch / file_type
            src.mkdir(parents=True, exist_ok=True)
            (src / "placeholder.txt").write_text("x", encoding="utf-8")
    config = {
        "dry_run": True,
        "delete_empty_directories": False,
        "movements": [
            {
                "source": str(base),
                "recursive": True,
                "change_key_map": [{"key": "algo2_name", "value": "Algo2"}],
                "transformation_rules": [
                    {
                        "type": "regex",
                        "pattern": r"(.*)/Algo/algo/(JPG|TIFF|PDF)(/.*)?$",
                        "replacement": r"$1/Algo/$2/Algo$3",
                    },
                    {
                        "type": "regex",
                        "pattern": r"(.*)/Algo/algo2/(JPG|TIFF|PDF)(/.*)?$",
                        "replacement": r"$1/Algo/$2/{algo2_name}$3",
                    },
                ],
                "filter_rules": [],
            }
        ],
    }
    path = E2E_CONFIG / "case3_recursive_by_type.json"
    write_json(path, config)
    return path


def build_case4() -> Path:
    docs_root = E2E_GENERATED / "case4_combo" / "docs_block" / "docs"
    docs_root.mkdir(parents=True, exist_ok=True)
    (docs_root / "a.txt").write_text("a", encoding="utf-8")
    (docs_root / "b.txt").write_text("b", encoding="utf-8")

    algo_root = E2E_GENERATED / "case4_combo" / "algo_block" / "Algo"
    for branch in ("algo", "algo2"):
        for file_type in ("JPG", "TIFF", "PDF"):
            src = algo_root / branch / file_type
            src.mkdir(parents=True, exist_ok=True)
            (src / "placeholder.txt").write_text("x", encoding="utf-8")

    config = {
        "dry_run": True,
        "delete_empty_directories": False,
        "movements": [
            {
                "source": str(docs_root),
                "recursive": True,
                "change_key_map": [],
                "transformation_rules": [
                    {"type": "extension", "extensions": [{"from": ".txt", "to": ".md"}]}
                ],
                "filter_rules": [{"type": "extension", "extensions": [".txt"]}],
            },
            {
                "source": str(algo_root),
                "recursive": True,
                "change_key_map": [{"key": "algo2_name", "value": "Algo2"}],
                "transformation_rules": [
                    {
                        "type": "regex",
                        "pattern": r"(.*)/Algo/algo/(JPG|TIFF|PDF)(/.*)?$",
                        "replacement": r"$1/Algo/$2/Algo$3",
                    },
                    {
                        "type": "regex",
                        "pattern": r"(.*)/Algo/algo2/(JPG|TIFF|PDF)(/.*)?$",
                        "replacement": r"$1/Algo/$2/{algo2_name}$3",
                    },
                ],
                "filter_rules": [],
            },
        ],
    }
    path = E2E_CONFIG / "case4_combo_case2_case3.json"
    write_json(path, config)
    return path


def verify_case1_real() -> tuple[bool, str]:
    base = E2E_GENERATED / "case1_single_file"
    ok = (base / "single.md").exists() and not (base / "single.txt").exists()
    return ok, "Debe existir single.md y no single.txt tras real-run."


def verify_case2_real() -> tuple[bool, str]:
    source_dir = E2E_GENERATED / "case2_folder_two_files" / "docs"
    ok = (
        (source_dir / "a.md").exists()
        and (source_dir / "b.md").exists()
        and not (source_dir / "a.txt").exists()
        and not (source_dir / "b.txt").exists()
    )
    return ok, "Deben existir a.md y b.md en docs, sin a.txt/b.txt."


def verify_case3_real() -> tuple[bool, str]:
    base = E2E_GENERATED / "case3_recursive_by_type" / "Algo"
    expected = [
        base / "JPG" / "Algo",
        base / "JPG" / "Algo2",
        base / "TIFF" / "Algo",
        base / "TIFF" / "Algo2",
        base / "PDF" / "Algo",
        base / "PDF" / "Algo2",
    ]
    ok = all(path.is_dir() for path in expected)
    return ok, "Deben existir rutas esperadas por tipo para Algo y Algo2."


def verify_case4_real() -> tuple[bool, str]:
    docs_dir = E2E_GENERATED / "case4_combo" / "docs_block" / "docs"
    docs_ok = (
        (docs_dir / "a.md").exists()
        and (docs_dir / "b.md").exists()
        and not (docs_dir / "a.txt").exists()
        and not (docs_dir / "b.txt").exists()
    )
    algo_base = E2E_GENERATED / "case4_combo" / "algo_block" / "Algo"
    algo_expected = [
        algo_base / "JPG" / "Algo",
        algo_base / "JPG" / "Algo2",
        algo_base / "TIFF" / "Algo",
        algo_base / "TIFF" / "Algo2",
        algo_base / "PDF" / "Algo",
        algo_base / "PDF" / "Algo2",
    ]
    algo_ok = all(path.is_dir() for path in algo_expected)
    ok = docs_ok and algo_ok
    return ok, "Caso combinado exige exito estricto de conversion de docs y reordenamiento por tipo."


def evaluate_case(name: str, config_path: Path, verify_real_fn: Callable[[], tuple[bool, str]]) -> CaseResult:
    dry_cmd = f"go run . --dry-run --rules {config_path}"
    dry_proc = run_cli(config_path, dry_run=True)
    dry_step = build_step_result(dry_proc, dry_cmd, "Dry-run ejecutado correctamente.")
    if dry_step.status != "PASS":
        real_step = StepResult(
            status="SKIPPED",
            reason="Real-run omitido porque dry-run no paso.",
            command=f"go run . --rules {config_path}",
            returncode=0,
            stdout="",
            stderr="",
            skipped=True,
        )
        return CaseResult(case_name=name, dry_run=dry_step, real_run=real_step, final_status="FAIL_DRY")
    real_cmd = f"go run . --rules {config_path}"
    real_proc = run_cli(config_path, dry_run=False)
    real_step = build_step_result(real_proc, real_cmd, "Real-run ejecutado correctamente.")
    if real_step.status == "PASS":
        ok, reason = verify_real_fn()
        if not ok:
            real_step = StepResult(
                status="FAIL",
                reason=f"Validacion post real-run fallo: {reason}",
                command=real_cmd,
                returncode=real_proc.returncode,
                stdout=real_proc.stdout.strip(),
                stderr=real_proc.stderr.strip(),
            )
    final_status = "PASS" if (dry_step.status == "PASS" and real_step.status == "PASS") else "FAIL_REAL"
    return CaseResult(case_name=name, dry_run=dry_step, real_run=real_step, final_status=final_status)


def summarize_output(text: str, max_chars: int = 350) -> str:
    cleaned = " ".join(text.split())
    if not cleaned:
        return "(sin salida)"
    if len(cleaned) <= max_chars:
        return cleaned
    return cleaned[: max_chars - 3] + "..."


def write_report_md(results: list[CaseResult]) -> None:
    lines: list[str] = []
    lines.append("# E2E CLI Reporte")
    lines.append("")
    lines.append("Ejecucion en cascada: dry-run y luego real-run si dry-run pasa.")
    lines.append("")
    lines.append("## Resultados")
    lines.append("")
    for res in results:
        lines.append(f"### {res.case_name}")
        lines.append(f"- Estado final: `{res.final_status}`")
        lines.append("")
        lines.append("#### Dry-run")
        lines.append(f"- Estado: `{res.dry_run.status}`")
        lines.append(f"- Motivo: {res.dry_run.reason}")
        lines.append(f"- Comando: `{res.dry_run.command}`")
        lines.append(f"- Exit code: `{res.dry_run.returncode}`")
        lines.append(f"- Stdout: `{summarize_output(res.dry_run.stdout)}`")
        lines.append(f"- Stderr: `{summarize_output(res.dry_run.stderr)}`")
        lines.append("")
        lines.append("#### Real-run")
        lines.append(f"- Estado: `{res.real_run.status}`")
        lines.append(f"- Motivo: {res.real_run.reason}")
        lines.append(f"- Comando: `{res.real_run.command}`")
        lines.append(f"- Exit code: `{res.real_run.returncode}`")
        lines.append(f"- Stdout: `{summarize_output(res.real_run.stdout)}`")
        lines.append(f"- Stderr: `{summarize_output(res.real_run.stderr)}`")
        lines.append("")
    REPORT_PATH.write_text("\n".join(lines), encoding="utf-8")


def write_report_json(results: list[CaseResult]) -> None:
    payload = {"results": [asdict(result) for result in results]}
    REPORT_JSON_PATH.write_text(json.dumps(payload, indent=2), encoding="utf-8")


def main() -> int:
    clean_and_prepare_dirs()
    case1 = build_case1()
    case2 = build_case2()
    case3 = build_case3()
    case4 = build_case4()
    results = [
        evaluate_case("Caso 1: archivo .txt -> .md", case1, verify_case1_real),
        evaluate_case("Caso 2: carpeta con 2 .txt -> .md", case2, verify_case2_real),
        evaluate_case("Caso 3: reorganizacion recursiva por tipo (JPG/TIFF/PDF)", case3, verify_case3_real),
        evaluate_case("Caso 4: combinacion de caso 2 y caso 3", case4, verify_case4_real),
    ]
    write_report_md(results)
    write_report_json(results)
    not_ok = [r for r in results if r.final_status != "PASS"]
    for res in results:
        print(f"[{res.final_status}] {res.case_name}")
    print(f"Reporte MD: {REPORT_PATH}")
    print(f"Reporte JSON: {REPORT_JSON_PATH}")
    if not not_ok:
        shutil.rmtree(E2E_GENERATED, ignore_errors=True)
        shutil.rmtree(E2E_CONFIG, ignore_errors=True)
        print("E2E OK: artefactos temporales limpiados.")
    else:
        print("E2E con fallos: artefactos conservados para diagnostico.")
    return 1 if not_ok else 0


if __name__ == "__main__":
    raise SystemExit(main())
