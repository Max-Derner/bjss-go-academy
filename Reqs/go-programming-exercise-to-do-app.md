# Go Programming Exercise - To-Do App

The goal of the To Do Application is to build and evolve a production-lite service over four phases.

1. A command line application using an in-memory data store.
2. Introduce a REST API to wrap the data store and use JSON file(s) to provide data store persistence.
3. Add a Web App and make the To Do multi user
4. Use a DB to provide the data store persistence.

![To Do App High-Level Architecture](./go-programming-exercise-to-do-app.svg)

## Use of Go Packages

This program accomplies the Go Academy and therefore intention is to leverage the Go standard library.  The exception to this are the following packages:

* [github.com/google/uuid] - For working with UUID/GUID
* [github.com/google/go-cmp/cmp] - For comparing things useful for unit tests
* [github.com/lib/pg] - a pure Go PostgreSQL driver

## Development Approach

While developing the To-Do App use Git to store your solution and use Git Tags to mark final commit for each phase of the project.
As you progress through the project, make regular commits with a commit message documentating your progress.

## Phase Guidance

Each phase builds upon the previous phase and is epected to continue to work through all phases.  For instance in phase 1 you build a CLI application, this application with small changes, should continue to work through phase 2, 3 and 4.

### Phase 1

* CLI works directly with the In-Memory Data Store

### Phase 2

* Wrap the Data Store with the V1 REST API.
* Use the [context] package to add a TraceID and [slog] to enable traceability of calls through the solution.
* At the ToDo level, use CSP to support concurrent reads and concurrent safe write.
* Use Parallel tests to validate that the solution is concurrent safe.
* Update the CLI App to use the REST API.
* Add an JSON Data Store and use a startup value to tell the REST API which data store to use.

### Phase 3

* Add a V2 API to the REST API that supports multiple users
* Use Parallel test to validate the solution is concurrent safe across multiple users.
* Add a Web API that supports multiple users
* Add multi-user support to the CLI App

### Phase 4

* Add an additional data store implementation using an external DB (PostgreSQL)

[github.com/google/uuid]: https://pkg.go.dev/github.com/google/uuid
[github.com/google/go-cmp/cmp]: https://pkg.go.dev/github.com/google/go-cmp/cmp
[github.com/lib/pg]: https://pkg.go.dev/github.com/lib/pq 
[context]: https://pkg.go.dev/context
[slog]: https://pkg.go.dev/log/slog