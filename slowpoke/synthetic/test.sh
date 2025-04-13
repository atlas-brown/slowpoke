for folder in *grpc*
do
    echo "$folder:$(ls $folder/yamls | wc -l)"
done