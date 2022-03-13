AVALANCHE_VERSION ?= v1.7.4
GOLANG_VERSION ?= 1.17.1
VMID ?= sqja3uK17MJxfC7AN8nGadBw9JK5BcrsNwNynsqP5Gih8M5Bm
BUILD_OUTPUT_FOLDER ?= ./build

# BUILD_OUTPUT_FOLDER should stay exported because we might want to change the location of resulting binary
build-plugin:
	./scripts/build.sh $(BUILD_OUTPUT_FOLDER)/$(VMID)

# Desired AVALANCHE_VERSION
avalanche-version:
	@echo $(AVALANCHE_VERSION)

# Desired GOLANG_VERSION
golang-version:
	@echo $(GOLANG_VERSION)
