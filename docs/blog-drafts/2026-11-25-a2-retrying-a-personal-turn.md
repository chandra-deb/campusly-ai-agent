# What happens when an assistant reply times out?

A timeout used to leave the prompt in the composer with a generic error.
The generic error left the next action unclear. A2 leaves the
prompt in the transcript with **Retry** and **Edit message**. Retry keeps
the same turn ID. Edit returns the text to the composer so the next send
has its own ID.

The browser test held the first synthetic reply until the loading state
was visible, returned a timeout, then sent the retry through Core and
Go. Both requests carried the same ID and text, and one prompt and reply
appeared together. A2 can repeat generation after a lost reply. Keeping a
single visible turn is different from saving a single durable turn, which
is the next slice. See the [A2 evidence](../evidence/A2.md).
