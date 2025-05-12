#!/bin/bash

set -euo pipefail

ARTIFACT_PATH="REPORT.md"
TMP_FILE=$(mktemp)

# Ensure the artifact exists
if [ ! -f "$ARTIFACT_PATH" ]; then
  echo "Artifact $ARTIFACT_PATH not found!"
  exit 1
fi

# Read pipeline timing data
PIPELINE_TIME=${SEMAPHORE_PIPELINE_CREATED_AT}
PIPELINE_ID=${SEMAPHORE_PIPELINE_ID}
INIT_DURATION=${SEMAPHORE_PIPELINE_INIT_DURATION}
QUEUE_DURATION=${SEMAPHORE_PIPELINE_QUEUEING_DURATION}
RUN_DURATION=${SEMAPHORE_PIPELINE_RUNNING_DURATION}

# Calculate phase boundaries
START_TS=$PIPELINE_TIME
INIT_END=$((START_TS + INIT_DURATION))
QUEUE_END=$((INIT_END + QUEUE_DURATION))
RUN_END=$((QUEUE_END + RUN_DURATION))

# Convert to ISO8601
format_time() {
  date -u -d "@$1" +"%Y-%m-%dT%H:%M:%S"
}

START=$(format_time "$START_TS")

# Create Mermaid Gantt entry
read -r -d '' CHART_ENTRY <<EOF
    section Pipeline $PIPELINE_ID
    Init       :active, init$PIPELINE_ID, $START, ${INIT_DURATION}s
    Queue      :active, queue$PIPELINE_ID, after init$PIPELINE_ID, ${QUEUE_DURATION}s
    Run        :active, run$PIPELINE_ID, after queue$PIPELINE_ID, ${RUN_DURATION}s
EOF

# Insert into or create REPORT.md
if ! grep -q "# Pipeline metrics" "$ARTIFACT_PATH"; then
  {
    echo "# Pipeline metrics"
    echo
    echo '```mermaid'
    echo "gantt"
    echo "    title Pipeline durations"
    echo "$CHART_ENTRY"
    echo '```'
    echo
    cat "$ARTIFACT_PATH"
  } > "$TMP_FILE"
  mv "$TMP_FILE" "$ARTIFACT_PATH"
else
  # Append to existing chart before closing ```
  awk -v entry="$CHART_ENTRY" '
  BEGIN {in_chart=0}
  /```mermaid/ {in_chart=1; print; next}
  /```/ {
    if (in_chart) {
      print entry
      in_chart=0
    }
    print
    next
  }
  {print}
  ' "$ARTIFACT_PATH" > "$TMP_FILE"
  mv "$TMP_FILE" "$ARTIFACT_PATH"
fi

echo "✅ REPORT.md updated with pipeline metrics for $PIPELINE_ID"
