#!/bin/bash

cd $(dirname $0)

# Get the IP address of a service
function kip() {
  local service
  if [ -z "$1" ]; then
    echo "Please provide a service name"
    return 1
  fi
  service=${1//_/}
  kubectl get svc "$service" | tail -n 1 | awk '{print $3}'
}

python3 populate.py --post_size 100 --number_posts_per_user 10 --social_graph $(kip social_graph) --compose_post $(kip compose_post)