
# Instructions
First you must setup your environments.

via docker -> set up in docker-compose.yml file.
via command line ->
export LINE_LOGIN_CODE={value}
export LINE_LOGIN_REDIRECT_URI={value}
export LINE_LOGIN_CLIENT_ID={value}
export LINE_LOGIN_CLIENT_SECRET={value}
export LINE_LOGIN_API_URL={value} (optional)
export HEALTHCHECK_REPORT_URL={value} (optional)

# Installation
via docker -> docker-compose up --build
via command line -> go install -v .

# Usage
via command line -> go-healthcheck {csv-file}

# Running Test
go test ./ -v
