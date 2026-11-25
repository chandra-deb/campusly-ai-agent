# A private chat that survives a reload

A1 and A2 kept the transcript in the page. Reloading meant starting over.
A3 saves each successful prompt and reply in a Core conversation owned
by the current president and organization. A new president for the same
organization cannot read the previous president's transcript.

The request reads the two most recent turns and releases the database
transaction before asking Go for a reply. On return, Core rechecks
presidency and locks the conversation long enough to save the turn.
If the same turn was already saved, the retry gets the original answer.
If someone deleted the conversation while generation was running, the
late reply cannot recreate it.

The owner chose retention until deletion for this synthetic lab. A
foreign key removes the turns when the parent conversation is deleted.
Long-term memory, backups and third-party model transfer are separate
design questions. See the [A3 retention decision](../decisions/0005-personal-conversation-retention.md).
