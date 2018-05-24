
# Default PCF Dev Org and Space connection
CF_ORG ?= pcfdev-org
CF_SPACE ?= pcfdev-space
CF_USERNAME ?= admin
CF_PASSWORD ?= admin
CF_URL ?= https://api.local.pcfdev.io

# Default Service Instance name
DB_SERVICE_INSTANCE ?= myService

# Optionally include a user variables file
-include uservars.mk

# Bin Paths
CF_CMD ?= $(shell which cf)
GO_CMD ?= $(shell which go)
DEP_CMD ?= $(shell which dep)

# Useful for debugging
showvars:
	@$(info CF_BIN: $(CF_BIN))
	@$(info CF_ORG: $(CF_ORG))
	@$(info CF_PASSWORD: $(CF_PASSWORD))
	@$(info CF_SPACE: $(CF_SPACE))
	@$(info CF_USERNAME: $(CF_USERNAME))
	@$(info CF_URL: $(CF_URL))
