# Heroes

[![Maintainability](https://api.codeclimate.com/v1/badges/b296b77da374de5180ae/maintainability)](https://codeclimate.com/github/bliuchak/heroes/maintainability)

Http server which provides basic CRUD functionality about superheroes (ID and Name).
I'm using here http package from stdlib, gorillamux and Redis for storage.

What's possible:
- Get single hero
- Get all heroes
- Create new hero
- Delete hero

## Motivation

Learn how to build good and practical http servers using goland and std http package.
I'd like to do at least one small change in this repo each day to impove it.

If YOU would like to join me please do! Let's learn this together :)

## How to run

You may use plain docker:

```bash
docker build -t heroes:0.0.1 .
docker run -d -p 6379:6379 --name heroredis --network=local redis:latest
docker run -e APP_PORT=3001 -e DB_HOST=heroredis -e DB_PORT=6379 --network=local heroes:0.0.1
```

Or docker-compose:

```bash
docker-compose build
docker-compose up
```

## License

MIT
