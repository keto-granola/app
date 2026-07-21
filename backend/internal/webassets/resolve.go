package webassets

type ResolvedAsset struct {
	Src string
	CSS []string
}

func ResolveEntry(loader *Loader, entryPath, clientURL string, devEnv bool) (*ResolvedAsset, error) {
	if devEnv {
		return &ResolvedAsset{Src: clientURL + "/" + entryPath}, nil
	}

	src, err := loader.Asset(entryPath)
	if err != nil {
		return nil, err
	}

	css, err := loader.AssetCSS(entryPath)
	if err != nil {
		return nil, err
	}

	return &ResolvedAsset{Src: src, CSS: css}, nil
}
