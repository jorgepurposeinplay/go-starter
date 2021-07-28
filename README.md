# Go Starter
This starter kit is designed to get you up and running with a project structure optimized for developing RESTful API services in Go. It is an opinionated Go starter kit built on top of Chi, using battle tested libraries proven to provide a good
foundation for a project written in Golang.

# Background
Our team builds a lot of services in go, and as services started growing, a common problem started to appear, each project had a different backbone than the other. In order to increase productivity and consistency we decided on what in our opinion
were considered the best tools, and started this kit.

# External packages
The starter kit uses the following Go packages. They can be easily replaceable, their uses are highly localized and abstrated.
  * Routing: [go-chi](https://github.com/go-chi/chi)
  * Logging [zap](https://github.com/uber-go/zap)
  * ORM [gorm](https://github.com/go-gorm/gorm)
  * CLI [cobra](https://github.com/spf13/cobra)
  * Configuration [viper](https://github.com/spf13/viper)
  * Auto docker recompile [reflex](https://github.com/acim/go-reflex)

![Test](https://github.com/oakeshq/go-starter/workflows/Test/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/qreasio/go-starter-kit)](https://goreportcard.com/report/github.com/qreasio/go-starter-kit)

## Get Started
Easy way to start would be using docker-compose. The command bellow will start two containers, the Go service and the PostgreSQL container. With the use of reflex docker containers can be used for local development aswell.
```
docker-compose -f ./d8t/docker-compose.dev.yml up
```
