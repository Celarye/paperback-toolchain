package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"

	assets "paperback.moe/paperback-cli/assets"
)

type BundlerConfiguration struct {
	OutputFolder string
}

type bundleContext struct {
	baseDir string
	config  *BundlerConfiguration
}

func (ctx bundleContext) srcDir() string {
	return filepath.Join(ctx.baseDir + "/src")
}

func (ctx bundleContext) bundlesDir() string {
	return filepath.Join(ctx.baseDir + "/bundles/" + ctx.config.OutputFolder)
}

func (ctx bundleContext) entryPoints() []api.EntryPoint {
	files, err := os.ReadDir(ctx.srcDir())
	if err != nil {
		panic(err)
	}

	entryPoints := []api.EntryPoint{}

	for _, item := range files {
		if strings.HasPrefix(item.Name(), ".") {
			continue
		}

		itemPath := filepath.Join(ctx.srcDir() + "/" + item.Name())

		pbconfigTSPath := filepath.Join(itemPath + "/pbconfig.ts")
		mainTSPath := filepath.Join(itemPath + "/main.ts")

		if _, err := os.Stat(pbconfigTSPath); err != nil {
			continue
		}

		if _, err := os.Stat(mainTSPath); err != nil {
			continue
		}

		entryPoints = append(entryPoints, api.EntryPoint{
			InputPath:  mainTSPath,
			OutputPath: filepath.Join(item.Name() + "/index"),
		})
	}

	return entryPoints
}

type repositoryVersionInfo struct {
	Toolchain string `json:"toolchain"`
	Types     string `json:"types"`
}

type repositoryInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type repository struct {
	BuiltWith  repositoryVersionInfo `json:"builtWith"`
	Repository repositoryInfo        `json:"repository"`
	BuildTime  uint64                `json:"buildTime"`

	Sources []any `json:"sources"`
}

func Bundle(baseDir string, config BundlerConfiguration) {
	ctx := &bundleContext{baseDir: baseDir, config: &config}

	log.Printf("bundler: base directory: %v", ctx.baseDir)
	mainEntryPoints := ctx.entryPoints()
	os.RemoveAll(ctx.bundlesDir())
	if err := os.Mkdir(ctx.bundlesDir(), os.ModePerm); err != nil {
		log.Fatalf("error writing shims.js: %v", err.Error())
	}

	shimFilePath := filepath.Join(ctx.bundlesDir(), "shims.js")
	if err := os.WriteFile(shimFilePath, assets.ShimsJS, 0644); err != nil {
		log.Fatalf("error writing shims.js: %v", err.Error())
	}

	// Transpile project
	mainBuildResult := api.Build(api.BuildOptions{
		Bundle:              true,
		EntryPointsAdvanced: mainEntryPoints,
		Format:              api.FormatIIFE,
		GlobalName:          "source",
		Metafile:            true,
		Outdir:              ctx.bundlesDir(),
		AbsWorkingDir:       ctx.baseDir,
		Write:               true,
		Inject:              []string{shimFilePath},
	})

	for _, err := range mainBuildResult.Errors {
		log.Printf(err.Text)
	}

	// Generate versioning info
	pbconfigEntryPoints := []api.EntryPoint{}

	for _, entryPoint := range mainEntryPoints {
		pbconfigEntryPoints = append(pbconfigEntryPoints, api.EntryPoint{
			InputPath: filepath.Join(filepath.Dir(entryPoint.InputPath), "pbconfig.ts"),
		})
	}

	now := time.Now()
	pbconfigBuildResult := api.Build(api.BuildOptions{
		Bundle:              true,
		EntryPointsAdvanced: pbconfigEntryPoints,
		Format:              api.FormatIIFE,
		AbsWorkingDir:       ctx.baseDir,
		Write:               false,
		TreeShaking:         api.TreeShakingTrue,
		Target:              api.ES2020,
		GlobalName:          "global",
	})
	log.Printf("bundler: project transpiled: %v", time.Since(now))

	sourceInfos := []any{}
	for index, outputFile := range pbconfigBuildResult.OutputFiles {
		sourceID := filepath.Dir(mainEntryPoints[index].OutputPath)
		now = time.Now()
		defer func() {
			log.Printf("bundler: generated info for %v: %v", sourceID, time.Since(now))
		}()

		jsctx := goja.New()
		_, err := jsctx.RunScript("", string(outputFile.Contents))
		if err != nil {
			fmt.Printf("1 err: %v\n", err)
			continue
		}

		defObject := jsctx.GlobalObject().Get("global").ToObject(jsctx).Get("default").ToObject(jsctx)

		defObject.Set("id", sourceID)

		jsonBytes, err := defObject.MarshalJSON()
		if err != nil {
			fmt.Printf("5 err: %v\n", err)
			continue
		}

		var sourceInfo any
		json.Unmarshal(jsonBytes, &sourceInfo)
		sourceInfos = append(sourceInfos, sourceInfo)

		infoPath := filepath.Join(ctx.bundlesDir(), sourceID, "info.json")
		staticsPath := filepath.Join(ctx.srcDir(), sourceID, "static")
		os.WriteFile(infoPath, jsonBytes, 0644)
		os.CopyFS(filepath.Join(ctx.bundlesDir(), sourceID, "static"), os.DirFS(staticsPath))
	}

	now = time.Now()
	repositoryInfo := repository{
		BuiltWith: repositoryVersionInfo{
			Toolchain: "",
			Types:     "",
		},

		Repository: repositoryInfo{
			Name:        "Paperback Extension Repository",
			Description: "",
		},

		Sources:   sourceInfos,
		BuildTime: uint64(time.Now().Unix()),
	}

	repositoryInfoBytes, err := json.Marshal(repositoryInfo)
	if err != nil {
		fmt.Printf("6 err: %v\n", err)
		return
	}

	infoPath := filepath.Join(ctx.bundlesDir(), "versioning.json")
	os.WriteFile(infoPath, repositoryInfoBytes, 0644)

	homepagePath := filepath.Join(ctx.bundlesDir(), "index.html")
	os.WriteFile(homepagePath, assets.HomepageHTML, 0644)

	log.Printf("bundler: generated versioning file: %v", time.Since(now))
}
