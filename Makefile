POSTPROCESS_MLX = close.mlx

.PHONY: all
all: m2e3dv6.view

m2e3dv6.view: output/m2e3dv6_body.ply output/m2e3dv6_pressuretongue.ply output/m2e3dv6_hardware.ply

output:
	mkdir output

output/%.stl: m2e3dv6 output
	cd output && ../m2e3dv6

m2e3dv6: $(wildcard *.go)
	go build

.PHONY: %.view
%.view: %.mlp
	meshlab $< 2>/dev/null

%.ply: %.stl $(POSTPROCESS_MLX)
	meshlab.meshlabserver -i $< -o $@ -s $(POSTPROCESS_MLX)

.PHONY: clean
clean:
	rm -rf output
	go clean
