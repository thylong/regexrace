#! /bin/bash
echo $GCLOUD_KEY | base64 --decode > gcloud.p12
gcloud auth activate-service-account $GCLOUD_EMAIL --key-file gcloud.p12
ssh-keygen -f ~/.ssh/google_compute_engine -N ""
