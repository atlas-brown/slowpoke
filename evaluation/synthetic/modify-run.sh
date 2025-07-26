#!/bin/bash

cd $(dirname $0)

find . -type f -name "run.sh" | while read -r file; do
  echo "Modifying $file"

  # First make the other replacements
  sed -i '' 's|cd $(dirname $0)/../..|cd $(dirname $0)/../../..|' "$file"
  sed -i '' 's|DIR=synthetic/\$EXP/results|DIR=./evaluation/results|' "$file"
  sed -i '' 's|python3 test.py -b synthetic|python3 src/main.py -b synthetic|' "$file"

  # Process the do/done blocks with proper indentation
  tmpfile=$(mktemp)
  awk '{
    # Add kubectl commands after do
    if (/^[[:space:]]*do[[:space:]]*$/) {
      print $0
      print "    kubectl delete deployments --all"
      print "    kubectl delete services --all"
    }
    # Add kubectl commands before done
    else if (/^[[:space:]]*done[[:space:]]*$/) {
      print "    kubectl delete deployments --all"
      print "    kubectl delete services --all"
      print $0
    }
    else {
      print $0
    }
  }' "$file" > "$tmpfile" && mv "$tmpfile" "$file"
done

echo "✅ All run.sh files modified."