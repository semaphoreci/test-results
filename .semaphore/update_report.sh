#!/bin/bash

set -euo pipefail
set -x

ARTIFACT_PATH="REPORT.md"
TMP_FILE=$(mktemp)

# Print input values for debug
echo "SEMAPHORE_PIPELINE_CREATED_AT: $SEMAPHORE_PIPELINE_CREATED_AT"
echo "SEMAPHORE_PIPELINE_ID: $SEMAPHORE_PIPELINE_ID"
echo "SEMAPHORE_PIPELINE_INIT_DURATION: $SEMAPHORE_PIPELINE_INIT_DURATION"
echo "SEMAPHORE_PIPELINE_QUEUEING_DURATION: $SEMAPHORE_PIPELINE_QUEUEING_DURATION"
echo "SEMAPHORE_PIPELINE_RUNNING_DURATION: $SEMAPHORE_PIPELINE_RUNNING_DURATION"

# Check artifact
if [ ! -f "$ARTIFACT_PATH" ]; then
  echo "❌ $ARTIFACT_PATH not found!"
  exit 1
fi

echo "📄 Current REPORT.md:"
cat "$ARTIFACT_PATH"

# Extract durations and time
PIPELINE_TIME=${SEMAPHORE_PIPELINE_CREATED_AT}
PIPELINE_ID=${SEMAPHORE_PIPELINE_ID}
INIT_DURATION=${SEMAPHORE_PIPELINE_INIT_DURATION}
QUEUE_DURATION=${SEMAPHORE_PIPELINE_QUEUEING_DURATION}
RUN_DURATION=${SEMAPHORE_PIPELINE_RUNNING_DURATION}

START_TS=$PIPELINE_TIME

# Convert to ISO 8601
format_time() {
  date -u -d "@$1" +"%Y-%m-%dT%H:%M:%S"
}

START=$(format_time "$START_TS")

# Mermaid entry
CHART_ENTRY=$(cat <<EOF
    section Pipeline $PIPELINE_ID
    Init       :active, init$PIPELINE_ID, $START, ${INIT_DURATION}s
    Queue      :active, queue$PIPELINE_ID, after init$PIPELINE_ID, ${QUEUE_DURATION}s
    Run        :active, run$PIPELINE_ID, after queue$PIPELINE_ID, ${RUN_DURATION}s
EOF
)

# Create or update chart
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

echo "✅ Updated REPORT.md with pipeline metrics"
