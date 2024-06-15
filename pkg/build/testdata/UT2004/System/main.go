package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/aldehir/ut2u/pkg/ini"
)

var command string
var iniFile string

func main() {
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	if command == "" {
		fmt.Println("Unknown command")
		os.Exit(1)
	}

	for _, arg := range os.Args[2:] {
		switch {
		case strings.HasPrefix(arg, "-ini="):
			_, iniFile, _ = strings.Cut(arg, "=")
		}
	}

	if !strings.EqualFold(command, "make") {
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}

	fmt.Printf("Processing %s\n", iniFile)
	cfg, err := ParseConfig(iniFile)
	if err != nil {
		fmt.Printf("Failed to parse %s: %s\n", iniFile, err)
		os.Exit(2)
	}

	packages, found := cfg.Values("Editor.EditorEngine", "EditPackages")
	if !found {
		fmt.Println("No EditPacakges found")
		os.Exit(3)
	}

	var buildPackages []string
	for _, pkg := range packages {
		if !slices.Contains(defaultPackages, pkg) {
			buildPackages = append(buildPackages, pkg)
		}
	}

	if len(buildPackages) == 0 {
		fmt.Println("No custom build packages defined")
		os.Exit(4)
	}

	for _, pkg := range buildPackages {
		err = BuildPackage(pkg)
		if err != nil {
			fmt.Printf("Failed to build %s: %s\n", pkg, err)
		}
	}
}

func ParseConfig(path string) (*ini.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ini.Parse(f)
}

func BuildPackage(name string) error {
	for _, ext := range makeExtensions {
		fn := name + ext
		f, err := os.Create(fn)
		if err != nil {
			return err
		}
		defer f.Close()
		f.WriteString("Dummy file")
		fmt.Printf("=> %s\n", fn)
	}
	return nil
}

var makeExtensions = []string{".u", ".ucl"}

var defaultPackages = []string{
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
