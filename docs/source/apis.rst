APIS
==================

GraphQL
-------

TODO Verify this from the gateway and assemble examples as GraphQLbin

REST
-----

Before you start
+++++++++++++++++
In order to manually test the interservice REST APIs, you need to generate a
bearer token. There's a CLI in the `login` service 
( https://gitlab.slade360emr.com/go/login ). For token generation to work,
your environment must have the *JWT_KEY* environment variable set to the same
value as it is on the prod/testing/staging service you are accessing. You can
find this value out by inspecting CI environment variables at
https://gitlab.slade360emr.com/groups/go/-/settings/ci_cd and looking for
*PROD_JWT_KEY*, *STAGING_JWT_KEY* and *TESTING_JWT_KEY*.

These JWT keys *expire after one minute* so you need to keep on regenerating
the token(s).

.. code::

    > go run cmd/interservice_login_client/main.go
    2020/11/08 10:31:14 Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MDQ4MjA3MzQsImlhdCI6MTYwNDgyMDY3NH0.8FeHwOMl2T3nN66c6pSPDWf4rln-p3AsNRp80A3MF-4

The examples have been developed using Postman ( https://www.postman.com/downloads/ ).

The POSTman collection that documents the interservice REST API can be found at
https://www.getpostman.com/collections/d023a368982c38ed7c66 .

Events
-------

Events will be used to enable rapid experimentation, rules for
*streaming analytics*.

For *rapid experimentation*, scripts that react to events published by the
feed shall be used to prototype e.g new customer communications, enrollment of
customers into experimental cohorts etc.

Our *rule engine* will apply pre-built or custom logic blocks to these events.

In order to get a *realtime view* of what's happening, these events shall also
be denormalized and streamed into our data lake or data warehouse.

TODO Go starter scripts to react to Events, deployed as Cloud Functions
