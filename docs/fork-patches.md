# Fork Patch Stack

This fork keeps local patches as a small linear stack on `bridgev2-main` above
upstream `mautrix/discord:main`.

Do not fold the patches into upstream snapshots or make ad hoc edits in the
deployed checkout. Put every local behavior change in its own commit on `main`
so it can be replayed cleanly when upstream moves.

## Upstream

- Upstream repo: `mautrix/discord`
- Upstream branch: `main`
- Maintained fork branch: `keithah/mautrix-discord:bridgev2-main`
- Sync branch used by automation: `sync/upstream-main`

## Runtime Patches

These are the behavior patches that must continue to apply:

1. Relay webhook/profile handling for Matrix-to-Discord sends.
2. Transient Discord REST and webhook retry/backoff.
3. Discord gateway reconnect throttling.
4. Bot-token session support for the bridge session.
5. Local media capability limit raised to 30 MiB.

Those commits contain:

- Relayed Matrix reactions through a configured Discord session.
- Relayed reaction redaction handling.
- Guardrails for disconnected relay sessions.
- Filtering invalid Matrix reaction keys before sending to Discord.
- Dropping Discord link-preview-only updates to avoid duplicate IRC output via Heisenbridge.
- Avoiding direct bot sends when relay mode should be used.
- Resolving relay webhook names from Matrix or matching ghost profiles instead of the relay bot profile.

The relay reaction patches depend on the matching mautrix-go fork hook until
that hook lands upstream.

## Maintenance Patches

The fork also contains maintenance-only commits, such as this document and the
GitHub Actions workflow. Those commits are intentionally part of the fork stack
so a rebuilt `main` still has the automation.

## Updating Manually

The GitHub Action does this daily. To reproduce locally:

```sh
git switch bridgev2-main
git pull --ff-only origin bridgev2-main
scripts/forward-port-fork.sh
go test ./...
git push origin sync/upstream-main
```

If the result looks good, merge the generated PR into `main`, or fast-forward
`bridgev2-main` to the sync branch:

```sh
git switch bridgev2-main
git merge --ff-only sync/upstream-main
git push origin bridgev2-main
```

If cherry-picking conflicts, resolve the conflict in the sync branch, continue
with `git cherry-pick --continue`, run tests, then update `bridgev2-main`.
