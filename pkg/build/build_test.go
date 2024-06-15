package build

import (
	"context"
	"testing"
)

//go:generate go build -o testdata/UT2004/System/UCC.exe testdata/UT2004/System/main.go

func TestBuild(t *testing.T) {
	builder, err := NewBuilder("testdata/UT2004")
	if err != nil {
		t.Error(err)
	}

	err = builder.AddPackage("ExamplePackage")
	if err != nil {
		t.Error(err)
	}

	err = builder.Build(context.Background())
	if err != nil {
		t.Errorf("Build failed: %s", err)
	}
}
