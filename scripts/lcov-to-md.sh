#!/bin/bash

# === CONFIG ===
LCOV_FILE="coverage.lcov"
OUTPUT_FILE="COVERAGE.md"
GOSEC_REPORT="gosec_report.txt"

# === CHECK FOR REQUIRED TOOLS ===
for tool in go gosec golint; do
  if ! command -v "$tool" &> /dev/null; then
    echo "❌ Required tool '$tool' not installed or not in PATH."
    exit 1
  fi
done

# === GIT INFO ===
BRANCH=$(git rev-parse --abbrev-ref HEAD)
COMMIT=$(git rev-parse HEAD)
AUTHOR=$(git log -1 --pretty=format:'%an')
DATE=$(git log -1 --pretty=format:'%ad')
MESSAGE=$(git log -1 --pretty=format:'%s')
CHANGED_FILES=$(git diff-tree --no-commit-id --name-only -r HEAD)

# === GO COVERAGE ===
COVERAGE_OUT="coverage.out"
go test -coverprofile="$COVERAGE_OUT" ./... > /dev/null
if [ ! -f "$COVERAGE_OUT" ]; then
  echo "❌ Coverage file not found. Tests may have failed."
  exit 1
fi

# === PARSE COVERAGE DATA ===
OVERALL_COVERED=0
OVERALL_TOTAL=0
declare -A FILE_HITS

while IFS= read -r line; do
  if [[ $line == "mode:"* ]]; then
    continue
  fi
  FILE=$(echo "$line" | cut -d':' -f1)
  RANGE=$(echo "$line" | cut -d':' -f2 | awk '{print $1}')
  COUNT=$(echo "$line" | awk '{print $3}')
  if [ "$COUNT" -gt 0 ]; then
    ((OVERALL_COVERED++))
  fi
  ((OVERALL_TOTAL++))
  FILE_NAME=$(basename "$FILE")
  FILE_HITS["$FILE_NAME"]=$((FILE_HITS["$FILE_NAME"] + (COUNT > 0 ? 1 : 0)))
done < "$COVERAGE_OUT"

OVERALL_COVERAGE=$(awk "BEGIN { printf \"%.2f\", ($OVERALL_COVERED/$OVERALL_TOTAL)*100 }")

# === go vet ===
VET_OUTPUT=$(go vet ./... 2>&1)

# === golint ===
LINT_OUTPUT=$(golint ./... 2>&1)

# === gosec ===
gosec ./... > "$GOSEC_REPORT" 2>&1

# === MERMAID COMMIT GRAPH ===
MERMAID_GRAPH="graph TD"
PREV_COMMIT=""
TOTAL_COMMITS=$(git rev-list --count HEAD)  # Get the total number of commits
DISPLAY_COMMITS=5  # Number of commits to display

if [ "$TOTAL_COMMITS" -gt "$DISPLAY_COMMITS" ]; then
  MERMAID_GRAPH+="\n  N[\"... $((TOTAL_COMMITS - DISPLAY_COMMITS)) more commits\"]"
fi

# Fetch the latest 5 commits
git log -n "$DISPLAY_COMMITS" --pretty=format:'%h|%s' | while IFS='|' read -r hash msg; do
  # Escape special characters for Mermaid rendering
  CLEAN_MSG=$(echo "$msg" | sed 's/"/\\"/g' | sed 's/</\\</g' | sed 's/>/\\>/g' | sed 's/`/\\`/g' | sed 's/<br>/ /g')
  MERMAID_GRAPH="${MERMAID_GRAPH}"$'\n'"  $hash[\"$hash: $CLEAN_MSG\"]"
  if [ -n "$PREV_COMMIT" ]; then
    MERMAID_GRAPH="${MERMAID_GRAPH}"$'\n'"  $hash --> $PREV_COMMIT"
  fi
  PREV_COMMIT="$hash"
done

# If there are more commits, link to "N more commits"
if [ "$TOTAL_COMMITS" -gt "$DISPLAY_COMMITS" ]; then
  MERMAID_GRAPH="${MERMAID_GRAPH}"$'\n'"  N --> $PREV_COMMIT"
fi

# === GENERATE MARKDOWN ===
cat << EOF > "$OUTPUT_FILE"
# 📈 Code Coverage Report

## 🔧 Commit Info
- **Branch**: \`$BRANCH\`
- **Commit**: \`$COMMIT\`
- **Author**: $AUTHOR
- **Date**: $DATE
- **Message**: _${MESSAGE}_

---

## 🕸️ Recent Commit History (Mermaid)

\`\`\`mermaid
$MERMAID_GRAPH
\`\`\`

---

## 📝 Changed Files

\`\`\`diff
$(for file in $CHANGED_FILES; do echo "+ $file"; done)
\`\`\`

---

## 🔎 Per-File Coverage Breakdown

| File | Coverage | Bar |
|------|----------|-----|
$(for file in "${!FILE_HITS[@]}"; do
  total_lines=$(grep -c "$file:" "$COVERAGE_OUT")
  covered_lines=${FILE_HITS["$file"]}
  if [ "$total_lines" -eq 0 ]; then
    percent=0
  else
    percent=$(awk "BEGIN { printf \"%.1f\", ($covered_lines/$total_lines)*100 }")
  fi

  bar_len=30
  filled_count=$(awk "BEGIN { printf \"%d\", ($percent/100)*$bar_len }")
  empty_count=$((bar_len - filled_count))
  filled_bar=$(printf "%0.s█" $(seq 1 $filled_count))
  empty_bar=$(printf "%0.s░" $(seq 1 $empty_count))
  bar="$filled_bar$empty_bar"

  echo "| \`$file\` | ${percent}% (${covered_lines}/${total_lines}) | $bar |"
done | sort -n)

---

## 📊 Coverage Summary

**Total Covered Statements**: $OVERALL_COVERED / $OVERALL_TOTAL → **$OVERALL_COVERAGE%**

---

## 🧪 go vet Analysis

\`\`\`
$VET_OUTPUT
\`\`\`

---

## 🧹 golint Results

\`\`\`
$LINT_OUTPUT
\`\`\`

---

## 🔐 gosec Security Report

\`\`\`
$(cat "$GOSEC_REPORT")
\`\`\`

EOF

# === CLEANUP ===
rm -f "$GOSEC_REPORT"

echo "✅ Report generated at $OUTPUT_FILE"
