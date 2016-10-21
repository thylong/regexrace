#! /bin/bash
if [ "$TRAVIS_BRANCH" = "master" ] && [ "$TRAVIS_PULL_REQUEST" = "false" ]; then
  gcloud --quiet components update kubectl
  # Push to Google container registry
  docker build -t gcr.io/$CLOUDSDK_CORE_PROJECT/$CLOUDSDK_CORE_PROJECT:v1 .
  gcloud docker push gcr.io/$CLOUDSDK_CORE_PROJECT/$CLOUDSDK_CORE_PROJECT:v1 > /dev/null
  # Deploy to the cluster
  gcloud container clusters get-credentials $CLOUDSDK_CORE_PROJECT
  kubectl apply -f config/k8s/mongo.yml
  kubectl apply -f config/k8s/$CLOUDSDK_CORE_PROJECT.yml
  kubectl rolling-update web-controller --image=gcr.io/$CLOUDSDK_CORE_PROJECT/$CLOUDSDK_CORE_PROJECT:v1 --image-pull-policy Always
fi
