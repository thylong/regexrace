#! /bin/bash
if [ "$TRAVIS_BRANCH" = "master" ] && [ "$TRAVIS_PULL_REQUEST" = "false" ]; then
  gcloud container clusters get-credentials $PROJECT_ID
  kubectl apply -f config/k8s/mongo.yml
  kubectl apply -f config/k8s/$PROJECT_ID.yml
  kubectl rolling-update web-controller --image=gcr.io/$PROJECT_ID/$PROJECT_ID:v1 --image-pull-policy Always
fi
