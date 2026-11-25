# Giving a planning chat a longer memory

An organizer can ask for a checklist, work through a few other details
and then ask what the second step was. Our first implementation sent
only the last two saved turns to the assistant. After a few detours the
checklist could still be on the screen while the assistant could no
longer see it.

For C1, Core walks backwards from the latest saved turn. It keeps at
most eight complete exchanges and stops when adding the next one
would exceed 12,000 UTF-8 bytes with the new prompt included. It then
sends them to Go in chronological order. The signed request names the
actor and organization, and Go checks the same limits independently.
The transcript remains separate from the text sent for each reply.

A summary might let us remember more distant events. It would also
require deciding when to regenerate it, how to handle mistakes and
what deleting the source should do. A bounded window gives us an
easier behavior to reproduce first. A model may still miss a detail
within that window. The local fake demonstrates which evidence
reaches the reply path. It cannot tell us how often a hosted model
would use that evidence correctly.

[Reproduction and review](../evidence/C1.md)
