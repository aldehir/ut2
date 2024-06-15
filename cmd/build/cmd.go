package build

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/aldehir/ut2u/pkg/build"
)

var deps []string

var buildCommand = &cobra.Command{
	Use:                   "build [pkg-path...]",
	Short:                 "Build a UT2004 Mutator/Mod",
	RunE:                  doBuild,
	DisableFlagsInUseLine: true,
}

func EnrichCommand(cmd *cobra.Command) {
	cmd.AddCommand(buildCommand)
}

func init() {
	buildCommand.Flags().StringSliceVarP(&deps, "dep", "d", []string{}, "Dependency")
}

func doBuild(cmd *cobra.Command, args []string) error {
	var err error
	var packages []string

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		// If no arguments, assume cwd is the package path
		path, err := filepath.Abs(cwd)
		if err != nil {
			return fmt.Errorf("error resolving path %s: %w", cwd, err)
		}

		packages = append(packages, path)
	} else {
		for _, arg := range args {
			path, err := filepath.Abs(arg)
			if err != nil {
				return fmt.Errorf("error resolving path %s: %w", arg, err)
			}
			packages = append(packages, path)
		}
	}

	if len(packages) == 0 {
		return fmt.Errorf("no packages defined")
	}

	// Assume the root path is the parent of the package path
	root := filepath.Dir(packages[0])
	fmt.Printf("Root: %s\n", root)

	builder, err := build.NewBuilder(root)
	if err != nil {
		return err
	}

	// Add any explicit dependencies
	for _, dep := range deps {
		err := builder.AddDependency(dep)
		if err != nil {
			return err
		}
	}

	// Add package to build
	for _, pkg := range packages {
		name := filepath.Base(pkg)
		err := builder.AddPackage(name)
		if err != nil {
			return err
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	return builder.Build(ctx)
}
