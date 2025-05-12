#!/bin/bash

set -euo pipefail

ARTIFACT_PATH="REPORT.md"
TMP_FILE=$(mktemp)

# Ensure the artifact exists
if [ ! -f "$ARTIFACT_PATH" ]; then
  echo "Artifact $ARTIFACT_PATH not found!"
  exit 1
fi

# Read pipeline timing data from Semaphore-provided environment variables
PIPELINE_TIME=${SEMAPHORE_PIPELINE_CREATED_AT}
PIPELINE_ID=${SEMAPHORE_PIPELINE_ID}
INIT_DURATION=${SEMAPHORE_PIPELINE_INIT_DURATION}
QUEUE_DURATION=${SEMAPHORE_PIPELINE_QUEUEING_DURATION}
RUN_DURATION=${SEMAPHORE_PIPELINE_RUNNING_DURATION}

# Calculate timestamps for start and phase boundaries
START_TS=$PIPELINE_TIME
INIT_END=$((START_TS + INIT_DURATION))
QUEUE_END=$((INIT_END + QUEUE_DURATION))
RUN_END=$((QUEUE_END + RUN_DURATION))

# Convert to ISO8601 for readability in the Gantt chart
format_time() {
    date -u -d "@$1" +"%Y-%m-%dT%H:%M:%S"
}

START=$(format_time "$START_TS")

# Create Gantt chart entry for this pipeline
CHART_ENTRY=$(cat <<EOF
    section Pipeline $PIPELINE_ID
    Init       :active, init$PIPELINE_ID, $START, ${INIT_DURATION}s
    Queue
