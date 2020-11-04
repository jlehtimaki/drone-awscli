# drone-awscli
Drone plugin to execute AWS Cli commands.

## Build
Build the binary with the following commands:

```export GO111MODULE=on
go mod download
go test
go build
```

## Docker

Build the docker image with:
```
docker build --rm=true -t lehtux/drone-awscli .
```

## Usage
```
docker run --rm lehtux/drone-awscli -e AWS_ACCESS_KEY=.... -e AWS_SECRET_KEY=.... -e PLUGIN_COMMAND="aws sts get-caller-identity"
```

## Parameters
| Paramenter            | Description                   |Required|          Type|
| -------------         |:-------------:                |:-------------:|:-----:|
| AWS_ACCESS_KEY        | AWS Access key                | YES           | String|
| AWS_SECRET_KEY        | AWS Access key secret         | YES           | String|
| PLUGIN_ASSUME_ROLE    | AWS Assume role               | NO            | String|
| PLUGIN_COMMANDS       | Commands to be run            | YES           |[]String|
| PLUGIN_SHELL          | Run AWS Cli in bash shell     | NO            |true/false|