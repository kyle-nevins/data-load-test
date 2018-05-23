# Makefile for data-load-test

# Include the Makefile PCF construct
include Makefile-pcf.mk

.PHONY: all

SHELL = bash

app:	ensure-deps	cf-org-space
	@$(CF_CMD) push
	@$(CF_CMD) bind-service crunchy-data-loader $(DB_SERVICE_INSTANCE)

check: $(GO_CMD) $(DEP_CMD)

# Check the DEP command
$(DEP_CMD):
	$(error Unable to find dep binary at DEP_CMD: $(DEP_CMD))
	false

ensure-deps: check
	$(DEP_CMD) ensure

# Check the GO command
$(GO_CMD):
	$(error Unable to find Go binary at GO_CMD: $(GO_CMD))
	false

cf:
	@$(CF_CMD) dev status

cf-org-space:
	$(info INFO: Using $(CF_CMD), assuming CF DEV is configured)
	@$(CF_CMD) login -a $(CF_URL) --skip-ssl-validation -u $(CF_USERNAME) -p $(CF_PASSWORD) -o $(CF_ORG) -s $(CF_SPACE)
	@$(CF_CMD) service-access
