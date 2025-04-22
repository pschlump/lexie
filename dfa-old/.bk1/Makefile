
all:
	go build

.json: input.json
	./input.pl <input.json >.json

coverage:
	go test -coverprofile=cover.out >,a
	sed -e 's,^_,_/../../../../..,' <cover.out >cover2.out
	go tool cover -html=cover2.out -o cover.html
	open cover.html

