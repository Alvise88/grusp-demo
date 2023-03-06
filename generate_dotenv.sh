#!/usr/bin/env bash

set -x

# Access argo server using service account
#SECRET=$(kubectl -n cicd get sa operate-workflow-sa  -o=jsonpath='{.secrets[0].name}')
#export ARGO_TOKEN="Bearer $(kubectl -n cicd get secret $SECRET -o=jsonpath='{.data.token}' | base64 --decode)"

ARGO_TOKEN="Bearer $(kubectl -n cicd get secret argoci.service-account-token -o=jsonpath='{.data.token}' | base64 --decode)"

export ARGO_IP=$(ifconfig eth0 | awk '/inet / {print $2; }' | cut -d ' ' -f 2 |  tr -d '\n')
# export ARGO_SERVER="${ARGO_IP}:443"
export ARGO_SERVER="argowf.grusp.io"

echo -n "" > .env # reset .env file

echo "ARGO_SECURE=true" >> .env
echo "ARGO_SERVER=${ARGO_SERVER}" >> .env
echo "ARGO_INSECURE_SKIP_VERIFY=true" >> .env
echo "ARGO_TOKEN=\"$ARGO_TOKEN\"" >> .env