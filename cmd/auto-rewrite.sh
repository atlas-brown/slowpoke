#!/bin/bash


# for each main.go in each folder in current directory
#   replace "github.com/eniac/mucache/pkg/cm" with "github.com/eniac/mucache/pkg/slowpoke"
#   replace "go cm.ZmqProxy()" with "slowpoke.SlowpokeInit()"
#   for every line starting with "func" except "func heartbeat(...)", insert a line after it slowpoke.SlowpokeCheck("function_name"), the function name can be extracted from func function_name(...)

benchmark=$1
if [ -z "$benchmark" ]; then
    echo "Usage: $0 <benchmark>"
    exit 1
fi
cd $(dirname $0)/$benchmark

# Find all main.go files in subdirectories
find . -type f -name "main.go" | while read -r file; do
    echo "Processing $file"

    # Replace import path
    sed -i 's|github.com/eniac/mucache/pkg/cm|github.com/eniac/mucache/pkg/slowpoke|g' "$file"

    # Replace go function call
    sed -i 's|go cm.ZmqProxy()|slowpoke.SlowpokeInit()|g' "$file"

    # Process function definitions to add slowpoke.SlowpokeCheck()
    awk '
    /^func / {
        if ($2 !~ /^heartbeat\(/) {
            match($2, /^([^(]*)/, arr);
            print $0;
            print "    slowpoke.SlowpokeCheck(\"" arr[1] "\");";
        } else {
            print $0;
        }
    }
    !/^func / { print $0 }
    ' "$file" > tmpfile && mv tmpfile "$file"
done

echo "Processing completed."

cd ..
