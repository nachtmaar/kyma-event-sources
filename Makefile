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

.PHONY: all clean

bin := mqttsource-controller

all: clean $(bin)

$(bin):
	@echo "+ Building $@"
	@CGO_ENABLED=0 go build -o $(bin) \
		$(GOBUILD_FLAGS) \
		./

clean:
	@echo "+ Cleaning"
	rm -f $(bin)
	@go clean -x -i ./

vendor: Gopkg.lock
	@echo '+ Pulling vendored dependencies'
	@dep ensure -v --vendor-only
