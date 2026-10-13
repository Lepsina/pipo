# Pinned versions of the external test-vector releases. Keep SPEC_VERSION in
# sync with internal/spectest.SpecVersion (the harness checks spectests/VERSION
# against it) and with the spec revision the code is written against.
SPEC_VERSION := v1.6.0-beta.0
BLS_TESTS_VERSION := v0.1.2

# Presets to download. minimal/ is not used by this project yet; add it with
# `make spectests SPEC_PRESETS="general mainnet minimal"` if that changes.
SPEC_PRESETS ?= general mainnet

SPEC_DIR := spectests
SPEC_URL := https://github.com/ethereum/consensus-spec-tests/releases/download/$(SPEC_VERSION)
BLS_TESTS_URL := https://github.com/ethereum/bls12-381-tests/releases/download/$(BLS_TESTS_VERSION)

# sha256 of the pinned release assets, so the vectors cannot be swapped
# silently. Update these whenever a version above changes.
SHA256_general := ad3de3429dbe65a0e23baea02108207adceaafe90679ec5a2ea1216f1e070cc6
SHA256_mainnet := e192eb2cd78eecd8a16784ee587e55e5f521bd6f58de600f0500c2aab7d2869b
SHA256_minimal := c28c30c32bddd0a24bb04fa80ceb4f92b859c899dd26927495b5d8b0b1fa5c1c
SHA256_bls_tests := a4ef909b6248395094949a3bcdeb4e1b5b999733d4436f2c5006edba99d29279

.PHONY: test spectests spectests-clean

test:
	go test ./...

# Download and extract the test vectors. Each archive is fetched once per
# version: a marker file records what is already unpacked.
spectests:
	@mkdir -p $(SPEC_DIR)
	@for preset in $(SPEC_PRESETS); do \
		marker=$(SPEC_DIR)/.$$preset-$(SPEC_VERSION); \
		if [ -f $$marker ]; then echo "$$preset $(SPEC_VERSION): already present"; continue; fi; \
		case $$preset in \
			general) want=$(SHA256_general) ;; \
			mainnet) want=$(SHA256_mainnet) ;; \
			minimal) want=$(SHA256_minimal) ;; \
			*) echo "unknown preset: $$preset" >&2; exit 1 ;; \
		esac; \
		tmp=$$(mktemp -d) || exit 1; \
		echo "downloading $$preset.tar.gz $(SPEC_VERSION)"; \
		curl -fL --progress-bar -o $$tmp/$$preset.tar.gz $(SPEC_URL)/$$preset.tar.gz || { rm -rf $$tmp; exit 1; }; \
		echo "$$want  $$tmp/$$preset.tar.gz" | sha256sum -c - || { rm -rf $$tmp; exit 1; }; \
		rm -rf $(SPEC_DIR)/$$preset; \
		tar -xzf $$tmp/$$preset.tar.gz -C $(SPEC_DIR) --strip-components=1 || { rm -rf $$tmp; exit 1; }; \
		rm -rf $$tmp; \
		touch $$marker; \
	done
	@marker=$(SPEC_DIR)/.bls-12-381-tests-$(BLS_TESTS_VERSION); \
	if [ -f $$marker ]; then \
		echo "bls-12-381-tests $(BLS_TESTS_VERSION): already present"; \
	else \
		tmp=$$(mktemp -d) || exit 1; \
		echo "downloading bls_tests_yaml.tar.gz $(BLS_TESTS_VERSION)"; \
		curl -fL --progress-bar -o $$tmp/bls_tests_yaml.tar.gz $(BLS_TESTS_URL)/bls_tests_yaml.tar.gz || { rm -rf $$tmp; exit 1; }; \
		echo "$(SHA256_bls_tests)  $$tmp/bls_tests_yaml.tar.gz" | sha256sum -c - || { rm -rf $$tmp; exit 1; }; \
		rm -rf $(SPEC_DIR)/bls-12-381-tests; \
		mkdir -p $(SPEC_DIR)/bls-12-381-tests; \
		tar -xzf $$tmp/bls_tests_yaml.tar.gz -C $(SPEC_DIR)/bls-12-381-tests || { rm -rf $$tmp; exit 1; }; \
		rm -rf $$tmp; \
		touch $$marker; \
	fi
	@printf 'consensus-spec-tests %s\nbls12-381-tests %s\n' "$(SPEC_VERSION)" "$(BLS_TESTS_VERSION)" > $(SPEC_DIR)/VERSION
	@echo "vectors ready in $(SPEC_DIR)/ ($$(cat $(SPEC_DIR)/VERSION | tr '\n' ' '))"

spectests-clean:
	rm -rf $(SPEC_DIR)
