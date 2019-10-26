# Munch Server
![travis](https://travis-ci.com/MunchApp/munchserver.svg?branch=master)

## Running Server

Run `go install` in .
Then, run `go run server.go`
or, for live reloading, `gin -p 80 run server.go`

## Clearing database

To clear things from your localhost database, run the following commands in the mongo shell
```
    use munch
    db.foodTrucks.deleteMany({})
    db.users.deleteMany({})
    db.reviews.deleteMany({})
```
