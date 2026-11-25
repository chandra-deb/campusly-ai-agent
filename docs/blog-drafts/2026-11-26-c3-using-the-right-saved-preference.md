# Carrying a saved preference into a new chat

After saving a preference for bullet points, the president should not
have to repeat it in a new planning chat. The first recall path applies
style preferences to a supported workshop checklist request. Core
loads them by the current workspace and president and sends just the
supported choices to Go. A question about the live workshop venue
still uses the activity record.

We could send all private memories with every request. That would
put unrelated context beside live answers and make deletion harder to
reason about. A specific request path gives us a clear first test.
The same checklist is rendered with bullets in the workspace that
saved them and with numbered steps in the other workspace.

A save or removal can happen while the model is busy. Core releases
its lock before calling Go, then checks the same choices again before
it commits a styled reply. If a preference was removed in the
meantime, the app asks the user to retry. This protects new replies.
The previous conversation remains readable until its owner deletes it.

[Reproduction and review](../evidence/C3.md)
