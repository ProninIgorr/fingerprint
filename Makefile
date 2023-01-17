.PHONY: run_register build_register clean_register run_detector build_detector clean_detector

run_register: build_register
	scripts/register.sh

build_register: clean_register
	go build ./cmd/fgp_register

clean_register:
	rm -f fgp_register

run_detector: build_detector
	scripts/detect.sh

build_detector: clean_detector
	go build ./cmd/fgp_detect

clean_detector:
	rm -f fgp_detect