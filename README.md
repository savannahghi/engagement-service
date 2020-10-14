# Debug micro-service

[![pipeline status](https://gitlab.slade360emr.com/go/debug/badges/develop/pipeline.svg)](https://gitlab.slade360emr.com/go/debug/-/commits/develop)
[![coverage report](https://gitlab.slade360emr.com/go/debug/badges/develop/coverage.svg)](https://gitlab.slade360emr.com/go/debug/-/commits/develop)

A library and service with debug and monitoring helpers e.g a debug endpoint, healthcheck endpoint etc.

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

- `GOOGLE_CLOUD_PROJECT`
- `FIREBASE_WEB_API_KEY`
