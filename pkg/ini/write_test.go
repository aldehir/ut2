package ini

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWrite(t *testing.T) {
	cfg := &Config{
		Sections: []*Section{
			{
				Name: "Engine.Engine",
				Items: []*Item{
					{"RenderDevice", []string{"D3DDrv.D3DRenderDevice"}},
				},
			},
			{
				Name: "Core.System",
				Items: []*Item{
					{
						Key: "Paths",
						Values: []string{
							"../System/*.u",
							"../Maps/*.ut2",
							"../Textures/*.utx",
							"../Sounds/*.uax",
						},
					},
				},
			},
		},
	}

	var buf strings.Builder

	err := cfg.Write(&buf)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(expectedWriteResult, buf.String()); diff != "" {
		t.Errorf("TestWrite() Mismatch (-want, +got):\n%s", diff)
	}
}

var expectedWriteResult = `[Engine.Engine]
RenderDevice=D3DDrv.D3DRenderDevice

[Core.System]
Paths=../System/*.u
Paths=../Maps/*.ut2
Paths=../Textures/*.utx
Paths=../Sounds/*.uax

`
