package config

import "path/filepath"

// Normaliza rutas linux windows convitiendolas a rutas limpias con filepath.Clean, esto es importante para evitar problemas con rutas como "workspace/inbox/../inbox/photos" o "workspace\\inbox\\..\\inbox\\photos"
func normalizePath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}
