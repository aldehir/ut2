package build

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aldehir/ut2u/pkg/ini"
)

type Builder struct {
	// UT2004 root path
	RootPath string

	// UT2004 system path
	SystemDir string

	// Path to UT2004 UCC executable
	UCCPath string

	// Packages to build
	Packages []Package

	// Dependencies to include in EditPackages
	Dependencies []string
}

type Package struct {
	Name string
	Path string
}

var ErrInvalidRoot = errors.New("invalid root")
var ErrInvalidPackage = errors.New("invalid package")

func NewBuilder(root string) (*Builder, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("error resolving path %s: %w", root, err)
	}

	// Check if root path is correct by checking for System/UCC.exe
	ucc := filepath.Join(absRoot, "System/UCC.exe")
	_, err = os.Stat(ucc)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRoot, root)
	}

	return &Builder{
		RootPath:  absRoot,
		SystemDir: filepath.Join(absRoot, "System"),
		UCCPath:   ucc,
	}, nil
}

func (b *Builder) AddPackage(name string) error {
	path := filepath.Join(b.RootPath, name)

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error reading path %s: %w", path, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: %s", ErrInvalidPackage, path)
	}

	b.Packages = append(b.Packages, Package{
		Name: name,
		Path: path,
	})

	return nil
}

func (b *Builder) AddDependency(name string) error {
	path := filepath.Join(b.SystemDir, name+".u")
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return fmt.Errorf("dependency %s not found", err)
	}

	b.Dependencies = append(b.Dependencies, name)
	return nil
}

func (b *Builder) buildDir() (string, error) {
	dir := filepath.Join(b.RootPath, "Build")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (b *Builder) generateBuildIni(path string) error {
	editPkgs := make([]string, len(defaultEditPackages))
	copy(editPkgs, defaultEditPackages)

	for _, pkg := range b.Dependencies {
		editPkgs = append(editPkgs, pkg)
	}

	for _, pkg := range b.Packages {
		editPkgs = append(editPkgs, pkg.Name)
	}

	cfg := generateConfig(editPkgs)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error generating config: %w", err)
	}
	defer f.Close()

	err = cfg.Write(f)
	if err != nil {
		return fmt.Errorf("error writing conifg: %w", err)
	}

	return nil
}

func (b *Builder) Build(ctx context.Context) error {
	buildDir, err := b.buildDir()
	if err != nil {
		return err
	}

	ini := filepath.Join(buildDir, "build.ini")
	err = b.generateBuildIni(ini)
	if err != nil {
		return err
	}

	build := &build{
		builder:  b,
		buildDir: buildDir,
		iniFile:  ini,
	}

	err = build.build(ctx)
	if err != nil {
		return err
	}

	return nil
}

type build struct {
	builder *Builder

	buildDir string
	iniFile  string
}

func (b *build) clean() error {
	var removeFiles []string
	for _, pkg := range b.builder.Packages {
		removeFiles = append(
			removeFiles,
			filepath.Join(b.builder.SystemDir, pkg.Name+".u"),
			filepath.Join(b.builder.SystemDir, pkg.Name+".ucl"),
			filepath.Join(b.builder.SystemDir, pkg.Name+".int"),
		)
	}

	for _, fp := range removeFiles {
		err := os.Remove(fp)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		}
	}

	return nil
}

func (b *build) build(ctx context.Context) error {
	if err := b.clean(); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, b.builder.UCCPath, "make", fmt.Sprintf("-ini=%s", b.iniFile))
	cmd.Dir = b.builder.SystemDir // UCC requires cwd to be the system path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateConfig(editPkgs []string) *ini.Config {
	return &ini.Config{
		Sections: []*ini.Section{
			{
				Name: "Engine.Engine",
				Items: []*ini.Item{
					{
						Key:    "EditorEngine",
						Values: []string{"Editor.EditorEngine"},
					},
				},
			},
			{
				Name: "Core.System",
				Items: []*ini.Item{
					{
						Key:    "SavePath",
						Values: []string{"../Save"},
					},
					{
						Key:    "CachePath",
						Values: []string{"../Cache"},
					},
					{
						Key:    "CacheExt",
						Values: []string{".uxx"},
					},
					{
						Key:    "CacheRecordPath",
						Values: []string{"../System/*.ucl"},
					},
					{
						Key:    "MusicPath",
						Values: []string{"../Music"},
					},
					{
						Key:    "SpeechPath",
						Values: []string{"../Speech"},
					},
					{
						Key: "Paths",
						Values: []string{
							"../System/*.u",
							"../Maps/*.ut2",
							"../Textures/*.utx",
							"../Sounds/*.uax",
							"../Music/*.umx",
							"../StaticMeshes/*.usx",
							"../Animations/*.ukx",
							"../Saves/*.uvx",
						},
					},
				},
			},
			{
				Name: "Editor.EditorEngine",
				Items: []*ini.Item{
					{
						Key:    "CacheSizeMegs",
						Values: []string{"32"},
					},
					{
						Key:    "EditPackages",
						Values: editPkgs,
					},
				},
			},
		},
	}
}

var defaultEditPackages = []string{
	"Core",
	"Engine",
	"Fire",
	"Editor",
	"UnrealEd",
	"IpDrv",
	"UWeb",
	"GamePlay",
	"UnrealGame",
	"XGame_rc",
	"XEffects",
	"XWeapons_rc",
	"XPickups_rc",
	"XPickups",
	"XGame",
	"XWeapons",
	"XInterface",
	"XAdmin",
	"XWebAdmin",
	"Vehicles",
	"BonusPack",
	"SkaarjPack_rc",
	"SkaarjPack",
	"UTClassic",
	"UT2k4Assault",
	"Onslaught",
	"GUI2K4",
	"UT2k4AssaultFull",
	"OnslaughtFull",
	"xVoting",
	"StreamlineFX",
	"UTV2004c",
	"UTV2004s",
	"OnslaughtBP",
}
