
# Instructions
First you must setup your environments.

via docker -> <br/>
> set up in docker-compose.yml file.<br/>

via command line -><br/>
> export LINE_LOGIN_CODE={value}<br/>
> export LINE_LOGIN_REDIRECT_URI={value}<br/> 
> export LINE_LOGIN_CLIENT_ID={value}<br/>
> export LINE_LOGIN_CLIENT_SECRET={value}<br/>
> export LINE_LOGIN_API_URL={value} (optional)<br/>
> export HEALTHCHECK_REPORT_URL={value} (optional)<br/>

# Installation
via docker -> 
> docker-compose up --build<br/>

via command line -> 
> go install -v ./

# Usage
via command line -> 
> go-healthcheck {csv-file}

# Running Test
> go test ./ -v
