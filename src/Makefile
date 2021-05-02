#!make

run:
	@go build -o=../bin/stackoverflow-questions .;
	../bin/stackoverflow-questions

test:
	go test -v ./...

GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

merge:
ifeq (,$(and $(filter Changes not staged for commit, $(shell git status)), $(filter Changes to be committed, $(shell git status))))
	- @echo "Changes ommitted/up-to-date for current working branch. Proceeding...";
ifdef to
ifeq (,$(filter $(to), $(GIT_BRANCH)))
ifneq (,$(filter $(to), $(shell git branch)))
	- @echo "Branch '${to}' found. Proceeding...";
	- $(eval CURRENT_BRANCH := $(GIT_BRANCH))
	- @git checkout ${to};
	- @git merge $(CURRENT_BRANCH);
	- @git checkout $(CURRENT_BRANCH);
	- @echo "All changes of '$(CURRENT_BRANCH)' merged with '$(to)'. Back to '$(CURRENT_BRANCH)'.";
else
	- @echo "Exited. Branch '${to}' not found.";
	- @exit 0;
endif
else
	- @echo "Exited. Current branch and merge-to branch cannot be same.";
	- @exit 0;
endif
else
	- @echo "Exited. Provide merge-to branch as to=<branch_name> and retry.";
	- @exit 0;
endif
else
	- @echo "Exited. Please do the add/rm/commit in current branch and retry.";
	- @exit 0;
endif

commit:
	- git add .
ifdef c
	- @echo ${c}
	- git commit -m "${c}"
else
	- git commit -m "Corrections"
endif
