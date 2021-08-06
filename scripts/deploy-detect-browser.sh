#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SOURCE_DIR="${DIR}/../launch/browser_detect"

gcloud functions \
  deploy detect_browser \
  --source=${SOURCE_DIR} \
  --runtime=python38 \
  --trigger-http \
  --region=europe-west3 \
  --allow-unauthenticated

