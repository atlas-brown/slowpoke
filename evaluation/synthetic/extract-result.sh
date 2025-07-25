# #!/bin/bash
# cd $(dirname $0)

# mkdir -p all-results

# for folder in *;
# do 
#     [[ -d "$folder/04-09-pokerpp-rm-deadlock-0maxconn"]] || continue
#     for log_file in "$folder/04-09-pokerpp-rm-deadlock-0maxconn"/*.log
#     do
#         filename=$(basename "$log_file")
#         # Check if the marker line exists
#         if grep -q "^\[test.py\] ++++++++++++++++++++++++++++++++ Summary:" "$log_file"; then
#             # Extract from the line containing the marker to the end of the file
#             # awk '/^\[test.py\] ++++++++++++++++++++++++++++++++ Summary:/,EOF' "$file" >> "$output_file"
#             # echo -e "\n--- Extracted from: $file ---\n" >> "$output_file"
#             echo "Extracting from $log_file" > all-results/"$filename"
#         fi
# done

#!/bin/bash
cd "$(dirname "$0")"

mkdir -p all-results

for folder in *; do 
    subfolder="$folder/04-09-pokerpp-rm-deadlock-0maxconn"
    [[ -d "$subfolder" ]] || continue

    for log_file in "$subfolder"/*.log; do
        [[ -f "$log_file" ]] || continue
        filename=$(basename "$log_file")

        # Check if the marker line exists
        if tail -n 50 "$log_file" | grep -q "^\[test.py\] ++++++++++++++++++++++++++++++++ Summary:" ; then
            # echo "Extracting from $log_file" > all-results/"$filename"
            # Extract content from the marker to the end
            tail -n 31 "$log_file" > all-results/"$filename"
        fi
    done
done
