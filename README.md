To do:

# deployment
- when i push to main, github actions needs to run tests. if tests pass, it needs to somehow get the binary onto my digital ocean droplet
- the vm needs to run build? or does this need to be copied straight tot the machine?
- vm needs to run the file
- what about env vars? just add manually for now?


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

GET /random => fetches a single quote at random for you 

GET /author?name=<author name> => fetches all quotes from authors matching that name

POST /author => inserts a new quote into the global quote db for everyone
** Content-Type must be application/json
** example of request body: {author":"Donald Trump","message":"CHINA!"}


running locally:

1. install binary from github
2. install sqlite from brew if not already done

to interact via cli with db:
sqlite3 db/quotedb.sqlite
