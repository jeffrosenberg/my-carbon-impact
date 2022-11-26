# My Carbon Impact

An app to simplify tracking your carbon footprint over time

## Design decisions and apologies

This project is meant to be a learning experience for me, more than anything functional.
Although I've chosen a subject area that interests me and which could have a real use,
it will most likely just serve as a platform for me to noodle on occasionally, as well
as a template for other projects I might want to start using a similar stack.

As a result, the design decisions made in this project are based more around the things
I'd like to practice than they are around any sort of technical optimization.
Some of the topics I intend to explore with this project include:

- NoSQL single-table design
  - It's highly likely that this schema would be better suited to a standard RDBMS
- Logging and metrics in AWS
- Managing deployments with Terraform
- Working with channels
  - May be a while before I get to this one

## Development notes

> These are just meant to be notes for myself, since I'm developing this sporadically.
> TODO: Remove if this project ever becomes more fully fleshed-out

### AWS / Terraform deployments

Resources are deployed to aws-west-2 (Oregon)

### Codebase status / progress

Resources existing so far:

- API Lambda (`cmd/api`):
  - Profile CRUD lambdas
    - Create and Get are fleshed out with persistence, others are stubs
  - "Scaffold" CRUD stubs
    - I don't even remember what these are, I think just meant to be a placeholder/template for future entities
- Web Lambda (`cmd/web`):
  - Profile page lambdas
    - Create and Get have very basic functionality but need more work
- Enabling code (`pkg` / `internal`):
  - `Profile` struct and initial methods: `go/internal/profile/profile.go`
  - `Vehicle` struct and initial methods: `go/internal/profile/vehicle.go`
    - I jumped the gun on creating this, I'm not using it yet
  - DynamoDb persistence: `go/internal/db`
    - Create and Get defined, no scan/query/update
  - EPA constants: `go/pkg/constants/epa_constants.go`
  - Logging: `go/pkg/logging/log.go`

### TODO

> Use this section for notes on next steps

2022-11-27: Refactor Lambdas in some way to use less boilerplate?

### Testing / executing

At present, the only code that can be executed is the `profile-*` lambdas.
Running manually in AWS Console, they will emit a success status code.

### Roadmap / development plan

This is just a high level list that itself could use refinement!

- Build "horizontally" (one piece of functionality at a time, across the stack)
- Start by building out the full stack:
  - Develop API Gateway endpoints in parallel with the lambdas
  - Persist data (DynamoDb single-table design)
  - Add logging / metrics / alerting / etc.
  - Build a simple UI
    - Server-rendered using Go templates
    - Something like Alpine/HTMX for smoother transitions?
- Create end-to-end tests and build pipeline
- Add functionality
  - Flesh out Profile entity and functionality
  - Create mileage event and tests