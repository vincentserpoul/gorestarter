

# Golang starter for REST APIs

Go is simple, fast, lean, typed, compiled, opinionated... It was invented at Google to ease the work for large development groups, and I think it does the job well.

This starter is as simple, lean as possible and is an example of a simple REST API, with vendoring, swaggering, concoursing, testing, benchmarking, linting and deploying.
It is based on the latest go version, 1.9.0.
It follows the best practices in the go community.

## How to use

Choose a project name and replace it in the script below.
Run the script.
You're ready to go!

```
git clone https://github.com/vincentserpoul/gorestarter YOURPROJECTNAME
cd YOURPROJECTNAME
rm -rf .git
find ./ -type f -exec sed -i -e 's/gorestarter/YOURPROJECTNAME/g' {} \;
find ./ -type f -exec sed -i -e 's/loft\/gorestarter/YOURPROJECTGROUP\/YOURPROJECTNAME/g' {} \;
git init
```

If you have a dependency on MySQL, and want to dev locally:
First install docker https://docs.docker.com/engine/installation/
Then (you might have to enable the experimental flag):

```
docker stack deploy --compose-file=docker/compose.yml gorestarter;
CONTAINER_NAME=$(docker ps --format '{{.Names}}' | grep percona) && docker exec -i $CONTAINER_NAME mysql -u root -e "CREATE DATABASE dev;GRANT ALL PRIVILEGES ON dev.* TO 'internal'@'%';";
```

MySQL is now available locally!

If you want to create a new resource:
* copy pkg/resourceone in a new folder pkg/yourresource
* add
```
	r.Mount("/v1", yourresource.Router(db))

	// Resourceone related things
	ddl := &yourresource.DDL{}
	err := ddl.MigrateUp(context.Background(), db)
	if err != nil {
		log.Fatal(err)
	}
```
to pkg/rest/serve.go

Please contribute, comment, post issues...

## Rules & opinions from a long time Golang usage and avid Golang news and articles reader

* Don't rely on too many external packages, go standard lib is very nice, *secure* and *simple*. Exceptions taken here:
    * "github.com/go-sql-driver/mysql" as we obviously need a specific driver for MySQL
    * "github.com/jmoiron/sqlx" as latest go1.8 named params not yet implemented in the mysql driver [coming very very soon](https://github.com/go-sql-driver/mysql/issues/561)
    * "github.com/cloudfoundry-community/go-cfenv" for cloud foundry env parsing
    * "github.com/segmentio/ksuid" for its specific sortable unique id generation (maybe switch to github.com/oklog/ulid, see [this article](https://blog.kowalczyk.info/article/JyRZ/generating-good-random-and-unique-ids-in-go.html) )
    * "github.com/go-chi/chi" for idiomatic routing with middleware
    * "github.com/sirupsen/logrus" for structured logging
* Did I already say: You probably don't need that external package, think twice.
* Don't think frameworks, think libraries.
* Don't use ORM, please. Learn SQL.
* You don't want global vars. Use wrappers and closures instead.
* Context should be used in only in very very few use cases (for now, the only I can thing of is a global requestid). Use wrappers and closures instead.
* Use gometalinter on everything (not a single error should subsist or write a linting exception).
* Test as much as possible (try to maintain your coverage above 80%)
* Don't hide errors!
* You probably don't need channels, use them carefully. They're powerful but add a lot of complexity.
* Benchmark your code, it's easy to do. Most of the time, it's not useful, but it can save you one day (sounds like tests)

## TODO List

[ ] Add more tests
[ ] Change the docker image for the MySQL service creation
[ ] Add Swagger using "github.com/yvasiyarov/swagger"
[ ] Add GRPC/protobuff impl
[ ] Add event listener (kafka listener?) + event launcher