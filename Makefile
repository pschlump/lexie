
all:
	go build

.json: input.json
	./input.pl <input.json >.json

