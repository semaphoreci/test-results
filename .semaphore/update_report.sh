#!/bin/bash

set -euo pipefail

ARTIFACT_PATH="REPORT.md"
TMP_FILE=$(mktemp)

# Fetch the artifact (adjust if needed for actual path or artifact download)
if [ ! -f "$ARTIFACT_PATH" ]; then
  echo "Artifact $ARTIFACT_PATH not found!"
  exit 1
fi

# Gather data
PIPELINE_TIME=${SEMAPHORE_PIPELINE_CREATED_AT}
PIPELINE_ID=${SEMAPHORE_PIPELINE_ID}
INIT_DURATION=${SEMAPHORE_PIPELINE_INIT_DURATION}
QUEUE_DURATION=${SEMAPHORE_PIPELINE_QUEUEING_DURATION}
RUN_DURATION=${SEMAPHORE_PIPELINE_RUNNING_DURATION}

# Calculate Gantt start and end points
START_TS=$PIPELINE_TIME
INIT_END=$((START_TS + INIT_DURATION))
QUEUE_END=$((INIT_END + QUEUE_DURATION))
RUN_END=$((QUEUE_END + RUN_DURATION))

# Convert Unix time to ISO8601 for Mermaid (optional: date -u -Iseconds)
format_time() {
    date -u -d "@$1" +"%Y-%m-%dT%H:%M:%S"
}

START=$(format_time "$START_TS")
INIT_END_TIME=$(format_time "$INIT_END")
QUEUE_END_TIME=$(format_time "$QUEUE_END")
RUN_END_TIME=$(format_time "$RUN_END")

# Create Gantt entries
CHART_ENTRY=$(cat <<EOF
    section Pipeline $PIPELINE_ID
    Init       :active, init$PIPELINE_ID, $START, $INIT_DURATIONs
    Queue      :active, queue$PIPELINE_ID, after init$PIPELINE_ID, $QUEUE_DURATIONs
    Run        :active, run$PIPELINE_ID, after queue$PIPELINE_ID, $RUN_DURATIONs
EOF
)

# Check if REPORT.md already has the chart
if ! grep -q "# Pipeline metrics" "$ARTIFACT_PATH"; then
    echo "# Pipeline metrics" > "$TMP_FILE"
    echo "" >> "$TMP_FILE"
    echo '```mermaid' >> "$TMP_FILE"
    echo "gantt" >> "$TMP_FILE"
    echo "    title Pipeline durations" >> "$TMP_FILE"
    echo "$CHART_ENTRY" >> "$TMP_FILE"
    echo '```' >> "$TMP_FILE"
    echo "" >> "$TMP_FILE"
    cat "$ARTIFACT_PATH" >> "$TMP_FILE"
    mv "$TMP_FILE" "$ARTIFACT_PATH"
else
    # Insert new chart row just before the closing ```
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

echo "Updated REPORT.md with pipeline metrics."
