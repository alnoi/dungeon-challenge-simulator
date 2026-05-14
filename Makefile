.PHONY: fmt lint test run sample-data

APP := ./cmd/app
CONFIG ?= config.json
EVENTS ?= events

fmt:
	gofmt -w cmd internal tests

lint:
	golangci-lint run

test:
	go test ./...

run:
	go run $(APP) -config $(CONFIG) -events $(EVENTS)

sample-data:
	@printf '%s\n' \
		'{' \
		'    "Floors": 2,' \
		'    "Monsters": 2,' \
		'    "OpenAt": "14:05:00",' \
		'    "Duration": 2' \
		'}' > $(CONFIG)
	@printf '%s\n' \
		'[14:00:00] 1 1' \
		'[14:05:00] 1 2' \
		'[14:06:00] 1 3' \
		'[14:07:00] 1 3' \
		'[14:08:00] 1 4' \
		'[14:08:00] 1 6' \
		'[14:12:00] 1 7' \
		'[14:15:00] 1 8' > $(EVENTS)
	@echo "Created $(CONFIG) and $(EVENTS)"
