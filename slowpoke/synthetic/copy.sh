# fanout-w3-http-sync
cd $(dirname $0)
for folder in *http*
do
    if [[ $folder == *fanout-w10* ]]; then
        continue
    fi
    IFS='-' read -r mode width protocol pattern <<< "$folder"
    # copy folder/run.sh to $mode-$width-grpc-$pattern/run.sh
    new_folder="$mode-$width-grpc-$pattern"
    cp $folder/run.sh $new_folder/run.sh
    rm -rf $new_folder/yamls
    cp -r /home/ubuntu/mucache/slowpoke/synthetic/fanout-w10-grpc-sync/yamls $new_folder/yamls
done