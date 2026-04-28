SHELL := /bin/sh

.PHONY: daemon-test daemon-run earth-backend-check link-check frontend-install frontend-build

daemon-test:
	cd edge_daemon && go test ./...

daemon-run:
	cd edge_daemon && go run ./cmd/agent

earth-backend-check:
	python -m compileall earth_mission_control/backend deep_space_link

link-check:
	python -m compileall deep_space_link

frontend-install:
	cd earth_mission_control/frontend && npm install

frontend-build:
	cd earth_mission_control/frontend && npm run build
