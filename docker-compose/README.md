## Pocket core deployments using docker-compose


This folder contains various docker-compose setups for deploying pocket-core-service in various ways


#### docker-compose.yml


It's the most simple configuration for deploying pocket-core using docker-compose you only just need to edit your configurations on the `command` section of `pocket-core-service` and use your custom chains_configuration on the env variable `POCKET_CORE_CHAINS`


$ docker-compose up


#### pocket-core-config-by-env-variable.yml


It's similar to `docker-compose.yml` but here we can also inject the options via env variable.


$ docker-compose -f pocket-core-chains-by-env-variable.yml up


#### pocket-core-config-by-s3.yml


Here we use a env variable configuration approach but we fetch the `chains.json` from a s3 location instead of giving the full content on the environment variable


$ docker-compose -f pocket-core-config-by-s3.yml up


#### pocket-core-config-by-volume-mapping.yml


In this setup we just map the file `chains.json` directly to the pocket-core-service configuration folder 


$ docker-compose -f pocket-core-config-by-volume-mapping.yml up
