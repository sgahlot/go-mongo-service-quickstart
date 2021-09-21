# go-service-quickstart
Sample microservice code, modeled after Quarkus quickstart apps, for MongoDB. It allows 
insertion, retrieval and deletion (TBD) operations on Fruit list.

_This app is using `go-kit` and `gorilla` libraries to provide the REST endpoints._

## Run the app
### Local

#### Run direct
`go run cmd/main.go`

#### Build executable
Use following command to build the executable:

* Mac

    `CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -o mongo-svc-quickstart ./cmd`

* Linux
 
    `CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o mongo-svc-quickstart ./cmd`

_Run the executable_

`./mongo-svc-quickstart`

#### Install MongoDB (in Docker) - Local ONLY
One can install MongoDB in Docker, to run the app locally. Execute following commands
to run MongoDB locally:

`docker compose -f resources/docker/mongo/mongo-docker-compose.yaml up -d`

_Above should fetch the latest image of MongoDB and run a container named mongo-test-db locally
The database should be initialized, but in case it is not, please run
following commands, in a terminal, to initialize the local Mongo DB:_

`docker exec -it mongo-test-db bash`

Once you're inside the mongo container, run following command:

`mongosh -u admin --authenticationDatabase admin`

For password, please provide the MONGO_INITDB_ROOT_PASSWORD value from `resources/docker/mongo/.env` file

Paste the contents of `resources/docker/mongo/mongo-init.sql` in the mongo prompt

To verify that the "sample-db" is created:

`show dbs`

To verify that the `fruit` collection is also created with 4 documents in it:

`db.fruit.find`

To exit the mongo shell as well as mongo docker container:
```
quit()
exit 
```

## Build Docker image for the service
* Build the image and run the container by running following command from root directory of the project:

  `docker compose -f go-svc-docker-compose.yaml up -d`

  _Above command will build the executable for the service as well as run the container
   for this image on port 9090_

* To stop and remove the container:

  `docker compose -f go-svc-docker-compose.yaml down`

**Alternate method to manually build and run the container:**

* Build the executable by following the [step in previous section](#build-executable).
* Run following command from root directory of the project to build docker image of the service:
  ```
  docker build -t go-mongo-quickstart:0.0.1-SNAPSHOT -f resources/docker/go/Dockerfile .
  ```
  _The image will be named `go-mongo-quickstart:0.0.1-SNAPSHOT`_
* Create and run a container using the image created in previous step:
  ```
  docker run --name go-mongo-fruit-app -d -p 9090:9090 --rm
      -e SERVICE_BINDING_ROOT=<BINDING_ROOT_DIR>
      -e DB_NAME=<DB_NAME_CONTAINING_FRUIT_COLLECTION>
      go-mongo-quickstart:0.0.1-SNAPSHOT
  ```
  _This will create a container named `go-mongo-fruit-app` listening on port 9090._  

## Tools used to perform API calls
You can either use `curl` or `httpie` to invoke the API from command line.
UI support is still a TODO. Examples provide use `httpie` tool

## Operations supported
### Insert fruit

`http POST localhost:9090/api/v1/fruits name="SOME FRUIT NAME" description="SOME DESCRIPTION"`

_above will gets translated to a POST call using name and description as JSON payload_

### Retrieve all fruits

`http http://localhost:9090/api/v1/fruits\?name\="ALL"`

_above will retrieve all the fruits from database_


## Environment variables used by the service (and binding client)
* `SERVICE_BINDING_ROOT`

  _Specfies the binding root containing a separate file for each value that's used by
   binding client to make the connection string_
* `DB_URL`

  _Database URL in case above property is NOT defined or want to run this service without binding client_
* `DB_NAME`

  _Name of the database from which the fruit collection is to be retrieved and used_
