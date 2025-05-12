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

# Compute timeline timestamps
START_TS=$PIPELINE_TIME
INIT_END=$((START_TS + INIT_DURATION))
QUEUE_END=$((INIT_END + QUEUE_DURATION))
RUN_END=$((QUEUE_END + RUN_DURATION))

START=$(format_time "$START_TS")
INIT_FINISH=$(format_time "$INIT_END")
QUEUE_FINISH=$(format_time "$QUEUE_END")
RUN_FINISH=$(format_time "$RUN_END")

# Timeline entry
read -r -d '' TIMELINE_ENTRY <<EOF || true
    section Pipeline ${PIPELINE_ID}
    ${START} : Init started
    ${INIT_FINISH} : Queue started
    ${QUEUE_FINISH} : Run started
    ${RUN_FINISH} : Run finished
EOF

# Ensure Pipeline metrics section exists
if ! grep -q "## Pipeline metrics" "$ARTIFACT_PATH"; then
  {
    echo -e "\n## 📊 Pipeline metrics\n"
    echo '```mermaid'
    echo 'timeline'
    echo "$TIMELINE_ENTRY"
    echo '```'
  } >> "$ARTIFACT_PATH"
else
  # Append entry to existing mermaid block
  awk -v entry="$TIMELINE_ENTRY" '
    BEGIN { inside = 0 }
    {
      print
      if ($0 ~ /```mermaid/) {
        inside = 1
      } else if (inside && $0 ~ /^```$/) {
        print entry
        inside = 0
      }
    }
  ' "$ARTIFACT_PATH" > "$TMP_FILE"
  mv "$TMP_FILE" "$ARTIFACT_PATH"
fi

echo "✅ Updated $ARTIFACT_PATH:"
cat "$ARTIFACT_PATH"
