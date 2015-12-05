# Otter

Otter is an opinionated configuration management framework for servers that run containers.

The data center is moving towards a container-centric world where the role of the host operating system becomes far 
more static and disposable. Modern configuration management tools (CM) solve some problems but create many more. 
Existing CM solutions with their complex abstractions and domain specific languages are built for servers with life 
spans of several months or years. In a container-centric data center the life span of a host operating system should 
last only as long as it can provide fundamental services without interruption (e.g. docker daemon, networking services, etc).

Otter only considers a host server to be in one of two categories: "master-node", "compute-node"

#### Master-Node
A master server is a server runs cluster-management tools e.g. Kubernetes master, mesos-master, etc.

#### Compute-Node
A compute server is a server that runs container workloads as well as any necessary orchestration services e.g.
Kubernetes node, mesos-slave, etc.

## Configuration Management

No Formula, Playbooks, Recipes, or Manifests. Otter draws influence from many successful CM systems but drastically 
simplifies them for the container world. Otter follows the "batteries-included" model to provide configuration support
for major container orchestration tools setup following best-practices.

Some features include:

* Declarative state
* File system operations
* Template rendering
* Package management
* Support for installing and configuring major Docker orchestration systems

## Host Secret Management

Otter provides end-to-end encryption for exposing sensitive information to servers in a dynamic and programmable way.

## Bootstrapping

Handles the initial bootstrapping of a host OS taking it from a new OS installation to being ready for handling 
container workload from a container orchestration tool.

## Monitoring

Otter provides inventory and host OS monitoring

## Installation

    git clone git@github.com/vektorlab/otter.git
    cd otter
    go install
    otter -c=example/otter.yml ls |jq .
    otter -c=example/otter.yml state
    

