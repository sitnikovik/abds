PR_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))

ifeq ($(firstword $(MAKECMDGOALS)),pr)
ifneq ($(words $(MAKECMDGOALS)),1)
$(foreach arg,$(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS)),$(eval $(arg):;@:))
endif
endif

.PHONY: pr
pr:
	@./scripts/pr/create.sh $(PR_ARGS)