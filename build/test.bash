#!/usr/bin/env bash
url=$(minikube service awesome-grid --url)
curl -X POST $url/wd/hub/session -d '{"desiredCapabilities":{"browserName":"chrome"}}' > output
cat output
if grep -q "\"status\": 0" output
then
   exit 0
else
   exit 1
fi