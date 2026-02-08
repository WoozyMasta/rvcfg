GO ?= go
LINTER  ?= golangci-lint
ALIGNER ?= betteralign
PROFILE_DIR ?= .profiles
PPROF_HTTP ?= :18080

.PHONY: test bench bench-parse verify vet fmt fmt-check lint align align-fix check tidy download tools release-notes diag-doc \
	profile-clean profile-parse-cpu profile-parse-mem profile-parse-all pprof-parse-cpu-web pprof-parse-mem-web pprof-parse-web

check: fmt-check vet lint align test

fmt:
	gofmt -w .

fmt-check:
	@gofmt -l . | tee /dev/stderr | read; \
	if [ $$? -eq 0 ]; then \
		echo "gofmt: files need formatting"; \
		exit 1; \
	fi

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

bench:
	$(GO) test -run=^$$ -bench 'Benchmark' -benchmem

bench-parse:
	$(GO) test -run=^$$ -bench 'BenchmarkParseVehicleConfig' -benchmem

verify:
	$(GO) mod verify

tidy:
	$(GO) mod tidy

download:
	$(GO) mod download

lint:
	$(LINTER) run ./...

align:
	$(ALIGNER) ./...

align-fix:
	$(ALIGNER) -apply ./...

tools:
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	$(GO) install github.com/dkorunic/betteralign/cmd/betteralign@latest
	$(GO) install github.com/google/pprof@latest

release-notes:
	@awk '\
	/^<!--/,/^-->/ { next } \
	/^## \[[0-9]+\.[0-9]+\.[0-9]+\]/ { if (found) exit; found=1; next } found { print } \
	' CHANGELOG.md

diag-doc:
	$(GO) run ./cmd/diagdocgen -out DIAGNOSTICS.md
