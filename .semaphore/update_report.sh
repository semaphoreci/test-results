#!/bin/bash
set -euo pipefail

ARTIFACT_PATH="REPORT.md"
TMP_FILE=$(mktemp)

# Print current report
echo "📄 Current $ARTIFACT_PATH:"
cat "$ARTIFACT_PATH"

# Capture environment variables
PIPELINE_TIME=${SEMAPHORE_PIPELINE_CREATED_AT}
PIPELINE_ID=${SEMAPHORE_PIPELINE_ID}
INIT_DURATION=${SEMAPHORE_PIPELINE_INIT_DURATION}
QUEUE_DURATION=${SEMAPHORE_PIPELINE_QUEUEING_DURATION}
RUN_DURATION=${SEMAPHORE_PIPELINE_RUNNING_DURATION}

# Skip if all durations are zero
if [[ "$INIT_DURATION" == "0" && "$QUEUE_DURATION" == "0" && "$RUN_DURATION" == "0" ]]; then
  echo "ℹ️ Skipping chart entry due to zero durations"
  exit 0
fi

# Format UNIX timestamp to ISO
format_time() {
  date -u -d "@$1" +"%Y-%m-%dT%H:%M:%S"
}

# Compute timestamps for each stage
START_TS=$PIPELINE_TIME
INIT_END=$((START_TS + INIT_DURATION))
QUEUE_END=$((INIT_END + QUEUE_DURATION))
RUN_END=$((QUEUE_END + RUN_DURATION))

START=$(format_time "$START_TS")
INIT_FINISH=$(format_time "$INIT_END")
QUEUE_FINISH=$(format_time "$QUEUE_END")
RUN_FINISH=$(format_time "$RUN_END")

# Gantt chart entry
read -r -d '' GANTT_ENTRY <<EOF || true
    section Pipeline ${PIPELINE_ID}
    Init :active, init_${PIPELINE_ID}, ${START}, ${INIT_DURATION}s
    Queue :active, queue_${PIPELINE_ID}, ${INIT_FINISH}, ${QUEUE_DURATION}s
    Run :active, run_${PIPELINE_ID}, ${QUEUE_FINISH}, ${RUN_DURATION}s
EOF

# Ensure Pipeline metrics section exists
if ! grep -q "## 📊 Pipeline metrics" "$ARTIFACT_PATH"; then
  {
    echo -e "\n## 📊 Pipeline metrics\n"
    echo "```mermaid"
    echo "gantt"
    echo "    title Pipeline durations"
    echo "$GANTT_ENTRY"
    echo "```"
  } >> "$ARTIFACT_PATH"
else
  # Append or update the mermaid chart inside the Pipeline metrics section
  awk -v entry="$GANTT_ENTRY" '
    BEGIN { inside = 0 }
    {
      print
      if ($0 ~ /## 📊 Pipeline metrics/) {
        inside = 1
      }
      if (inside && $0 ~ /```mermaid/) {
        print "```mermaid"
        print "gantt"
        print "    title Pipeline durations"
        print entry
        inside = 0
      }
    }
  ' "$ARTIFACT_PATH" > "$TMP_FILE"
  mv "$TMP_FILE" "$ARTIFACT_PATH"
fi

echo "✅ Updated $ARTIFACT_PATH:"
cat "$ARTIFACT_PATH"
