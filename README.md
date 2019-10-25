# Munch Server

## Running Server

Run `go install` in .
Then, run `go run server.go`
or, for live reloading, `gin -p 80 run server.go`


To clear things from your lcoalhost database, in command prompt run the following commands:
    mongo.exe
    use munch
    db.foodTrucks.deleteMany({})
