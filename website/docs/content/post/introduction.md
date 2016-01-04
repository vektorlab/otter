+++
description = "First post"
weight = 1
type = "post"
class="post first"
+++

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
