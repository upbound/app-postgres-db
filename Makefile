# Project Setup
PROJECT_NAME := app-postgres-db
PROJECT_REPO := github.com/upbound/$(PROJECT_NAME)

PLATFORMS ?= linux_amd64
include build/makelib/common.mk

include build/makelib/output.mk

# Setup Go
GO_STATIC_PACKAGES = $(GO_PROJECT)/cmd/app-postgres-db
GO_SUBDIRS += cmd
GO111MODULE = on
include build/makelib/golang.mk

# Docker images
DOCKER_REGISTRY = upbound
IMAGES = app-postgres-db
include build/makelib/image.mk

# We want submodules to be set up the first time `make` is run.
# We manage the build/ folder and its Makefiles as a submodule.
# The first time `make` is run, the includes of build/*.mk files will
# all fail, and this target will be run. The next time, the default as defined
# by the includes will be run instead.
fallthrough: submodules
	@echo Initial setup complete. Running make again . . .
	@make

# Update the submodules, such as the common build scripts.
submodules:
	@git submodule sync
	@git submodule update --init --recursive
