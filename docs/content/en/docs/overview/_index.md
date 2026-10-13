---
title: "Overview"
linkTitle: "Overview"
weight: 1
description: >
  What N42 Gateway is and how it works.
---

## HAProxy

> [!INFO] TODO
>
> About the HAProxy load balancer itself and why HAProxy.

## Ingress Controllers

> [!INFO] TODO
>
> About the Kubernetes' ingress controllers, what N42 Gateway watches in the k8s cluster and how it builds HAProxy configuration.

## Releases

N42 Gateway is released as a Docker image on top of Alpine's flavor of
[HAProxy image](https://hub.docker.com/_/haproxy).

See the [individual releases](https://github.com/n42-gateway/n42-gateway/releases)  to read the changelog per release.

Images are built by a [GitHub Workflow](https://github.com/n42-gateway/n42-gateway/actions/workflows/image.yaml)
and deployed to [ghcr.io](https://github.com/n42-gateway/n42-gateway/pkgs/container/n42-gateway) and
[quay.io](https://quay.io/repository/n42-gateway/n42-gateway) whenever a tag is applied.
The `latest` tag will always point to the latest stable version while `canary` tag will always
point to the latest beta-quality and release-candidate versions.

Before the beta-quality releases, the source code could also be tagged and images deployed.
The `snapshot` tag will always point to the latest tagged version, which could be a release,
a beta-quality or a development version.

### Signatures and attestations

Every image is signed with [cosign](https://docs.sigstore.dev/cosign/) keyless signing
and comes with signed attestations: a [SLSA](https://slsa.dev/) build provenance and an
[SPDX](https://spdx.dev/) SBOM (Software Bill of Materials). They let you check that
an image was built by the N42 Gateway GitHub Workflow from this repository, which
commit it was built from, and which packages and Go modules it contains, e.g. to scan
it for vulnerabilities or to enforce an admission policy in your cluster. Signatures
and attestations are stored next to the image in both registries, the attestations
on GitHub as well.

Verify the signature, issued to the image workflow running on a tag:

```shell
cosign verify ghcr.io/n42-gateway/n42-gateway:<tag> \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  --certificate-identity-regexp '^https://github\.com/n42-gateway/n42-gateway/\.github/workflows/image\.yaml@refs/tags/'
```

Verify the build provenance attestation, using the [GitHub CLI](https://cli.github.com/):

```shell
gh attestation verify oci://ghcr.io/n42-gateway/n42-gateway:<tag> \
  --repo n42-gateway/n42-gateway
```

Verify the SBOM attestation and print the SBOM, the one of the `linux/amd64` image:

```shell
gh attestation verify oci://ghcr.io/n42-gateway/n42-gateway:<tag> \
  --repo n42-gateway/n42-gateway \
  --predicate-type https://spdx.dev/Document/v2.3 \
  --format json --jq '.[].verificationResult.statement.predicate'
```

Every platform's image also has its own, unsigned, SBOM and build provenance in the
image index, covered by the image signature. Read the SBOM of a platform:

```shell
docker buildx imagetools inspect ghcr.io/n42-gateway/n42-gateway:<tag> \
  --format '{{ json (index .SBOM "linux/amd64").SPDX }}'
```

The same commands apply to `quay.io/n42-gateway/n42-gateway`.

## Where should I go next?

* [Getting Started]({{% relref "/docs/getting-started" %}}): Get started with N42 Gateway!
