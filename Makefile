# Copyright 2019 The Kyma Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PKG := github.com/antoineco/kyma-event-sources


cmds := $(wildcard cmd/*)

.PHONY: all clean $(cmds)

all: clean $(cmds)

GOBUILD_FLAGS := -v
$(cmds): vendor
	@echo "+ Building $(notdir $@)"
	@CGO_ENABLED=0 go build -o $(notdir $@) \
		$(GOBUILD_FLAGS) \
		$(PKG)/$@

clean:
	@echo "+ Cleaning"
	rm -rf bin/
	rm -f $(notdir $(cmds))
	@go clean -x -i $(PKG)/cmd/...

vendor: Gopkg.lock
	@echo '+ Pulling vendored dependencies'
	@dep ensure -v --vendor-only


.PHONY: test

testpkgs := $(shell go list ./... | grep -vE '/apis/|/client/|/cmd/|/testing')

GOTEST_FLAGS := -cover
test: vendor
	@echo "+ Performing unit tests"
	@go test \
		$(GOTEST_FLAGS) \
		$(testpkgs)


#####################
#                   #
#  Code generation  #
#                   #
#####################

# see https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md#generate-code

# e.g. "sources/v1alpha1 sources/v1alpha2"
API_GROUPS := sources/v1alpha1
# generates e.g. "PKG/apis/sources/v1alpha1 PKG/apis/sources/v1alpha2"
api-import-paths := $(foreach group,$(API_GROUPS),$(PKG)/apis/$(group))

generators := deepcopy client lister informer injection
generators_bin := $(foreach x,$(generators),bin/$(x)-gen)

.PHONY: codegen $(generators) $(generators_bin)

codegen: $(generators) injection

# http://blog.jgc.org/2007/06/escaping-comma-and-space-in-gnu-make.html
comma := ,
space :=
space +=

# doc: https://godoc.org/k8s.io/code-generator/cmd/deepcopy-gen
deepcopy: bin/deepcopy-gen
	@echo "+ Generating deepcopy funcs for $(API_GROUPS)"
	@bin/deepcopy-gen \
		--go-header-file hack/boilerplate.go.txt \
		--input-dirs $(subst $(space),$(comma),$(api-import-paths))

client: bin/client-gen
	@echo "+ Generating clientsets for $(API_GROUPS)"
	@rm -rf client/generated/clientset
	@bin/client-gen \
		--go-header-file hack/boilerplate.go.txt \
		--input $(subst $(space),$(comma),$(API_GROUPS)) \
		--input-base $(PKG)/apis \
		--clientset-path $(PKG)/client/generated/clientset

lister: bin/lister-gen
	@echo "+ Generating listers for $(API_GROUPS)"
	@rm -rf client/generated/lister
	@bin/lister-gen \
		--go-header-file hack/boilerplate.go.txt \
		--input-dirs $(subst $(space),$(comma),$(api-import-paths)) \
		--output-package $(PKG)/client/generated/lister

informer: bin/informer-gen
	@echo "+ Generating informers for $(API_GROUPS)"
	@rm -rf client/generated/informer
	@bin/informer-gen \
		--go-header-file hack/boilerplate.go.txt \
		--input-dirs $(subst $(space),$(comma),$(api-import-paths)) \
		--output-package $(PKG)/client/generated/informer \
		--versioned-clientset-package $(PKG)/client/generated/clientset/internalclientset \
		--listers-package $(PKG)/client/generated/lister

injection: bin/injection-gen
	@echo "+ Generating injection for $(API_GROUPS)"
	@rm -rf client/generated/injection
	@bin/injection-gen \
		--go-header-file hack/boilerplate.go.txt \
		--input-dirs $(subst $(space),$(comma),$(api-import-paths)) \
		--output-package $(PKG)/client/generated/injection \
		--versioned-clientset-package $(PKG)/client/generated/clientset/internalclientset \
		--external-versions-informers-package $(PKG)/client/generated/informer/externalversions

$(generators_bin): vendor
	@if [ -d ./vendor/k8s.io/code-generator/cmd/$(notdir $@) ]; then \
		go build -o bin/$(notdir $@) ./vendor/k8s.io/code-generator/cmd/$(notdir $@); \
	else \
		go build -o bin/$(notdir $@) ./vendor/knative.dev/pkg/codegen/cmd/$(notdir $@); \
	fi
