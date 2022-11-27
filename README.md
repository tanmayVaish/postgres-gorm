# postgres-gorm

It is a pretty basic integration repository. Here, We are integrating PostgreSQL with Go Backend Server Using GORM (A Go ORM).

Moreover, I have created some endpoints listed below:

``` /api/create_books
    /api/get_books/
    /api/get_book/:id
    /api/delete_books/:id
```

## Using Postgres

There also a `docker-compose.yaml` file which you can utilize to test out the server. 

- Just make sure you have `docker` installed in your system.
- Then, In your repository, type:

```fish
docker compose up
```

> Note: I am using `5434` as the `PORT` for running PostgreSQL. Change it according you.

# Testing the Server

There is a file named `thunder-collection_Postgres_Gorm.json` which you can import it to your Thunder Client Extension and Test out the endpoints.
