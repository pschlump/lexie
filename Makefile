
all:
	go build

.json: input.json
	./input.pl <input.json >.json

all_dependencies:
	go list -m all >,a

graph_dependencies:
	go mod graph

