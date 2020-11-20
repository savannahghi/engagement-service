#!/bin/sh
# grant Cloud Pub/Sub the permission to create tokens
export PUBSUB_SERVICE_ACCOUNT="service-${GOOGLE_PROJECT_NUMBER}@gcp-sa-pubsub.iam.gserviceaccount.com"
gcloud projects add-iam-policy-binding ${GOOGLE_CLOUD_PROJECT} --member="serviceAccount:${PUBSUB_SERVICE_ACCOUNT}" --role='roles/iam.serviceAccountTokenCreator'
