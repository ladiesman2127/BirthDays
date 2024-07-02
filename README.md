# Birthdays
Golang REST API service for getting friends birthday notifications.

Built with Gin and MongoDB

# Commands
Supports the next REST API commands:

1. POST /api/signup/ - Create a new user.
2 .POST /api/auth/ - Auth to get JWT tokens to use the commands below
3. POST /api/user/:id - Get specific user data by his ID
4. POST /api/users - Get all users data
5. POST /api/addFriend/:id - Add the user to the current user's friends list; therefore, he will get a notification about birthdays.
6. POST /api/removeFriend/:id - Remove the user from the current user's friends list; therefore, he will not get a notification about birthdays.
7. POST /api/birthdays - Show all friends that have birthdays today.

## Project structure
- api - contains files to define routes and server side interaction-
- cmd - contains file to start the service
- internal - all project files that should not be used from the outside
   - models - types that are used to structure data
   - controllers - functions that can interact with models and api
  - database - mongoDB helpers
- utils - helper files that can be used outside of the project
- tests - unit-test files
- Dockerfile - file to containerize project
- docker-compose - file to create composed container with mongoDB

## How to use

1. You can use Docker to make it easier
```sh
docker-compose up --build .
```
2. Or you can use your own mongoDB on local machine (make sure that you've change mongodb port in .env file).

   For that you can use command
```sh
go run ./cmd/app/main.go
```
# How to run tests
To run tests use command
```sh
go test ./tests
```
