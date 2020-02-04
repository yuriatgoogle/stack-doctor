#Copyright 2019 Google LLC
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#https://www.apache.org/licenses/LICENSE-2.0
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.

# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.12 base image
FROM golang:1.12

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/yuriatgoogle/stack-doctor/opentelemetry-traces-logs

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Install the packages
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["go", "run", "main.go"]