# Jason's Game Land Service
A service for land related requests in conjunction with Jason's Game.

This is a command line application that can run in three different modes:

* Bootstrapper
* Service
* Client

## Configuration
Before using the application, you must create a configuration file, config.json, containing
the cryptographic keys of bootstrapper and service nodes. You can use config.json.example as a
starting point and generate keys with `jasons-game-land-service generate-key`.

## Usage
A bootstrapper instance must be running first:

```
$ jasons-game-land-service bootstrapper -p34001
```

Then a service instance may be launched:

```
$ jasons-game-land-service service -p34002
```

Once the service is running and ready for requests, you may interact with it via a client:
```
$ jasons-game-land-service client build-portal
* Sending request to build portal...
* Successfully sent request to build portal!
```
