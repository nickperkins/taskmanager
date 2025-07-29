#!/bin/bash

API_URL="http://localhost:8080/tasks"

random_sleep() {
  sleep "$(awk -v min=0.2 -v max=0.6 'BEGIN{srand(); print min+rand()*(max-min)}')"
}

echo "--- Show current tasks (initial) ---"
resp=$(curl -s $API_URL)
echo "$resp" | jq . || echo "[WARN] Invalid or empty response: $resp"
random_sleep

echo "--- Add a task with a user-supplied ID ---"

USER_ID="demo-script-id-001"

USER_ID_PAYLOAD='{"id":"'$USER_ID'","title":"Task with explicit user ID","description":"Created with explicit user-supplied id at step: Add a task with a user-supplied ID","completed":false}'
echo "[curl] curl -X POST -H 'Content-Type: application/json' -d '$USER_ID_PAYLOAD' $API_URL"
echo "[payload] $USER_ID_PAYLOAD"
USER_ID_RESP=$(curl -s -X POST -H "Content-Type: application/json" -d "$USER_ID_PAYLOAD" $API_URL)
echo "$USER_ID_RESP" | jq .
random_sleep

echo "--- Add two new tasks ---"

TASK1_PAYLOAD='{"title":"Task created as first auto-generated","description":"Created as first auto-generated task at step: Add two new tasks","completed":false}'
echo "[curl] curl -X POST -H 'Content-Type: application/json' -d '$TASK1_PAYLOAD' $API_URL"
echo "[payload] $TASK1_PAYLOAD"
TASK1_RESP=$(curl -s -X POST -H "Content-Type: application/json" -d "$TASK1_PAYLOAD" $API_URL)
echo "$TASK1_RESP" | jq .

TASK2_PAYLOAD='{"title":"Task created as second auto-generated","description":"Created as second auto-generated task at step: Add two new tasks","completed":false}'
echo "[curl] curl -X POST -H 'Content-Type: application/json' -d '$TASK2_PAYLOAD' $API_URL"
echo "[payload] $TASK2_PAYLOAD"
TASK2_RESP=$(curl -s -X POST -H "Content-Type: application/json" -d "$TASK2_PAYLOAD" $API_URL)
echo "$TASK2_RESP" | jq .
TASK2_ID=$(echo "$TASK2_RESP" | jq -r '.id')
random_sleep

echo "--- Show current tasks (after adding two) ---"
resp=$(curl -s $API_URL)
echo "$resp" | jq . || echo "[WARN] Invalid or empty response: $resp"
random_sleep

echo "--- Delete the first task ---"

# Dynamically delete the first task in the current list
FIRST_ID=$(curl -s $API_URL | jq -r '.[0].id')
if [ -n "$FIRST_ID" ] && [ "$FIRST_ID" != "null" ]; then
  echo "[curl] curl -X DELETE $API_URL/$FIRST_ID"
  DELETE_RESP=$(curl -s -X DELETE $API_URL/"$FIRST_ID")
  if [ -n "$DELETE_RESP" ]; then
    echo "$DELETE_RESP" | jq .
  else
    echo "Deleted first task with id: $FIRST_ID"
  fi
else
  echo "[WARN] No task found to delete as first task."
fi
random_sleep

echo "--- Add a third task ---"

TASK3_PAYLOAD='{"title":"Task created as third auto-generated","description":"Created as third auto-generated task at step: Add a third task","completed":false}'
echo "[curl] curl -X POST -H 'Content-Type: application/json' -d '$TASK3_PAYLOAD' $API_URL"
echo "[payload] $TASK3_PAYLOAD"
TASK3_RESP=$(curl -s -X POST -H "Content-Type: application/json" -d "$TASK3_PAYLOAD" $API_URL)
echo "$TASK3_RESP" | jq .
random_sleep

echo "--- Update the second task ---"
UPDATE_PAYLOAD='{"title":"Task created as second auto-generated (updated)","description":"Updated at step: Update the second task","completed":true}'
echo "[curl] curl -X PUT -H 'Content-Type: application/json' -d '$UPDATE_PAYLOAD' $API_URL/$TASK2_ID"
echo "[payload] $UPDATE_PAYLOAD"
UPDATE_RESP=$(curl -s -X PUT -H "Content-Type: application/json" -d "$UPDATE_PAYLOAD" $API_URL/$TASK2_ID)
echo "$UPDATE_RESP" | jq .
random_sleep

echo "--- Show current tasks (after update) ---"
resp=$(curl -s $API_URL)
echo "$resp" | jq . || echo "[WARN] Invalid or empty response: $resp"
random_sleep

echo "--- Delete all tasks ---"
# Fetch all current tasks and delete each by ID
ALL_IDS=$(curl -s $API_URL | jq -r '.[].id')
for id in $ALL_IDS; do
  echo "[curl] curl -X DELETE $API_URL/$id"
  DELETE_RESP=$(curl -s -X DELETE $API_URL/$id)
  if [ -n "$DELETE_RESP" ]; then
    echo "$DELETE_RESP" | jq .
  else
    echo "Deleted task with id: $id"
  fi
  sleep 0.2
done

echo "--- Show current tasks (final) ---"
resp=$(curl -s $API_URL)
echo "$resp" | jq . || echo "[WARN] Invalid or empty response: $resp"
