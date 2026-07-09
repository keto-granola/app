package webassets

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed dist/.vite/manifest.json
var manifestFS embed.FS

//go:embed dist/assets
var assetsFS embed.FS

const islandsEntryDir = "src/islands/entries/"

func IslandEntryPath(name string) string {
	return islandsEntryDir + name + ".ts"
}

type ManifestEntry struct {
	File string   `json:"file"`
	CSS  []string `json:"css,omitempty"`
}

type Manifest map[string]ManifestEntry

type Loader struct {
	manifest Manifest
}

func New() (*Loader, error) {
	data, err := manifestFS.ReadFile("dist/.vite/manifest.json")
	if err != nil {
		return nil, fmt.Errorf("read vite manifest: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parse vite manifest: %w", err)
	}

	return &Loader{manifest: m}, nil
}

func (l *Loader) Asset(islandEntryPath string) (string, error) {
	entry, ok := l.manifest[islandEntryPath]
	if !ok {
		return "", fmt.Errorf("no manifest entry for %q", islandEntryPath)
	}

	return "/" + entry.File, nil
}

func (l *Loader) AssetCSS(islandEntryPath string) ([]string, error) {
	entry, ok := l.manifest[islandEntryPath]
	if !ok {
		return nil, fmt.Errorf("no manifest entry for %q", islandEntryPath)
	}

	out := make([]string, len(entry.CSS))
	for i, c := range entry.CSS {
		out[i] = "/" + c
	}

	return out, nil
}

func AssetsHandler() (http.Handler, error) {
	sub, err := fs.Sub(assetsFS, "dist/assets")
	if err != nil {
		return nil, err
	}

	return http.FileServer(http.FS(sub)), nil
}
