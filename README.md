!Note
project not deployed because digital ocean is refusing to allow outgoing smtp and that is what my auth currently hinges upon :( 
Need to find a workaround at some point

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

# Instructions

This application can be run either natively or using Docker. Supported architectures: `arm64`, `amd64`

## Option 1: Docker 
1. Pull the image:
   ```
   docker pull jamoowen/quoteapi:latest
   ```
2. Create a `.env` file in your project directory with your secrets:
   ```
   SECRET_KEY=your_secret_here
   # Add other required environment variables
   ```
3. Run with docker compose:
   ```
   docker compose up
   ```

## Option 2: Native Installation
1. Clone the repository
2. Install Go (if not already installed)
3. Create a `.env` file with your secrets (as above)
4. Use the following make commands:

   ```bash
   # Initialize the database (required first time)
   make init-db

   # Build the application
   make build

   # Run the server directly without building
   make run
   
   # Or build and start the server
   make build
   make start-server

   # Run tests
   make test

   # Clean build artifacts
   make clean

   # For development: recreate database with sample data
   make recreate-dev-db
   ```

Required Environment Variables:
- SECRET_KEY: Your secret key
- [List any other required env variables]

Note: The application uses SQLite for data storage. Database files will be created in the `db` directory.
