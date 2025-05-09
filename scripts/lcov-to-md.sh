#!/bin/bash

LCOV_FILE="coverage.lcov"
OUTPUT_FILE="COVERAGE.md"

# Git info
BRANCH=${SEMAPHORE_GIT_BRANCH:-$(git rev-parse --abbrev-ref HEAD)}
COMMIT=$(git rev-parse HEAD)
AUTHOR=$(git log -1 --pretty=format:'%an')
DATE=$(git log -1 --pretty=format:'%ad')
MESSAGE=$(git log -1 --pretty=format:'%s')
CHANGED_FILES=$(git diff-tree --no-commit-id --name-only -r HEAD)

# Temp file for sorting
TMP_FILE=$(mktemp)

OVERALL_TOTAL=0
OVERALL_COVERED=0
CURRENT_FILE=""
FILE_TOTAL=0
FILE_COVERED=0

while IFS= read -r line; do
  case "$line" in
    SF:*)
      CURRENT_FILE=$(basename "${line#SF:}")
      FILE_TOTAL=0
      FILE_COVERED=0
      ;;
    DA:*)
      count=$(echo "$line" | cut -d',' -f2)
      FILE_TOTAL=$((FILE_TOTAL + 1))
      OVERALL_TOTAL=$((OVERALL_TOTAL + 1))
      if [ "$count" -gt 0 ]; then
        FILE_COVERED=$((FILE_COVERED + 1))
        OVERALL_COVERED=$((OVERALL_COVERED + 1))
      fi
      ;;
    end_of_record)
      if [ -n "$CURRENT_FILE" ]; then
        if [ "$FILE_TOTAL" -gt 0 ]; then
          percent=$(awk "BEGIN { printf \"%.2f\", ($FILE_COVERED/$FILE_TOTAL)*100 }")
        else
          percent="0.00"
        fi
        printf "%07.2f|%s|%d|%d\n" "$percent" "$CURRENT_FILE" "$FILE_COVERED" "$FILE_TOTAL" >> "$TMP_FILE"
      fi
      ;;
  esac
done < "$LCOV_FILE"

if [ "$OVERALL_TOTAL" -gt 0 ]; then
  OVERALL_COVERAGE=$(awk "BEGIN { printf \"%.2f\", ($OVERALL_COVERED/$OVERALL_TOTAL)*100 }")
else
  OVERALL_COVERAGE="0.00"
fi

# 🧮 Extract metrics
LINES=($(tail -n 30 /tmp/system-metrics))

CPU_VALUES=()
MEM_VALUES=()
SYSTEM_DISK_VALUES=()
DOCKER_DISK_VALUES=()
SHM_VALUES=()

for line in "${LINES[@]}"; do
  CPU_VALUES+=("$(echo "$line" | grep -oP 'cpu:\K[0-9.]+')")
  MEM_VALUES+=("$(echo "$line" | grep -oP 'mem:\s*\K[0-9.]+')")
  SYSTEM_DISK_VALUES+=("$(echo "$line" | grep -oP 'system_disk:\s*\K[0-9.]+')")
  DOCKER_DISK_VALUES+=("$(echo "$line" | grep -oP 'docker_disk:\s*\K[0-9.]+')")
  shm=$(echo "$line" | grep -oP 'shared_memory:\s*\K[0-9]+')
  SHM_VALUES+=("$(awk "BEGIN { printf \"%.2f\", ($shm/512)*100 }")")
done

print_ascii_chart() {
  local label="$1"
  local -n values=$2
  local height=10
  local max_width=60
  local total=${#values[@]}

  # Downsample if too wide
  local step=$(( total > max_width ? total / max_width : 1 ))

  local downsampled=()
  local timestamps=()

  for ((i=0; i<total; i+=step)); do
    downsampled+=("${values[i]}")
    timestamps+=("$((i))")  # Could use actual times if stored
  done

  local width=${#downsampled[@]}

  echo "## 📊 $label Usage (last ${#values[@]} samples)"
  echo '```text'

  for ((level=height; level>=0; level--)); do
    threshold=$((level * 10))
    printf "%4d%% |" "$threshold"
    for val in "${downsampled[@]}"; do
      val_int=$(awk "BEGIN { print int($val) }")
      if (( val_int >= threshold )); then
        printf " █"
      else
        printf "  "
      fi
    done
    echo
  done

  # Axis
  printf "      +"
  for ((i=0; i<width; i++)); do
    printf "--"
  done
  echo

  # Label every 10th tick
  printf "       "
  for ((i=0; i<width; i++)); do
    if (( i % 10 == 0 )); then
      printf "%-2d" $((i * step))
    else
      printf "  "
    fi
  done
  echo
  echo '```'
}

# Write markdown report
{
  echo "# 📈 Code Coverage Report"
  echo
  echo "## 🔧 Commit Info"
  echo "- **Branch**: \`$BRANCH\`"
  echo "- **Commit**: \`$COMMIT\`"
  echo "- **Author**: $AUTHOR"
  echo "- **Date**: $DATE"
  echo "- **Message**: _${MESSAGE}_"
  echo

  # ASCII usage graphs
  print_ascii_chart "CPU" CPU_VALUES
  print_ascii_chart "Memory" MEM_VALUES
  print_ascii_chart "System Disk" SYSTEM_DISK_VALUES
  print_ascii_chart "Docker Disk" DOCKER_DISK_VALUES
  print_ascii_chart "Shared Memory (as % of 512MB)" SHM_VALUES

  echo "---"
  echo
  echo "## 🧵 Workflow Debug Info"
  echo
  echo "| Variable | Value |"
  echo "|----------|-------|"
  for var in SEMAPHORE_GIT_BRANCH SEMAPHORE_GIT_COMMITTER SEMAPHORE_JOB_ID SEMAPHORE_PROJECT_NAME SEMAPHORE_PIPELINE_ID SEMAPHORE_WORKFLOW_ID SEMAPHORE_GIT_SHA SEMAPHORE_GIT_REPO_NAME; do
    val="${!var}"
    echo "| \`$var\` | \`$val\` |"
  done
  echo
  echo "---"
  echo
  echo "## 📝 Changed Files"
  echo
  echo '```diff'
  for file in $CHANGED_FILES; do echo "+ $file"; done
  echo '```'
  echo
  echo "---"
  echo
  echo "## 🔍 Per-File Coverage"
  echo
  echo "| File | Coverage | Visual |"
  echo "|------|----------|--------|"

  sort "$TMP_FILE" | while IFS='|' read -r padded file covered total; do
    percent=$(echo "$padded" | sed 's/^0*//')
    bar_len=20
    filled=$(awk "BEGIN { printf \"%d\", ($percent/100)*$bar_len }")
    empty=$((bar_len - filled))
    filled_bar=$(yes '🟩' | head -n "$filled" | tr -d '\n')
    empty_bar=$(yes '⬜' | head -n "$empty" | tr -d '\n')
    echo "| \`$file\` | $percent% ($covered/$total) | $filled_bar$empty_bar |"
  done

  echo
  echo "---"
  echo
  echo "## 📊 Total Coverage"
  echo "**$OVERALL_COVERED / $OVERALL_TOTAL → $OVERALL_COVERAGE%**"
} > "$OUTPUT_FILE"

rm -f "$TMP_FILE"

echo "✅ Report saved to $OUTPUT_FILE"
