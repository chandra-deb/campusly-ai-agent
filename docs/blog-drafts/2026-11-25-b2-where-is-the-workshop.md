# Finding the next workshop from its schedule

The activity title alone cannot tell us which workshop comes next.
A workshop can have several sessions, or its first session may have
been cancelled. The assistant asks Core for the next scheduled session
in the selected organization. The answer shows exactly where the
record has a venue, with a button to open the activity.

The test uses two organizations. The other one has an earlier
workshop, and the current one has a cancelled session before Room 12.
The answer has to skip both. A named activity without a venue gets
a clear missing-data message. The [B2 checks](../evidence/B2.md) cover
those branches and show the resulting chat.
