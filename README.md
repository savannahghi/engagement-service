# Feed service

[![pipeline status](https://gitlab.slade360emr.com/go/feed/badges/develop/pipeline.svg)](https://gitlab.slade360emr.com/go/feed/-/commits/develop)
[![coverage report](https://gitlab.slade360emr.com/go/feed/badges/develop/coverage.svg)](https://gitlab.slade360emr.com/go/feed/-/commits/develop)

A service that fetches and preprocesses content for the feed,library and faqs section in Bewell app.

## Interservice API

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/d023a368982c38ed7c66)

## JSON Schema Files
This project uses JSON Schema to validate inputs and outputs.

### Schema file hosting

In order for module references to work, the schema files in 
`graph/feed/schema` need to be hosted. We use 
https://firebase.google.com/docs/hosting for that.

The schema files are hosted at https://schema.healthcloud.co.ke/
e.g https://schema.healthcloud.co.ke//event.schema.json .

This is in the `bewell-app` project. You need to 
`npm install -g firebase-tools` and `firebase login` first. After that,
any time the schema files change, run `firebase deploy` to host the
updated files.

## Environment variables

For local development, you need to *export* the following env vars:

```bash
# Google Cloud Settings
export GOOGLE_APPLICATION_CREDENTIALS="<a path to a Google service account JSON file>"
export GOOGLE_CLOUD_PROJECT="<the name of the project that the service account above belongs to>"
export FIREBASE_WEB_API_KEY="<an API key from the Firebase console for the project mentioned above>"

# Go private modules
export GOPRIVATE="gitlab.slade360emr.com/go/*,gitlab.slade360emr.com/optimalhealth/*"
```

The server deploys to Google Cloud Run. For Cloud Run, the necessary environment
variables are:

- `GHOST_CMS_API_ENDPOINT`
- `GHOST_CMS_API_KEY`

## Service architecture

The design of this service aspires to follow the principles of _domain driven
design_ and _hexagonal architecture_.

For the feed, the domain object is _feed.Feed_ . The aggregate is 
_feed.FeedAggregate_. There's a _feed.Repository_ interface that can be used
to adapt to alternative databases. There's a _feed.NotificationService_
interface that can be used to adapt to alternative message buses.

The communications to the outside world occur over both REST and GraphQL. At
the time of writing, there's a plan to add gRPC and messaging ports.
