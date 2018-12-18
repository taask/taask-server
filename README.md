![Logo](https://user-images.githubusercontent.com/5942370/50045447-d14f3400-0060-11e9-8e98-78cfdcf85a75.png)

# Taask Core

:wave: Welcome!

Taask Core is an open source system for running arbitrary jobs on any infrastructure. In other words, `cloud native tasks-as-a-service`.

## What does Taask Core do?
- :white_check_mark: Runs arbitrary workloads (tasks) written in any language
- :white_check_mark: Runs those tasks in a fault tolerant way
- :white_check_mark: Operates with no single point of failure
- :white_check_mark: Operates on any infrastructure
- :white_check_mark: Keeps task data encrypted, in transit and at rest

## What doesn't it do?
- :no_entry_sign: Act as a message bus
- :no_entry_sign: Replace Kubernetes
- :no_entry_sign: Orchestrate servers

## But is it serverless?
That depends on your definition, but in general, no. If you deploy the control plane and autoscale it, and use a runner that delegates to containers (such as runner-k8s), then it resembles serverless since it scales to 0 and allows for usage-based compute.

More precisely, Taask Core is a Functions-as-a-service platform that is tuned towards heavy, long-running workloads. Taask does stand for "Tasks as a service... k", after all :smile:

## Project status
:warning: Taask is in *pre-Alpha* and should not be used for critical workloads. When all critical components have been implemented, and the platform's security has been fully reviewed, it will graduate to alpha.

Taask Core has three goals:
- Security
- Reliability
- Speed

Taask Core will not graduate to alpha until its security has been proven. It will not graduate to beta until its reliability has been proven. It will not graduate to stable until it has been further optimized for production traffic without compromising security or reliability. This is the guiding principle the project follows.

## Components
Taask Core is comprised of three components:
- taask-server (this project)
- runners
- clients

And optional components:
- Service mesh (linkerd by default - for obervability and transport security)
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
| (produce tasks)   |    (control plane)      |     (compute plane)     |
| (consume results) |  [multiple partners]    |  [scales automatically] |
|                   |                         |                         |
|                   |                         |                         |
|                   |                         |                         |
|                   |                         |                         |
```

### Server
The Taask Core server (this project) is the main component of the Taask control plane.
It is responsible for consuming tasks, tracking their state, and scheduling them onto the compute plane.
taask-server operates in a "managing partner" cluster. Multiple instances work together to share management tasks and remain fault-tolerant.
Clients can communicate with any instance of taask-server, and the tasks they produce will be distributed across the partners and their runners.
Runners register themselves with individual instances of taask-server.

### Runners
The Taask compute plane is comprised of one or more runners. Runners register themselves with taask-server to make themselves available for tasks.
Runners communicate with taask-server using gRPC, and bi-directionally stream data for optimal performance.
Tasks are scheduled to runners, they are executed, and the results returned.
Runners can be written in any language using first-party and third-party runner libraries.

	- [runner-k8s](https://github.com/taask/runner-k8s): Runs tasks as Kubernetes Jobs using container images
	- [runner-golang](https://github.com/taask/runner-golang): Go library for developing custom runners

### Child Runners
Runners such as runner-k8s can delegate tasks to _child runners_, which are short-lived, ephemeral runners meant to execute one task in their lifetime.

### Clients
Clients generate tasks, which are arbitrary JSON, and send them to taask-server to be executed.
Clients can stream task status updates from taask-server.
Clients encrypt task data before submitting them, ensuring tasks are never transmitted or stored in a decrypted state.

	- [client-golang](https://github.com/taask/client-golang): Go library for producing tasks and communicating with taask-server

### Tasks
Tasks are arbitrary JSON. Clients produce tasks, runners use the task JSON as input for their `Run` function, and return arbitrary result JSON.

### Security
Task data is opaque and encrypted in transit and at rest. Security is the top priority of Taask Core.
- Client libraries generate a _client keypair_, used to authenticate with taask-server.
- Clients encrypt task data with a symmetric _task key_, and then encrypt the _task key_ with a _task keypair_.
	- The _task key_ and _task keypair_ are cached in the memory of the client who created the task. Eventually, the _task keypair_ will be optionally stored in a shared location, accessible to any clients with access to the storage. This will allow clients to become more ephemeral.
- The _task key_ is also encrypted with the _master runner public key_. Both encrypted versions of the _task key_ are sent with the task.
- Task data remains encrypted until it is recieved by the runner, at which point it is decrypted and provided to the `Run` function of the runner. 
  - If the assigned runner delegates the task, the task JSON is not decrypted until it reaches its child runner.
- taask-server owns the _master runner private key_, and can decrypt the _task key_.
- Every runner and child runner generates a _runner keypair_ upon startup, used to authenticate with taask-server and decrypt task data.
- When a task is assigned to a runner, taask-server re-encrypts the _task key_ using the runner's _runner public key_
- The runner decrypts the _task key_ and then uses it to decrypt the task data.
- The task key is then used to encrypt the task result data, and the result is returned to the client.

### Maintainers
- Connor Hicks [@cohix](https://github.com/cohix)