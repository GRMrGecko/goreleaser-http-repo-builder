# goreleaser-http-repo-builder

This tool was written out of the need to build a release repository compatible with [go-selfupdate](https://github.com/creativeprojects/go-selfupdate) with releases built by [goreleaser](https://goreleaser.com/).

## Example Usage

The command has extensive help available, the following is an example of building a release and adding it to a new repo.

```bash
goreleaser release --snapshot --skip=publish
mkdir repo
goreleaser-http-repo-builder add-release --repo=repo/ --release=dist/
```

After adding a release, you can copy the repo to your web server for update distrobution.

## Example Goreleaser Config

While there is good [documentation available](https://goreleaser.com/customization/) that I'd recommend reading, the following provides some examples that may be helpful in generating a release that is compatible with go-selfupdate.

- The checksums file name defaults to preappend the project name, which is not compatible if you wish to use the checksums to verify an update.
- If you're signing releases with an ECDSA key, this is what I found works best.
- If you need to specify the version manually, you can edit the version template. By default, goreleaser will use the git tag to determine the version.

```yaml
version: 2

before:
  hooks:
    - go mod tidy

builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    name_template: "{{ .ProjectName }}_{{ .Os }}_{{ .Arch }}"
    wrap_in_directory: true

checksum:
    name_template: "checksums.txt"

signs:
  - artifacts: all
    cmd: openssl
    args:
        - dgst
        - -sha256
        - -sign
        - "signing.key"
        - -out
        - ${signature}
        - ${artifact}

snapshot:
    version_template: "v0.1.2"
```
