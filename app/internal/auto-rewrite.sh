# for each *.go in each folder in current directory
#  if there's "github.com/eniac/mucache/pkg/state" or "github.com/eniac/mucache/pkg/invoke", replace it with "github.com/eniac/mucache/pkg/slowpoke", if both are in the file, replace one and delete the other
#  replace "state.SetState" with "slowpoke.SetState"
#  replace "state.GetState" with "slowpoke.GetState"
#  replace "state.GetBulkStateDefault[" with "slowpoke.GetBulkStateDefault["
#  replace "state.GetBulkState[" with "slowpoke.GetBulkState["
#  replace "state.SetBulkState" with "slowpoke.SetBulkState"
#  replace "invoke.Invoke" with "slowpoke.Invoke"

#!/bin/bash

benchmark=$1
if [ -z "$benchmark" ]; then
    echo "Usage: $0 <benchmark>"
    exit 1
fi
cd $(dirname $0)/$benchmark

# Find all *.go files in subdirectories
find . -type f -name "*.go" | while read -r file; do
    echo "Processing $file"

    # Check if both "state" and "invoke" imports exist in the file
    if grep -q "github.com/eniac/mucache/pkg/state" "$file" && grep -q "github.com/eniac/mucache/pkg/invoke" "$file"; then
        # Replace "state" import with "slowpoke" and remove "invoke"
        sed -i 's|github.com/eniac/mucache/pkg/state|github.com/eniac/mucache/pkg/slowpoke|' "$file"
        sed -i '/github.com\/eniac\/mucache\/pkg\/invoke/d' "$file"
    else
        # Replace either "state" or "invoke" with "slowpoke"
        sed -i 's|github.com/eniac/mucache/pkg/state|github.com/eniac/mucache/pkg/slowpoke|g' "$file"
        sed -i 's|github.com/eniac/mucache/pkg/invoke|github.com/eniac/mucache/pkg/slowpoke|g' "$file"
    fi

    # Replace function and variable calls
    sed -i 's|state.SetState|slowpoke.SetState|g' "$file"
    sed -i 's|state.GetState|slowpoke.GetState|g' "$file"
    sed -i 's|state.GetBulkStateDefault\[|slowpoke.GetBulkStateDefault\[|g' "$file"
    sed -i 's|state.GetBulkState\[|slowpoke.GetBulkState\[|g' "$file"
    sed -i 's|state.SetBulkState|slowpoke.SetBulkState|g' "$file"
    sed -i 's|invoke.Invoke|slowpoke.Invoke|g' "$file"

done

echo "Processing completed."

cd ..