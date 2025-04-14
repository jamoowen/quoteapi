To do:

- rate limiter
- tests for rate limiter
- auth middleware
- connection pool

# Quote Api

REST API for returning quotes from authors

## Version

v0.0.1

## Authentication

# you must first sign up for an api key

in your browser go to /authenticate

1. Enter your email
2. Check your emails for the OTP sent (might be in spam)
3. Go back to the page and enter your email again and the OTP

!Note if you already have an api key, doing this again will delete it and recreate a new one

## Routes

GET /random
GET /random/author

running locally:

1. install binary from github
2. install sqlite from brew if not already done

to interact via cli with db:
sqlite3 db/quotedb.sqlite
