# Otter

![Otter](website/home/static/img/otter.png)

![License](https://img.shields.io/badge/license-MIT-blue.svg) ![Build](https://travis-ci.org/vektorlab/otter.svg?branch=master) [![Docker Repository on Quay](https://quay.io/repository/vektorlab/otter/status "Docker Repository on Quay")](https://quay.io/repository/vektorlab/otter)

Otter is a high-performance opinionated configuration management framework written in Go for servers that run containers.

The data center is moving towards a container-centric world where the role of the host operating system becomes far 
more static and disposable. Modern configuration management tools (CM) solve some problems but create many more. 
Existing CM solutions with their complex abstractions and domain specific languages are built for servers with life 
spans of several months or years. In a container-centric data center the life span of a host operating system should 
last only as long as it can provide fundamental services without interruption (e.g. docker daemon, networking services, etc).

## Configuration Management
Otter draws influence from many successful CM systems but drastically simplifies them for the container world. 
Otter follows the "batteries-included" model to provide configuration support for major container orchestration 
tools setup following best-practices.

Some features include:

* Declarative state
* File system operations
* Template rendering
* Package management
* Support for installing and configuring major Docker orchestration systems

##Roadmap

### Bootstrapping
Handles the initial bootstrapping of a host OS taking it from a new OS installation to being ready for handling 
container workload from a container orchestration tool.

### Remote Execution
Largely concurrent remote execution.

## Disclaimer
**Otter is under active development and is not production-ready**

