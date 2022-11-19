# My Carbon Impact

An app to simplify tracking your carbon footprint over time

## Development notes

> These are just meant to be notes for myself, since I'm developing this sporadically.
> TODO: Remove if this project ever becomes more fully fleshed-out

### AWS / Terraform deployments

Resources are deployed to aws-west-2 (Oregon)

### Codebase status / progress

Resources existing so far:

- Lambda (`cmd` functions):
  - Profile CRUD stubs
  - "Scaffold" CRUD stubs
    - I don't even remember what these are, I think just meant to be a placeholder/template for future entities
- Enabling code (`pkg` / `internal` functions):
  - `Profile` struct and initial methods: `go/internal/profile/profile.go`
  - `Vehicle` struct and initial methods: `go/internal/profile/vehicle.go`
  - EPA constants: `go/pkg/epa/epa_constants.go`
  - Logging: `go/pkg/logging/log.go`

### Testing / executing

At present, the only code that can be executed is the `profile-*` lambdas.
Running manually in AWS Console, they will emit a success status code.

### TODO / roadmap

This is just a high level list that itself could use refinement!

- Develop API Gateway endpoints in parallel with the lambdas
- Flesh out Profile entity and functionality
- Create mileage event and tests
- Persist data (Cockroach? Dynamo?)
- Build a simple UI