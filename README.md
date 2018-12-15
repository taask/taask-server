![Logo](https://user-images.githubusercontent.com/5942370/50045447-d14f3400-0060-11e9-8e98-78cfdcf85a75.png)

# Taask Core

:wave: Welcome!

Taask Core is an open source system for running arbitrary jobs on any infrastructure. In other words, cloud native tasks-as-a-service.

## What does Taask Core do?
- :white_check_mark: Runs arbitrary workloads (tasks) written in any language
- :white_check_mark: Runs those tasks in a fault tolerant way
- :white_check_mark: Operates with no single point of failure
- :white_check_mark: Operates on any infrastructure
- :white_check_mark: Keeps task data encrypted, end to end

## What doesn't it do?
- :no_entry_sign: Act as a message bus
- :no_entry_sign: Replace Kubernetes
- :no_entry_sign: Orchestrate servers

## Project status
:warning: Taask is in *Alpha* (v0.0.4) and should not be used for critical workloads

## Components
Taask Core is comprised of three components:
- taask-server (this project)
- runners
- clients

And optional components:
- Prometheus (for metrics collection)
- Grafana (for metrics visualization)
- EFK (Elasticsearch, Fluentbit, Kibana - for aggregated logging)

### Layout
```
|                   |       Prometheus        |                         |
|                   |         Grafana         |                         |
|                   |    (collect metrics) ---|-----------|             |
|                   |           |             |           |             |
|                   |           |             |           |             |
|     Clients  <----|------> Server <---------|---->    Runners         |
| (produce tasks)   |    (control layer)      |     (compute layer)     |
| (consume results) |  [multiple partners]    |  [scales automatically] |
|                   |                         |                         |
|                   |                         |                         |
|                   |                         |                         |
|                   |                         |                         |
```

### Server
The Taask Core server (this project) is the main component of the Taask control plane.
It is responsible for consuming tasks, tracking their state, and scheduling them onto the compute layer.
taask-server operates in a "managing partner" cluster. Multiple instances work together to share management tasks.
Clients can communicate with any instance of taask-server, and the tasks they produce will be distributed across the partners and their runners.
Runners register themselves with individual instances of taask-server, recieve tasks to be executed, and return their results.

### Runners
Taask Runners register themselves with taask-core to make themselves available for tasks.
Runners communicate with taask-server using gRPC, and bi-directionally stream data for optimal performance
Tasks are scheduled to runners, they are executed, and the results returned.
Runners can be written in any language using first-party and third-party runner libraries.

  - [runner-k8s](https://github.com/taask/runner-k8s): Runs tasks as Kubernetes Jobs, using any container image
  - [runner-golang]((https://github.com/taask/runner-golang): Go library for developing custom runners

### Child Runners
Runners such as runner-k8s can delegate tasks to _child runners_, which are short-lived, ephemeral runners meant to execute one task in their lifetime.

### Clients
Clients generate tasks, which are arbitrary JSON, and send them to taask-server to be processed.
Clients can stream task status updates from taask-server.
Clients encrypt task data by default, ensuring tasks are never transmitted or stored in a decrypted state.
