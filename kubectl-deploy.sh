#! /bin/bash
if [ "$TRAVIS_BRANCH" = "master" ] && [ "$TRAVIS_PULL_REQUEST" = "false" ]; then
  gcloud container clusters get-credentials $CLOUDSDK_CORE_PROJECT
  kubectl apply -f config/k8s/mongo.yml
  kubectl apply -f config/k8s/$CLOUDSDK_CORE_PROJECT.yml
  kubectl rolling-update web-controller --image=gcr.io/$CLOUDSDK_CORE_PROJECT/$CLOUDSDK_CORE_PROJECT:v1 --image-pull-policy Always
fi
