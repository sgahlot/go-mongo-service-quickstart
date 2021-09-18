# go-service-quickstart
Sample microservice code, modeled after Quarkus quickstart apps, for MongoDB. It allows 
insertion, retrieval and deletion (TBD) operations on Fruit list.

_This app is using `go-kit` and `gorilla` libraries to provide the REST endpoints._

## Run the app:
### Local
To run the codebase locally against MongoDB, execute following commands:

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

## Tools used to perform API calls
You can either use `curl` or `httpie` to invoke the API from command line.
UI support is still a TODO. Examples provide use `httpie` tool

## Operations supported:
### Insert fruit

`http POST localhost:9090/api/v1/fruits name="SOME FRUIT NAME" description="SOME DESCRIPTION"`

_above will gets translated to a POST call using name and description as JSON payload_

### Retrieve all fruits

`http http://localhost:9090/api/v1/fruits\?name\="ALL"`

_above will retrieve all the fruits from database_
