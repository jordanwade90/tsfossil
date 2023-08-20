# tsfossil

A [Tailscale][ts] [tsnet][] wrapper for the [Fossil][fsl] VCS.

[ts]: https://tailscale.com
[tsnet]: https://pkg.go.dev/tailscale.com/tsnet
[fsl]: https://www.fossil-scm.org

```console
docker build -t tsfossil .
docker run -e TS_AUTHKEY=... -e TS_HOSTNAME=... tsfossil
```

## Configuring the container

| Environment variable | Use                                  | Default value      |
|----------------------|--------------------------------------|--------------------|
| TS_AUTHKEY           | The Tailscale authkey to log in with | *none*             |
| TS_HOSTNAME          | The hostname to log in as            | container hostname |
| TS_STATE_DIR         | Tailscale state directory            | tsnet default      |

## Configuring the repo

The repo is located at `/museum/repo.fossil`.
If no repo exists at the time the first Tailscale user attempts to access it,
a new repo will be created and that user will be made the Setup user of it.
The new repo is configured to use automatic sign-in as the Tailscale username.

If an existing repo is copied into the container,
that repo needs to be configured to use `REMOTE_USER` authentication
to have automatic sign-in.

## Persistence

Make `/museum` a volume to persist the Fossil repo.
Set TS_STATE_DIR and make it a volume to persist the Tailscale state directory.
