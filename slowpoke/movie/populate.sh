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

python3 populate.py \
  --cast_info $(kip cast_info) \
  --compose_review $(kip compose_review) \
  --frontend $(kip frontend) \
  --movie_id $(kip movie_id) \
  --movie_info $(kip movie_info) \
  --movie_reviews $(kip movie_reviews) \
  --page $(kip page) \
  --plot $(kip plot) \
  --review_storage $(kip review_storage) \
  --unique_id $(kip unique_id) \
  --user $(kip user) \
  --user_reviews $(kip user_reviews)
