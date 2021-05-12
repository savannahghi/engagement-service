#!/bin/sh

# Decrypt the service account file
# --batch to prevent interactive command --yes to assume "yes" for questions
mkdir -p ~/secrets
gpg --quiet --batch --yes --decrypt --passphrase="$SECRET_PASSPHRASE" \
  --output $CI_PROJECT_DIR/bewell-app-testing.json bewell-app-testing.json.gpg