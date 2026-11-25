# Editing a proposed workshop

A workshop request often starts with only a title and a rough capacity.
The local assistant returns editable fields, leaving missing details
for the president. The implementation uses a small typed proposal so the form can show
each field and point out values that need correction.

The browser test shortened the suggested title to one letter, saw the
validation message, and replaced it with the intended title. The Core
test checked that asking for a draft creates no activity and leaves no
draft fields in the saved transcript. See the [B4 checks](../evidence/B4.md).
