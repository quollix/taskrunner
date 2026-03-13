## Change Control

- Default mode is diagnosis-only.
- Never modify files, run formatters that change files, or execute write operations unless explicitly approved by the user in the same turn.
- Before any edit, present:
  1. files to change
  2. brief reason
  3. exact command(s) or patch scope
- Wait for an approval message that contians with `DOIT` before making changes.

### If no approval
- Continue with read-only investigation only.
- Provide findings, root cause, and a proposed patch diff, but do not apply it.

### Safety fallback
- If instructions conflict, this change-control section takes precedence.

### Optional strict mode
- One `DOIT` only authorizes one edit batch; require a new `DOIT` for subsequent edits.
