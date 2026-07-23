#!/usr/bin/env bash
set -euo pipefail

API="${API:-http://localhost:8080}"

echo "=== DAWAI E2E Smoke Test ==="

# 1. Login
echo "[1] Login..."
TOKEN=$(curl -s -X POST "$API/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@test.local","password":"TestPass123!"}' \
  | jq -r '.data.accessToken // empty')

if [ -z "$TOKEN" ]; then echo "FAIL: login (API unreachable or creds bad)"; exit 1; fi
echo "  PASS: got JWT"

# 2. Create subject
echo "[2] Create subject..."
SUBJECT_ID=$(curl -s -X POST "$API/api/subjects" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"Violin 101","description":"Beginner violin"}' \
  | jq -r '.data.id // empty')
[ -z "$SUBJECT_ID" ] && { echo "FAIL: create subject"; exit 1; }
echo "  PASS: subject $SUBJECT_ID"

# 3. Rubric component
echo "[3] Add rubric component..."
COMP_ID=$(curl -s -X POST "$API/api/subjects/$SUBJECT_ID/rubric" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"name":"Intonation","description":"Pitch accuracy","scale_min":1,"scale_max":5,"weight":1}' \
  | jq -r '.data.id // empty')
[ -z "$COMP_ID" ] && { echo "FAIL: add rubric"; exit 1; }
echo "  PASS: rubric $COMP_ID"

# 4. List students
echo "[4] List students..."
STUDENT_IDS=$(curl -s "$API/api/students" -H "Authorization: Bearer $TOKEN" | jq -r '.data[].id')
COUNT=$(echo "$STUDENT_IDS" | grep -c . || true)
echo "  INFO: $COUNT students in school"

# 5. Submit assessment (skip if no students seeded)
if [ "$COUNT" -gt 0 ]; then
  STUDENT_ID=$(echo "$STUDENT_IDS" | head -1)
  echo "[5] Submit assessment for $STUDENT_ID..."
  ASSESS_ID=$(curl -s -X POST "$API/api/assessments" \
    -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
    -d "{\"student_id\":\"$STUDENT_ID\",\"subject_id\":\"$SUBJECT_ID\",\"scores\":[{\"rubric_component_id\":\"$COMP_ID\",\"score\":4}],\"feedback\":\"Good progress\"}" \
    | jq -r '.data.id // empty')
  [ -z "$ASSESS_ID" ] && { echo "FAIL: submit assessment"; exit 1; }
  echo "  PASS: assessment $ASSESS_ID"

  echo "[6] Fetch assessment..."
  RESULT=$(curl -s "$API/api/assessments/$ASSESS_ID" -H "Authorization: Bearer $TOKEN" | jq -r '.success')
  [ "$RESULT" = "true" ] && echo "  PASS: fetched" || { echo "FAIL: fetch"; exit 1; }
else
  echo "[5-6] SKIP: no students seeded (create via POST /api/students first)"
fi

echo ""
echo "=== SMOKE COMPLETE (happy path OK) ==="
