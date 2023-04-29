# Go env variables
export GOFLAGS		?= -mod=vendor
export GO111MODULE	?= on

# Service env variables
export METRICS_ADDR			?= :8081


# Build variables
.DEFAULT_GOAL		:= help
IMAGE_NAME			?= merkle-tree
APP_NAME			?= merkle-tree
COMMIT_ID			?= snapshot
BUILD_VERSION		?= 0.0.0-snapshot
DOCKER_REGISTRY		?= 127.0.0.1:5000

TEST_REPORT_DIR			:= target/test
COVERAGE_FILE_SUFFIX	:= coverage.txt

# List of tools that can be installed with go get
TOOLS_DIR		:= .tools/
GOTESTSUM		:= ${TOOLS_DIR}gotest.tools/gotestsum@v1.7.0

${GOTESTSUM}:
	$(eval TOOL=$(@:%=%))
	@echo Installing ${TOOL}...
	go install $(TOOL:${TOOLS_DIR}%=%)
	@mkdir -p $(dir ${TOOL})
	@cp ${GOBIN}/$(firstword $(subst @, ,$(notdir ${TOOL}))) ${TOOL}

COVER_PKGS			= $(subst ${SPACE},${COMMA},$(shell go list ./...))
UNIT_TEST_FLAGS		= -race -cover
UNIT_TEST_FLAGS		+= -coverprofile=${TEST_REPORT_DIR}/unit_test_${COVERAGE_FILE_SUFFIX}
INT_TEST_FLAGS		= -race -cover -coverpkg=${COVER_PKGS} -tags=integration -run='^TestIntegration'
INT_TEST_FLAGS		+= -coverprofile=${TEST_REPORT_DIR}/integration_test_${COVERAGE_FILE_SUFFIX}

.PHONY: help
help:  ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: run
run: ## Run service locally with default values
	go run cmd/main.go

builds/${APP_NAME}:
	env GOOS=linux CGO_ENABLED=0 go build -o build/_output/bin/${APP_NAME} cmd/main.go


.PHONY: build
build: builds/${APP_NAME}

.PHONY: test
test: ${GOTESTSUM} ## Run unit tests
	@mkdir -p ${TEST_REPORT_DIR}
	${GOTESTSUM} --jsonfile ${TEST_REPORT_DIR}/units-tests-output.log --junitfile=${TEST_REPORT_DIR}/unit-junit.xml -- ${UNIT_TEST_FLAGS} ./...


.PHONY: docker
docker: builds/${APP_NAME}
	docker build --rm -t ${APP_NAME} .

.PHONY: lint
lint:
	mkdir -p target
ifeq (, $(shell which golangci-lint))
	docker run --rm -v $(shell pwd):/work -w /work neo-docker-releases.repo.lab.pl.alcatel-lucent.com/go-tools:1.15.5-39 golangci-lint run --out-format checkstyle 2>&1 | tee target/lint-report.xml
else
	golangci-lint run --out-format checkstyle 2>&1 | tee target/lint-report.xml
endif


.PHONY: docker-build
docker-build: build ## Build service docker image
	docker build --rm -t ${APP_NAME} --file build/Dockerfile .
	docker tag ${APP_NAME} ${DOCKER_REGISTRY}/${APP_NAME}:latest
	docker tag ${APP_NAME} ${DOCKER_REGISTRY}/${APP_NAME}:${BUILD_VERSION}

.PHONY: docker-push
docker-push: ## Publish docker image
	docker push ${DOCKER_REGISTRY}/${APP_NAME}:latest
	docker push ${DOCKER_REGISTRY}/${APP_NAME}:${BUILD_VERSION}

.PHONY: vendor
vendor: ## Update vendor folder to match go.mod
	go mod tidy
	go mod vendor
