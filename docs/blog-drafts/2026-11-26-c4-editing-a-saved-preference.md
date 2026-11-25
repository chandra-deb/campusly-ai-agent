# Making a saved preference easy to change

A president may want short checklists this week and more detail for a
larger event. The memory panel now shows the saved choice, how it was
saved and when it last changed. The owner can pick a new value,
apply the edit or remove the entry.

When two panels have the same row open, whichever saves second may
still be looking at the older version. The edit request carries the
timestamp it saw. Core compares that timestamp while it holds the
workspace's role lock. If the stored value has already changed, the
panel asks the owner to reload. The row keeps its identity and
original provenance across accepted edits.

In the synthetic chat, a concise checklist has two steps. Editing the
preference to Expanded adds an equipment step in a new chat.
Removing the preference gives the next chat the default checklist.
These fixtures make the effect of an edit visible without requiring a
model to discover the meaning of the settings.

[Reproduction and review](../evidence/C4.md)
