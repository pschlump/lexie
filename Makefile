
all:
	go build

.json: input.json
	./input.pl <input.json >.json

all_dependencies:
	go list -m all >,a

graph_dependencies:
	go mod graph

go_test:
	( cd com ; go test )
	( cd in ; go test )
	( cd mt ; go test )
	( cd pbread ; go test )
	( cd re ; go test )
	( cd st ; go test )
	( cd tok ; go test )
	( cd nfa ; go test )
	( cd smap ; go test )
	( cd dfa ; go test )
	( cd eval ; go test )
	echo ""
	echo "PASS"
