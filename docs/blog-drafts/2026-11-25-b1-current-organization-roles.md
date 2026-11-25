# Where the assistant gets the list of organizations

A president can belong to more than one organization. The selected
workspace tells Campusly where the conversation is happening. To answer
“Which organizations do I manage?”, the assistant needs a current list
of president assignments for that account.

The local Go service recognizes that question and sends back an intent
with the request's scope. Core checks that the caller still leads the
selected organization, then loads the active assignments for the caller.
The UI marks the current workspace in the result.

The answer lives in the open chat. A new request after a role change
reads the current assignments, and reopening a saved conversation
restores only its planning transcript. This keeps a useful answer in
view without turning last week's assignment into this week's source.
The [local evidence](../evidence/B1.md) records the cases and limits.
