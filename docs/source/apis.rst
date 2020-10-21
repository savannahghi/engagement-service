APIS
==================

GraphQL
--------
The development/testing GraphQL service can be found at https://feed-testing-uyajqt434q-ew.a.run.app/ide .

TODO Verify GraphQL and assemble examples as Grapqhlbin

REST
-----
TODO Verify REST and assemble POSTman collection

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
