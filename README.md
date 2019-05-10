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

## Internals
### Bootstrapper
Bootstrappers create a relay type libp2p host for other peers to bootstrap against.

### Service
Each service node creates a libp2p host that bootstraps itself against bootstrappers that it
finds via its configuration. After bootstrapping, it spawns an actor to handle requests from peers
and finally remote routing gets set up (using the Tupelo SDK). The remote routing subsystem routes
messages to the actor as they arrive over the wire from other peers.

### Client
The client detects service nodes via its configuration, and creates a corresponding representation
of a remote actor for each of them. After doing so it creates a libp2p host and bootstraps against
the bootstrappers it can find via its configuration. Finally, it starts the Tupelo SDK
remote routing subsystem and sends its request to the service, for example to build a portal.
The request gets passed as a message to the actor representing the actor running remotely in 
the service.
