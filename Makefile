.PHONY: build
build: ## bash hack/build.sh
	bash hack/build.sh

.PHONY: test
test: ## bash hack/test.sh
	bash hack/all_test.sh

help: ## help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort