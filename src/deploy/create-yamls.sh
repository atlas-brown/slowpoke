#!/bin/bash

cd $(dirname $0)

export MUCACHE_TOP=${MUCACHE_TOP:-$(git rev-parse --show-toplevel --show-superproject-working-tree)}

## Make sure that the containers exist and have been built first
export docker_io_username=yizhengx
export application_namespace=${1?application name not given, e.g., social}
export single_node_flag=${2:-12}

# Keep this for now as it might be used for condition checks
export cm_enabled=false
export ttl=0

mkdir -p ../${application_namespace}/yamls
echo "Caution: This will overwrite any existing files in ../${application_namespace}/yamls"

declare -A all_services
all_services["singleservice"]="service"
all_services["twoservices"]="caller callee"
all_services["chain"]="service1 service2 service3 service4 backend"
all_services["chain3"]="service1 service2 backend"
all_services["chain4"]="service1 service2 service3 backend"
all_services["star"]="frontend backend1 backend2 backend3 backend4"
all_services["fanin"]="frontend1 frontend2 frontend3 frontend4 backend"
all_services["loadcm"]="stub loader"

all_services["movie"]="cast_info compose_review frontend movie_id movie_info movie_reviews page plot review_storage unique_id user user_reviews"
all_services["hotel"]="frontend profile rate reservation search user"
all_services["social"]="post_storage home_timeline user_timeline social_graph compose_post"
all_services["boutique"]="cart checkout currency email frontend payment product_catalog recommendations shipping"

## TODO: Add acmeair, and onlineboutique

services=(${all_services[$application_namespace]})

## Services
for idx in "${!services[@]}"; do
  app_name=${services[$idx]}
  app_name_no_underscores=${app_name//_/}
  if [ "$single_node_flag" = "standard" ]; then
    node_idx=$((idx + 1))
  else
    NUM=$single_node_flag
    node_idx=$((idx%NUM + 1))
  fi
  NODE_IDX="${node_idx}" \
    CM_ENABLED="$cm_enabled" \
    EXPIRATION_TTL="$ttl" \
    APP_NAMESPACE="$application_namespace" \
    APP_NAME="$app_name" \
    APP_NAME_NO_UNDERSCORES="$app_name_no_underscores" \
    APP_NAME_NO_UNDERSCORES_UPPERCASE="${app_name_no_underscores^^}" \
    SLOWPOKE_PRERUN="\${SLOWPOKE_PRERUN}" \
    SLOWPOKE_POKER_BATCH_THRESHOLD="\${SLOWPOKE_POKER_BATCH_THRESHOLD}" \
    envsubst <"./app.yaml" > "../${application_namespace}/yamls/$app_name.yaml"
done