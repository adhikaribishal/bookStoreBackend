# bookStoreBackend

Golang backend for BookStore

# Setup
1. Install the required packages using
    ```bash
    go mod download
    ```
2. Create a ```.env``` file in the project root directory and copy the variable names from ```env_sample``` file. Add environment variables value in the ```.env``` file accordingly.

3. Run the server using ```go run .``` in the project's root directory.

# Description

This project is solely for practicing API development in golang. The project is for a book store. It supports CRUD operation for a book. You must first create a user and login before performing any operation on a book.

User authentication is handled using JWTs. For generating JWT tokens and authentication I am using ```jwt-go``` package by [dgrijalva](https://www.github.com/dgrijalva/jwt-go).

API routing is handled using Alex Wagner's ShiftPath technique. You can learn more about the technique [here](https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html). In this technique he presents a technique involving a small ShiftPath() helper that returns the first path segment, and shifts the rest of the URL down. The current handler switches on the first path segment, then delegates to sub-handlers which do the same thing on the rest of the URL.
