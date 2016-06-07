Comstock is a Heroku log drain that publishes to Kafka, written in Go.

## Setup

Set the required environment variables.

```
heroku config:set COMSTOCK_KAFKA_TOPIC=logs
heroku config:set COMSTOCK_USERNAME=admin
heroku config:set COMSTOCK_PASSWORD=changeme
```

Create a Heroku Kafka cluster and a new topic. 

```
heroku plugins:install heroku-kafka
heroku addons:create heroku-kafka:beta-dev
heroku kafka:wait
heroku kafka:create logs
```

Once the application is deployed and running, attach the drain to another running Heroku application.

```
heroku drains:add https://COMSTOCK_USERNAME:COMSTOCK_PASSWORD@mylogdrain.herokuapp.com/logs
```

Watch your logs stream in

```
heroku kafka:tail logs
```
