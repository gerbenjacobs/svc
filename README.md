# svc

An opinionated Go service and folder structure setup.

## Running

- Copy the config.yml.example to config.yml; `cp config.yml.example config.yml`
- Create a secret token; `echo -n "my secret, make sure it's not too short" | openssl dgst -sha256`
- Set up the database with Docker Compose; `docker-compose up -d`
- Run the application; `go run cmd/main.go`

## Examples

You can use the User service to test some endpoints out.

### Creating a user

```shell script
curl -v -X POST -d '{"name": "Gerben"}' http://127.0.0.1:8000/v1/user
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
curl -v -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI5M2Q1YmFlZC0zZjJkLTQ0YTYtYjNlMC0wMjg0MWZjZjMwYjkiLCJuYmYiOjE1NzY2NjU5ODd9.UtIQNBpLBCRVH65LriP9uqKds-jrKJzmIcILQv082yc" http://127.0.0.1:8000/v1/user
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