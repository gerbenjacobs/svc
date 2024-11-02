# svc

An opinionated Go service and folder structure setup, inspired by @prep.

Note: This is **not** an actual framework; however you are free to use, adapt and learn (parts of) it.

## Rationale

There are 4 main directories of which 3 are synonymous to a layer:

- `cmd`: contains executables, most likely just an application (http, grpc), but also CLI tools
- `handler`: contains all files related to handling incoming requests or triggers
- `services`: contains your business logic
- `storage`: contains all repository and storage related operations

Domain models live in the root of the project, although people have created a `models` folder
before.

The last 3 folders each contain a self-titled go file that is the entry point of their layer.
For example `handler/handler.go` has the actual endpoints and the mux router,
`services/services.go` and `storage/storage.go` contain the interfaces that are being used
throughout the application.

Ideally *handlers* only know about *services* and services only know about *storage*.
Use the interfaces instead of actual implementations.

Every layers knows about the *Domain Models* and you should have these be the types that are
transferred between the Service and Storage layer. It's perfectly ok to then create custom
storage DAOs and custom input or output models for dealing with HTTP (i.e. an almost exact
copy of app.User but without the Password field). Just make sure that the service and storage
methods do the transformation back and forth.

Separation of concern and an explicit clarity is what this structure gives you.

*Throughout the code I've written comments prefixed with `Rationale:` to explain a bit about the code.*

## Running this example

- Copy the config.yml.example to config.yml
  - `cp config.yml.example config.yml`
- Create a secret token
  - `echo -n "my secret, make sure it's not too short" | openssl dgst -sha256`
- Set up the database with Docker Compose
  - `docker-compose up -d`
- Run the application (make sure the database container is up and running)
  - `go run cmd/app/main.go`

## Examples

You can use the User service to test some endpoints out.
Change the token when retrieving your user to the one you received when creating a user.

### Creating a user

```shell script
curl -i -X POST -d '{"name": "Gerben"}' http://127.0.0.1:8000/v1/user
```

```json
{
 "id": "93d5baed-3f2d-44a6-b3e0-02841fcf30b9",
 "name": "Gerben",
 "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI5M2Q1YmFlZC0zZjJkLTQ0YTYtYjNlMC0wMjg0MWZjZjMwYjkiLCJuYmYiOjE1NzY2NjU5ODd9.UtIQNBpLBCRVH65LriP9uqKds-jrKJzmIcILQv082yc",
 "created_at": "2019-12-18T10:46:27.002421Z",
 "updated_at": "2019-12-18T10:46:27.002421Z"
}
```

### Retrieving your user

```shell script
curl -i -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI5M2Q1YmFlZC0zZjJkLTQ0YTYtYjNlMC0wMjg0MWZjZjMwYjkiLCJuYmYiOjE1NzY2NjU5ODd9.UtIQNBpLBCRVH65LriP9uqKds-jrKJzmIcILQv082yc" http://127.0.0.1:8000/v1/user
```

```json
{
 "id": "93d5baed-3f2d-44a6-b3e0-02841fcf30b9",
 "name": "Gerben",
 "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI5M2Q1YmFlZC0zZjJkLTQ0YTYtYjNlMC0wMjg0MWZjZjMwYjkiLCJuYmYiOjE1NzY2NjU5ODd9.UtIQNBpLBCRVH65LriP9uqKds-jrKJzmIcILQv082yc",
 "created_at": "2019-12-18T10:46:27Z",
 "updated_at": "2019-12-18T10:46:27Z"
}
```

### Converting our users to a JSON output file using the CLI

```shell
go run cmd/cli/main.go
```

```shell
2019/12/18 12:58:23 Found 1 users
2019/12/18 12:58:23 [93d5baed-3f2d-44a6-b3e0-02841fcf30b9] Gerben - 2019-12-18T10:46:27Z
2019/12/18 12:58:23 Finished writing to output.json
```
