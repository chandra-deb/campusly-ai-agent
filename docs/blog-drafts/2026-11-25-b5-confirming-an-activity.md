# Keeping one activity when confirmation is retried

A president reviews the title, description, type, capacity and
organization before creating an activity. Core signs those exact fields
and checks the current president assignment again when the user confirms.
An edited value cannot use the earlier signature.

The activity, discussion and receipt commit together. If the request
is repeated after the commit, Core can return the same activity ID.
The test also injected a failure after the discussion was flushed.
The write rolled back, and a retry using the same approval succeeded
after the synthetic failure was removed. See the [local checks](../evidence/B5.md).
