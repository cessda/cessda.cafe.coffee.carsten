# Carsten's CESSDA CAFE Coffee Machine
# Copyright CESSDA-ERIC 2019
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not
# use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# build the binary
FROM golang:latest as builder
WORKDIR /go/src/coffee-api
COPY *.go /go/src/coffee-api/
RUN make prep
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build

# package the binary into a container
FROM scratch
COPY --from=builder /go/src/coffee-api/coffee-api /coffee-api
CMD ["/coffee-api"]

