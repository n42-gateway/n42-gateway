# CHANGELOG v0.17 branch

* [Major improvements](#major-improvements)
* [Upgrade notes - read before upgrade from v0.16!](#upgrade-notes)
* [Contributors](#contributors)
* [v0.17.0-alpha.3](#v0170-alpha3)
  * [Reference](#reference-a3)
  * [Release notes](#release-notes-a3)
  * [Improvements](#improvements-a3)
  * [Fixes](#fixes-a3)
* [v0.17.0-alpha.2](#v0170-alpha2)
  * [Reference](#reference-a2)
  * [Release notes](#release-notes-a2)
  * [Improvements](#improvements-a2)
  * [Fixes](#fixes-a2)
* [v0.17.0-alpha.1](#v0170-alpha1)
  * [Reference](#reference-a1)
  * [Release notes](#release-notes-a1)
  * [Improvements](#improvements-a1)
  * [Fixes](#fixes-a1)

## Major improvements

Highlights of this version:

* Branding changed from HAProxy Ingress to N42 Gateway 
* Embedded HAProxy version update from 2.8 to 3.0.
* Gateway API compliant implementation. Support of HTTPRoute, TLSRoute and TCPRoute APIs, although the last one does not have compliance tests yet.
* Support of multiple HTTP(S) frontends on distinct TCP port number.
* Use of `add server` and `del server` API calls during scale-out and scale-in operations.

## Upgrade notes

Breaking backward compatibility from v0.16:

* The branding update need a special attention when upgrading to `v0.17.0-alpha.3` and newer, see the [migration guide](https://n42-gateway.github.io/v0.17/docs/migration-guide/).
* HAProxy versions older than 2.6 are no longer supported in the External HAProxy deployment.
* Gateway API updated from v1.0 to v1.6, which drops support to Gateway and HTTPRoute v1alpha2.
* The default HTTP and HTTPS frontends are created only if an Ingress or Gateway API resource references it. The old behavior of always having HTTP(S) binding their TCP ports can be achieved using [`create-default-frontends`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#bind-port) configuration key.
* All the CORS response headers are removed if Origin request header is missing or not authorized. This improves compliance and simplifies the configuration.
* Plain HTTP Passthrough, formerly known as Fronting Proxy, was redesigned for the multiple HTTP(S) frontends support, and its configuration was simplified. Give it a special attention during tests and check the new documentation: [HTTP Passthrough](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#http-passthrough).

A refactor was made on internal HAProxy model to support multiple HTTP(S) frontends, this is the biggest refactor since the v0.8 one in the converter code. Although there are no known backward compatibility changes beyond the ones already reported, it is suggested to thoroughly observe v0.17 on test and staging environments before migrate to production. Do not hesitate to file an issue if you find a misbehavior.

## Contributors

* Alexander Stephan ([alexanderstephan](https://github.com/alexanderstephan))
* Arthur Le Roux ([arthlr](https://github.com/arthlr))
* Ask Bjørn Hansen ([abh](https://github.com/abh))
* Bagas Purwa S ([bapung](https://github.com/bapung))
* Christian Menges ([Garfield96](https://github.com/Garfield96))
* CodeOpsAI ([CodeOpsAI](https://github.com/CodeOpsAI))
* hayden ([onelapahead](https://github.com/onelapahead))
* Ian Roberts ([ianroberts](https://github.com/ianroberts))
* Joao Morais ([jcmoraisjr](https://github.com/jcmoraisjr))
* Josh Soref ([jsoref](https://github.com/jsoref))
* Lola Delannoy ([spnngl](https://github.com/spnngl))
* Mia Mouret ([mia-mouret](https://github.com/mia-mouret))
* Nadia Santalla ([nadiamoe](https://github.com/nadiamoe))
* Pedro Gonçalves ([PerGon](https://github.com/PerGon))
* Silvio Knizek ([killermoehre](https://github.com/killermoehre))
* Vladimir Kozhukalov ([kozhukalov](https://github.com/kozhukalov))

# v0.17.0-alpha.3

## Reference (a3)

* Release date: `2026-09-16`
* Helm chart: `--version 0.17.0-alpha.3 --devel`
* Image (GHCR): `ghcr.io/n42-gateway/n42-gateway:v0.17.0-alpha.3`
* Image (Quay): `quay.io/n42-gateway/n42-gateway:v0.17.0-alpha.3`
* Embedded HAProxy version: `3.0.27`
* GitHub release: `https://github.com/n42-gateway/n42-gateway/releases/tag/v0.17.0-alpha.3`

## Release notes (a3)

This is the third tag of the v0.17 branch, which brings a major update:

- Branding changed from HAProxy Ingress to N42 Gateway. This update impacts new deployments. The [migration guide](https://n42-gateway.github.io/v0.17/docs/migration-guide/) covers some upgrade scenarios, with or without backward compatibility, with or without high availability.

Some other notable features or improvements were added:

- New `--full-controller-name` and `--configuration-class` command-line options help to migrate from HAProxy Ingress without breaking backward compatibility.
- Christian added a controller-runtime configuration that strips managed fields from the cached objects, which reduces memory usage.
- hayden added enforcement mode configuration keys, allowing to configure distinct responses on allow and deny lists. [doc](https://n42-gateway.github.io/v0.17/docs/configuration/keys/#allowlist).
- Alexander added a new configuration key, which performs backend server rename on dynamic scaling using pre-allocated slots and naming as `IP` or `POD`. This option requires HAProxy 3.5 dev2 or newer. [doc](https://n42-gateway.github.io/v0.17/docs/configuration/keys/#backend-server-naming).

Some fixes were also merged:

- A number of CVEs from Go's stdlib and N42 Gateway dependencies.
- Arthur reported and fixed the repopulation of the IP address on Ingress and Gateway status when a single replica looses the leader election.
- CodeOpsAI fixed a missing permission on the close issue scheduler.
- Christian found the source of a race in the test code and reenabled the `-race` command-line option in the test run.
- Silvio fixed the rendering of a link in the home page.

Dependencies:

- embedded haproxy from 3.0.23 to 3.0.27
- client-go from v0.36.1 to v0.37.0
- go from 1.26.3 to 1.26.8

## Improvements (a3)

New features and improvements since `v0.17.0-alpha.2`:

* Bump golang.org/x/crypto from 0.51.0 to 0.52.0 [#1484](https://github.com/n42-gateway/n42-gateway/pull/1484) (dependabot[bot])
* add full controller name command-line option [#1487](https://github.com/n42-gateway/n42-gateway/pull/1487) (jcmoraisjr)
  * Command-line options:
    * [`--full-controller-name`](https://n42-gateway.github.io/v0.17/docs/configuration/command-line/#ingress-class)
* Bump golang.org/x/crypto from 0.52.0 to 0.53.0 [#1489](https://github.com/n42-gateway/n42-gateway/pull/1489) (dependabot[bot])
* Bump actions/setup-go from 6 to 7 [#1504](https://github.com/n42-gateway/n42-gateway/pull/1504) (dependabot[bot])
* Bump actions/stale from 10 to 11 [#1510](https://github.com/n42-gateway/n42-gateway/pull/1510) (dependabot[bot])
* Bump k8s.io/client-go from 0.36.1 to 0.36.3 [#1507](https://github.com/n42-gateway/n42-gateway/pull/1507) (dependabot[bot])
* Bump github.com/go-logr/logr from 1.4.3 to 1.4.4 [#1505](https://github.com/n42-gateway/n42-gateway/pull/1505) (dependabot[bot])
* Bump golang.org/x/crypto from 0.53.0 to 0.55.0 [#1502](https://github.com/n42-gateway/n42-gateway/pull/1502) (dependabot[bot])
* Bump sigs.k8s.io/gateway-api from 1.5.1 to 1.6.1 [#1508](https://github.com/n42-gateway/n42-gateway/pull/1508) (dependabot[bot])
* Bump github.com/prometheus/client_golang from 1.23.2 to 1.24.1 [#1514](https://github.com/n42-gateway/n42-gateway/pull/1514) (dependabot[bot])
* Bump github.com/stretchr/testify from 1.11.1 to 1.12.0 [#1515](https://github.com/n42-gateway/n42-gateway/pull/1515) (dependabot[bot])
* Bump k8s.io/client-go from 0.36.3 to 0.36.4 [#1516](https://github.com/n42-gateway/n42-gateway/pull/1516) (dependabot[bot])
* Bump github.com/stretchr/testify from 1.12.0 to 1.12.1 [#1517](https://github.com/n42-gateway/n42-gateway/pull/1517) (dependabot[bot])
* N42 Gateway branding - Step 1/2 [#1523](https://github.com/n42-gateway/n42-gateway/pull/1523) (jcmoraisjr)
* Strip managed fields from cached objects [#1526](https://github.com/n42-gateway/n42-gateway/pull/1526) (Garfield96)
* Add configuration class command-line option [#1527](https://github.com/n42-gateway/n42-gateway/pull/1527) (jcmoraisjr)
  * Command-line options:
    * [`--configuration-class`](https://n42-gateway.github.io/v0.17/docs/configuration/command-line/#configuration-class)
* Enforcement mode for configuring either (default) deny/403 behavior, reject, or the silent-drop behavior like WAF [#1512](https://github.com/n42-gateway/n42-gateway/pull/1512) (onelapahead)
  * Configuration keys:
    * [`allowlist-enforcement-mode`](https://n42-gateway.github.io/v0.17/docs/configuration/keys/#allowlist)
    * [`denylist-enforcement-mode`](https://n42-gateway.github.io/v0.17/docs/configuration/keys/#allowlist)
* Bump k8s.io/client-go from 0.36.4 to 0.37.0 [#1521](https://github.com/n42-gateway/n42-gateway/pull/1521) (dependabot[bot])
* Make server rename part of dynamic updates [#1481](https://github.com/n42-gateway/n42-gateway/pull/1481) (alexanderstephan)
  * Configuration keys:
    * [`backend-server-rename`](https://n42-gateway.github.io/v0.17/docs/configuration/keys/#backend-server-naming)
* N42 Gateway branding - Step 2/2 [#1528](https://github.com/n42-gateway/n42-gateway/pull/1528) (jcmoraisjr)
* Bump golang.org/x/crypto from 0.55.0 to 0.57.0 [#1533](https://github.com/n42-gateway/n42-gateway/pull/1533) (dependabot[bot])
* exercise the new image hub credentials [88f2c59](https://github.com/n42-gateway/n42-gateway/commit/88f2c59707716ce0fd0bc7019d88eed2191e73f4) (Joao Morais)
* bump go from 1.26.3 to 1.26.8 [0a209a4](https://github.com/n42-gateway/n42-gateway/commit/0a209a47f25457e32d09be383435741ee71f55d4) (Joao Morais)
* bump haproxy from 3.0.23 to 3.0.27 [ce38468](https://github.com/n42-gateway/n42-gateway/commit/ce3846848e7862d14590b6b518e9341311f73f62) (Joao Morais)
* bump dependencies [404a62f](https://github.com/n42-gateway/n42-gateway/commit/404a62f85f62f3923060b806fb7c7a286c56bed7) (Joao Morais)

Chart improvements since `v0.17.0-alpha.2`:

* N42 Gateway branding [#118](https://github.com/n42-gateway/charts/pull/118) (jcmoraisjr)

## Fixes (a3)

* Fix not rendered markdown [#1480](https://github.com/n42-gateway/n42-gateway/pull/1480) (killermoehre)
* fix: repopulate ingress status after leader re-acquisition [#1486](https://github.com/n42-gateway/n42-gateway/pull/1486) (arthlr)
* Fix permissions-missing in closeissues.yaml [#1498](https://github.com/n42-gateway/n42-gateway/pull/1498) (CodeOpsAI)
* Fix data race in TestHAProxyProcsLoop and enable -race flag for tests [#1525](https://github.com/n42-gateway/n42-gateway/pull/1525) (Garfield96)
* fix ghcr authentication [4b52fab](https://github.com/n42-gateway/n42-gateway/commit/4b52fab2f1d55c76b0b655f081f602c89c8d374e) (Joao Morais)

# v0.17.0-alpha.2

## Reference (a2)

* Release date: `2026-05-17`
* Helm chart: `--version 0.17.0-alpha.2 --devel`
* Image (Quay): `quay.io/jcmoraisjr/haproxy-ingress:v0.17.0-alpha.2`
* Image (Docker Hub): `docker.io/jcmoraisjr/haproxy-ingress:v0.17.0-alpha.2`
* Embedded HAProxy version: `3.0.23`
* GitHub release: `https://github.com/jcmoraisjr/haproxy-ingress/releases/tag/v0.17.0-alpha.2`

## Release notes (a2)

This is the second tag of the v0.17 branch, which brings some new features:

- Ian changed the image build in a way it cross-compile to the target platform, from the build one. This change makes the image generation lightning fast: image build time decreased from more than 1 hour to about 5 minutes.
- Added `--watch-ingress` command-line option, allowing to fully disable Ingress API watch and parsing.
- Ian configured all the writable folders as emptyDir, which makes the controller to work on containers having read only file system.
- Lola added VPA (VerticalPodAutoscaler) configuration option.

A number of fixes were also merged, most of them already reported on v0.16 release:

- Bjørn reported and fixed a frontend name matching issue on v0.17. If any of the hosts use ssl-passthrough, paths scoped configurations fail to match if the match applies only on some paths of a backend. This issue was introduced on the new multi frontends: the frontend name now comes from the configuration, instead of hardcoded, and the configuration didn't cover the ssl-passthrough naming.
- Gabriel reported a backend leak issue that can potentially consume all the internal port numbers used by external authentication. This leak only happens on DNS based backends used by http and https protos, and only when the addresses behind the DNS use to change - e.g. using https://somesvc.somenamespace.svc from a headless service, the IP addresses will change whenever the pods behind that service scales.
- Updating base image and Go, which fixes a number of reported CVEs on OS libraries and Go's stdlib.
- Nadia reported that external authentication, if placed in the frontend via `auth-external-placement` configuration key, uses exact path match despite of the path configuration. Backend placed configuration (the default placement) does not have this problem. It is recommended to update HAProxy Ingress asap if you use external authentication placed in the frontend.
- Florian reported that idle metric collector can crash the controller if haproxy eventually reports more than 100 on its metric. This happens because the controller did not check the boundaries and a counter metric would become negative, making Prometheus client to crash. See also https://github.com/haproxy/haproxy/issues/3339.
- A race can happen in the controller start using master-worker mode, when checking if the master socket is already available. In case of an error reading the socket, the controller checks its presence, returning the original error if it was created in this time frame and was found. This behavior makes the controller delay 30 extra seconds to become ready.
- Nadia reported and fixed a case-sensitive match in external authentication header match, which is expected to be case-insensitive.
- Ian reported and fixed the starting user configured in the controller image, changing from the `haproxy` name to its UID `99`. The UID continues the same, but not using the name allows to configure the container runtime to run as non root without the need to specify an UID.
- Logan reported that a PDB resource was always being created despite of being configured, this happened because chart was comparing the `maxUnavailable` to a declared zero, which is also the value when it is not configured.
- The service account is now configured only in the controller container for security reasons, sidecar containers does not have the service account anymore.

Dependencies:

- embedded haproxy from 3.0.18 to 3.0.23
- client-go from v0.35.2 to v0.36.1
- controller-runtime from v0.23.3 to v0.24.1
- go from 1.26.1 to 1.26.3

## Improvements (a2)

New features and improvements since `v0.17.0-alpha.1`:

* Always run `go build` on the BUILDPLATFORM architecture [#1434](https://github.com/jcmoraisjr/haproxy-ingress/pull/1434) (ianroberts)
* update metrics page and dashboard [#1435](https://github.com/jcmoraisjr/haproxy-ingress/pull/1435) (jcmoraisjr)
* Bump k8s.io/client-go from 0.35.2 to 0.35.3 [#1440](https://github.com/jcmoraisjr/haproxy-ingress/pull/1440) (dependabot[bot])
* Bump golang.org/x/crypto from 0.49.0 to 0.50.0 [#1449](https://github.com/jcmoraisjr/haproxy-ingress/pull/1449) (dependabot[bot])
* Bump k8s.io/client-go from 0.35.3 to 0.35.4 [#1454](https://github.com/jcmoraisjr/haproxy-ingress/pull/1454) (dependabot[bot])
* doc: fix default value for watch gateway [#1457](https://github.com/jcmoraisjr/haproxy-ingress/pull/1457) (jcmoraisjr)
* fail initialization if the discovery client fails [#1458](https://github.com/jcmoraisjr/haproxy-ingress/pull/1458) (jcmoraisjr)
* start decoupling global and ingress converters [#1459](https://github.com/jcmoraisjr/haproxy-ingress/pull/1459) (jcmoraisjr)
* Bump k8s.io/client-go from 0.35.4 to 0.36.0 [#1463](https://github.com/jcmoraisjr/haproxy-ingress/pull/1463) (dependabot[bot])
* pin dependencies from makefile [#1467](https://github.com/jcmoraisjr/haproxy-ingress/pull/1467) (jcmoraisjr)
* add watch-ingress command line option [#1460](https://github.com/jcmoraisjr/haproxy-ingress/pull/1460) (jcmoraisjr)
  * Command-line options:
    * [`--watch-ingress`](https://haproxy-ingress.github.io/v0.17/docs/configuration/command-line/#watch-ingress)
* add ip mode command-line option [#1466](https://github.com/jcmoraisjr/haproxy-ingress/pull/1466) (jcmoraisjr)
  * Command-line options:
    * [`--ip-mode`](https://haproxy-ingress.github.io/v0.17/docs/configuration/command-line/#ip-mode)
* add build image step in build workflow [#1442](https://github.com/jcmoraisjr/haproxy-ingress/pull/1442) (jcmoraisjr)
* delay test of image build [#1477](https://github.com/jcmoraisjr/haproxy-ingress/pull/1477) (jcmoraisjr)
* Bump sigs.k8s.io/controller-runtime from 0.23.1-0.20260424122448-c8b4b9d61fbd to 0.24.1 [#1474](https://github.com/jcmoraisjr/haproxy-ingress/pull/1474) (dependabot[bot])
* Bump go.uber.org/zap from 1.27.1 to 1.28.0 [#1475](https://github.com/jcmoraisjr/haproxy-ingress/pull/1475) (dependabot[bot])
* Bump golang.org/x/crypto from 0.50.0 to 0.51.0 [#1476](https://github.com/jcmoraisjr/haproxy-ingress/pull/1476) (dependabot[bot])
* docs: fix watch-ingress entry [c5b5e97](https://github.com/jcmoraisjr/haproxy-ingress/commit/c5b5e97fa897fafa3183f1caf9ee52b3ebcc5d2b) (Joao Morais)
* bump embedded haproxy from 3.0.18 to 3.0.23 [7199d7b](https://github.com/jcmoraisjr/haproxy-ingress/commit/7199d7be3f920acf42aa3e36f230b5d2f721cd2d) (Joao Morais)
* bump go from 1.26.1 to 1.26.3 [b9a6519](https://github.com/jcmoraisjr/haproxy-ingress/commit/b9a6519a1d0518cc30507423d03c3d0ca9d1e47b) (Joao Morais)
* bump client-go from v0.36.0. to v0.36.1 [5d4a598](https://github.com/jcmoraisjr/haproxy-ingress/commit/5d4a59879eb3c9fa6289bbf3317206ed0754df87) (Joao Morais)

Chart improvements since `v0.17.0-alpha.1`:

* feat: always mount writeable folders from emptyDir [#106](https://github.com/haproxy-ingress/charts/pull/106) (ianroberts)
* feat: add VerticalPodAutoscaler resource for controller [#107](https://github.com/haproxy-ingress/charts/pull/107) (spnngl)
* create pdb only if max unavailable is defined [#108](https://github.com/haproxy-ingress/charts/pull/108) (jcmoraisjr)
* manually mount sa only in controller pod [#109](https://github.com/haproxy-ingress/charts/pull/109) (jcmoraisjr)

## Fixes (a2)

* Use numeric USER in Dockerfile [#1431](https://github.com/jcmoraisjr/haproxy-ingress/pull/1431) (ianroberts)
* match sni using single dns label [#1436](https://github.com/jcmoraisjr/haproxy-ingress/pull/1436) (jcmoraisjr)
* fix HTTPS frontend name mismatch breaking backend pathID matching [#1433](https://github.com/jcmoraisjr/haproxy-ingress/pull/1433) (abh)
* convert user-provided auth external header names to lowercase [#1429](https://github.com/jcmoraisjr/haproxy-ingress/pull/1429) (nadiamoe)
* add tracking for http based auth backends [#1462](https://github.com/jcmoraisjr/haproxy-ingress/pull/1462) (jcmoraisjr)
* parameterize the eventually timeout and interval [#1468](https://github.com/jcmoraisjr/haproxy-ingress/pull/1468) (jcmoraisjr)
* adding boundary in the idle_pct metric [#1456](https://github.com/jcmoraisjr/haproxy-ingress/pull/1456) (jcmoraisjr)
* fix ip address parsing for ip mode [#1472](https://github.com/jcmoraisjr/haproxy-ingress/pull/1472) (jcmoraisjr)
* fix request match on frontend based external auth [#1470](https://github.com/jcmoraisjr/haproxy-ingress/pull/1470) (jcmoraisjr)
* fix race checking if haproxy socket is missing [#1471](https://github.com/jcmoraisjr/haproxy-ingress/pull/1471) (jcmoraisjr)

# v0.17.0-alpha.1

## Reference (a1)

* Release date: `2026-03-15`
* Helm chart: `--version 0.17.0-alpha.1 --devel`
* Image (Quay): `quay.io/jcmoraisjr/haproxy-ingress:v0.17.0-alpha.1`
* Image (Docker Hub): `docker.io/jcmoraisjr/haproxy-ingress:v0.17.0-alpha.1`
* Embedded HAProxy version: `3.0.18`
* GitHub release: `https://github.com/jcmoraisjr/haproxy-ingress/releases/tag/v0.17.0-alpha.1`

## Release notes (a1)

This is the first tag of the v0.17 branch, which brings a number of new features:

* Embedded HAProxy version updated from 2.8 to 3.0. We always upgrade HAProxy to the next minor, non EOL version on every new minor version of HAProxy Ingress. Note also that the HAProxy version can be easily changed using external deployment. See the [external deployment documentation](https://haproxy-ingress.github.io/v0.17/docs/examples/external-haproxy/).
* This is the first release to reach Gateway API compliance. HAProxy Ingress implements all core and some extended features of HTTPRoute and TLSRoute. TCPRoute is also supported, however Gateway API does not have compliance tests for this API yet.
* Multiple HTTP and HTTPS ports can be configured at the same time, their hostnames and paths don't conflict each other. This is a trivial configuration on Gateway API, and uses annotations to reach the same behavior on Ingress API. Read more about configurable listening port via Ingress API in the [HTTP Frontends configuration keys documentation](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#http-frontends).
* Use the `add server` and `del server` API calls to scale backend server replicas. This new implementation not only avoids the need of pre allocated slots, but also uses the correct server name despite how the backend server name is configured. Read more about the new dynamic update in the new [Dynamic Scaling documentation](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#dynamic-scaling).

There are some known breaking changes since v0.16, see the [upgrade notes](#upgrade-notes) section in this changelog. We should have at least one other alpha/snapshot release, and we will report any new breaking changes on its release notes in the case it happens.

Dependencies:

- embedded haproxy from 2.8.16 to 3.0.18
- client-go from v0.34.1 to v0.35.2
- controller-runtime from v0.22.3 to v0.23.3
- go from 1.25.3 to 1.26.1

## Improvements (a1)

New features and improvements since `v0.16`:

* configure more than one http(s) ports [#1295](https://github.com/jcmoraisjr/haproxy-ingress/pull/1295) (jcmoraisjr)
* update haproxy from 2.8.16 to 3.0.12 [#1313](https://github.com/jcmoraisjr/haproxy-ingress/pull/1313) (jcmoraisjr)
* Spelling [#1321](https://github.com/jcmoraisjr/haproxy-ingress/pull/1321) (jsoref)
* Bump golangci/golangci-lint-action from 7 to 9 [#1317](https://github.com/jcmoraisjr/haproxy-ingress/pull/1317) (dependabot[bot])
* Bump sigs.k8s.io/controller-runtime from 0.22.3 to 0.22.4 [#1315](https://github.com/jcmoraisjr/haproxy-ingress/pull/1315) (dependabot[bot])
* Bump golang.org/x/sync from 0.17.0 to 0.18.0 [#1316](https://github.com/jcmoraisjr/haproxy-ingress/pull/1316) (dependabot[bot])
* Bump golang.org/x/crypto from 0.43.0 to 0.44.0 [#1320](https://github.com/jcmoraisjr/haproxy-ingress/pull/1320) (dependabot[bot])
* Bump k8s.io/client-go from 0.34.1 to 0.34.2 [#1319](https://github.com/jcmoraisjr/haproxy-ingress/pull/1319) (dependabot[bot])
* Bump golang.org/x/crypto from 0.44.0 to 0.45.0 [#1322](https://github.com/jcmoraisjr/haproxy-ingress/pull/1322) (dependabot[bot])
* Add check-spelling v0.0.25 [#1323](https://github.com/jcmoraisjr/haproxy-ingress/pull/1323) (jsoref)
* add http-ports-local config for http(s) port override [#1314](https://github.com/jcmoraisjr/haproxy-ingress/pull/1314) (jcmoraisjr)
* Bump go.uber.org/zap from 1.27.0 to 1.27.1 [#1324](https://github.com/jcmoraisjr/haproxy-ingress/pull/1324) (dependabot[bot])
* Bump actions/checkout from 5 to 6 [#1325](https://github.com/jcmoraisjr/haproxy-ingress/pull/1325) (dependabot[bot])
* update gateway api from v1.0.0 to v1.4.0 [#1326](https://github.com/jcmoraisjr/haproxy-ingress/pull/1326) (jcmoraisjr)
* move gateway api version check to the cache [#1329](https://github.com/jcmoraisjr/haproxy-ingress/pull/1329) (jcmoraisjr)
* make status update synchronous [#1330](https://github.com/jcmoraisjr/haproxy-ingress/pull/1330) (jcmoraisjr)
* add custom status response backend [#1332](https://github.com/jcmoraisjr/haproxy-ingress/pull/1332) (jcmoraisjr)
* add gateway status update [#1333](https://github.com/jcmoraisjr/haproxy-ingress/pull/1333) (jcmoraisjr)
* drop apiserver dependency [#1336](https://github.com/jcmoraisjr/haproxy-ingress/pull/1336) (jcmoraisjr)
* Bump golang.org/x/sync from 0.18.0 to 0.19.0 [#1337](https://github.com/jcmoraisjr/haproxy-ingress/pull/1337) (dependabot[bot])
* Bump golang.org/x/crypto from 0.45.0 to 0.46.0 [#1338](https://github.com/jcmoraisjr/haproxy-ingress/pull/1338) (dependabot[bot])
* add http passthrough config keys [#1340](https://github.com/jcmoraisjr/haproxy-ingress/pull/1340) (jcmoraisjr)
  * Configuration keys:
    * [`bind-http-passthrough`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#bind)
    * [`http-passthrough`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#http-passthrough)
    * [`http-passthrough-port`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#http-passthrough)
* add option to always create default frontends [#1342](https://github.com/jcmoraisjr/haproxy-ingress/pull/1342) (jcmoraisjr)
  * Configuration keys:
    * [`create-default-frontends`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#bind-port)
* Bump k8s.io/client-go from 0.34.2 to 0.34.3 [#1347](https://github.com/jcmoraisjr/haproxy-ingress/pull/1347) (dependabot[bot])
* Add connection timeout command-line option [#1348](https://github.com/jcmoraisjr/haproxy-ingress/pull/1348) (jcmoraisjr)
  * Command-line options:
    * [`--connection-timeout`](https://haproxy-ingress.github.io/v0.17/docs/configuration/command-line/#timeout)
* update gateway api crds on integration test [#1344](https://github.com/jcmoraisjr/haproxy-ingress/pull/1344) (jcmoraisjr)
* add custom http frontend config keys [#1341](https://github.com/jcmoraisjr/haproxy-ingress/pull/1341) (jcmoraisjr)
  * Configuration keys:
    * [`allow-local-bind`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#bind)
    * [`http-frontend`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#http-frontends)
    * [`http-frontends`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#http-frontends)
* add reload with websocket connection test [#1349](https://github.com/jcmoraisjr/haproxy-ingress/pull/1349) (jcmoraisjr)
* Bump k8s.io/client-go from 0.34.3 to 0.35.0 [#1352](https://github.com/jcmoraisjr/haproxy-ingress/pull/1352) (dependabot[bot])
* Bump golang.org/x/crypto from 0.46.0 to 0.47.0 [#1358](https://github.com/jcmoraisjr/haproxy-ingress/pull/1358) (dependabot[bot])
* Bump golang.org/x/crypto from 0.47.0 to 0.48.0 [#1368](https://github.com/jcmoraisjr/haproxy-ingress/pull/1368) (dependabot[bot])
* Bump k8s.io/client-go from 0.35.0 to 0.35.1 [#1370](https://github.com/jcmoraisjr/haproxy-ingress/pull/1370) (dependabot[bot])
* configure default backend per host [#1376](https://github.com/jcmoraisjr/haproxy-ingress/pull/1376) (jcmoraisjr)
* add header handling config keys [#1374](https://github.com/jcmoraisjr/haproxy-ingress/pull/1374) (jcmoraisjr)
  * Configuration keys:
    * [`request-add-headers`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#headers)
    * [`request-del-headers`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#headers)
    * [`request-set-headers`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#headers)
    * [`response-add-headers`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#headers)
    * [`response-del-headers`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#headers)
    * [`response-set-headers`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#headers)
* add redirect scheme filter on gateway api [#1375](https://github.com/jcmoraisjr/haproxy-ingress/pull/1375) (jcmoraisjr)
* add cors filter on gateway api [#1377](https://github.com/jcmoraisjr/haproxy-ingress/pull/1377) (jcmoraisjr)
* add header modifier filters on gateway api [#1378](https://github.com/jcmoraisjr/haproxy-ingress/pull/1378) (jcmoraisjr)
* add reference grant support on gateway api [#1379](https://github.com/jcmoraisjr/haproxy-ingress/pull/1379) (jcmoraisjr)
* Bump sigs.k8s.io/controller-runtime from 0.22.4 to 0.23.1 [#1362](https://github.com/jcmoraisjr/haproxy-ingress/pull/1362) (dependabot[bot])
* bump haproxy and k8s versions on integration tests [#1383](https://github.com/jcmoraisjr/haproxy-ingress/pull/1383) (jcmoraisjr)
* improve dynamic update by adding/deleting servers [#1363](https://github.com/jcmoraisjr/haproxy-ingress/pull/1363) (jcmoraisjr)
  * Configuration keys:
    * [`dynamic-scaling`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#dynamic-scaling)
* update gateway api from v1.4.1 to v1.5.0 [#1392](https://github.com/jcmoraisjr/haproxy-ingress/pull/1392) (jcmoraisjr)
* Bump k8s.io/client-go from 0.35.1 to 0.35.2 [#1399](https://github.com/jcmoraisjr/haproxy-ingress/pull/1399) (dependabot[bot])
* remove CORS headers if Origin is not allowed [#1393](https://github.com/jcmoraisjr/haproxy-ingress/pull/1393) (jcmoraisjr)
* Bump k8s.io/klog/v2 from 2.130.1 to 2.140.0 [#1412](https://github.com/jcmoraisjr/haproxy-ingress/pull/1412) (dependabot[bot])
* Bump sigs.k8s.io/controller-runtime from 0.23.1 to 0.23.3 [#1413](https://github.com/jcmoraisjr/haproxy-ingress/pull/1413) (dependabot[bot])
* Bump golang.org/x/sync from 0.19.0 to 0.20.0 [#1414](https://github.com/jcmoraisjr/haproxy-ingress/pull/1414) (dependabot[bot])
* Bump docker/login-action from 3 to 4 [#1415](https://github.com/jcmoraisjr/haproxy-ingress/pull/1415) (dependabot[bot])
* Bump docker/build-push-action from 5 to 7 [#1416](https://github.com/jcmoraisjr/haproxy-ingress/pull/1416) (dependabot[bot])
* Bump docker/setup-qemu-action from 3 to 4 [#1417](https://github.com/jcmoraisjr/haproxy-ingress/pull/1417) (dependabot[bot])
* Bump docker/setup-buildx-action from 2 to 4 [#1418](https://github.com/jcmoraisjr/haproxy-ingress/pull/1418) (dependabot[bot])
* drop misplaced and noop http-server-close option [#1403](https://github.com/jcmoraisjr/haproxy-ingress/pull/1403) (jcmoraisjr)
* improve gateway api doc [#1409](https://github.com/jcmoraisjr/haproxy-ingress/pull/1409) (jcmoraisjr)
* configure CA for requests missing SNI [#1408](https://github.com/jcmoraisjr/haproxy-ingress/pull/1408) (jcmoraisjr)
  * Configuration keys:
    * [`auth-tls-default-secret`](https://haproxy-ingress.github.io/v0.17/docs/configuration/keys/#auth-tls)
* make cors allow header rule more flexible [#1401](https://github.com/jcmoraisjr/haproxy-ingress/pull/1401) (jcmoraisjr)
* add cors headers on proxy generated responses [#1402](https://github.com/jcmoraisjr/haproxy-ingress/pull/1402) (jcmoraisjr)
* remove default certificate for tlsroute api [#1396](https://github.com/jcmoraisjr/haproxy-ingress/pull/1396) (jcmoraisjr)
* Bump embedded haproxy from 3.0.12 to 3.0.18 [#1423](https://github.com/jcmoraisjr/haproxy-ingress/pull/1423) (jcmoraisjr)
* Bump Gateway API from 1.5.0 to 1.5.1 [#1425](https://github.com/jcmoraisjr/haproxy-ingress/pull/1425) (jcmoraisjr)
* Bump dependencies [#1424](https://github.com/jcmoraisjr/haproxy-ingress/pull/1424) (jcmoraisjr)
* Bump go from 1.25.3 to 1.26.1 [#1422](https://github.com/jcmoraisjr/haproxy-ingress/pull/1422) (jcmoraisjr)
* Bump HAProxy doc link from 2.8 to 3.0 [#1426](https://github.com/jcmoraisjr/haproxy-ingress/pull/1426) (jcmoraisjr)
* drop nbproc configuration [#1427](https://github.com/jcmoraisjr/haproxy-ingress/pull/1427) (jcmoraisjr)
* drop htx configuration [#1428](https://github.com/jcmoraisjr/haproxy-ingress/pull/1428) (jcmoraisjr)

Chart improvements since `v0.16`:

* PodDisruptionBudget - Add feature to control PDB via maxUnavailable as well [#95](https://github.com/haproxy-ingress/charts/pull/95) (PerGon)
* Toggle controller service [#99](https://github.com/haproxy-ingress/charts/pull/99) (kozhukalov)
* add list and watch permission to namespace and node apis [#101](https://github.com/haproxy-ingress/charts/pull/101) (jcmoraisjr)
* Do not check k8s version when creating IngressClass [#102](https://github.com/haproxy-ingress/charts/pull/102) (kozhukalov)
* add read access to ReferenceGrant API [#103](https://github.com/haproxy-ingress/charts/pull/103) (jcmoraisjr)
* add daemonset extraHostPorts [#98](https://github.com/haproxy-ingress/charts/pull/98) (bapung)
* add gatewayclass resource [#104](https://github.com/haproxy-ingress/charts/pull/104) (jcmoraisjr)

## Fixes (a1)

* fix full-sync trigger [#1331](https://github.com/jcmoraisjr/haproxy-ingress/pull/1331) (jcmoraisjr)
* fix status update during shutdown [#1335](https://github.com/jcmoraisjr/haproxy-ingress/pull/1335) (jcmoraisjr)
* check if process is the same waiting for reload [#1355](https://github.com/jcmoraisjr/haproxy-ingress/pull/1355) (jcmoraisjr)
* fixes wildcard match for gateway api [#1372](https://github.com/jcmoraisjr/haproxy-ingress/pull/1372) (jcmoraisjr)
* consider http header filter when sorting paths [#1373](https://github.com/jcmoraisjr/haproxy-ingress/pull/1373) (jcmoraisjr)
* fix hostname and certificate handling on gateway api [#1380](https://github.com/jcmoraisjr/haproxy-ingress/pull/1380) (jcmoraisjr)
* fix: report all node IPs in status for dual-stack support [#1381](https://github.com/jcmoraisjr/haproxy-ingress/pull/1381) (mia-mouret)
* fix: move logic to capture request origin header earlier [#1385](https://github.com/jcmoraisjr/haproxy-ingress/pull/1385) (ianroberts)
* Reflect back request origin when credentials enabled [#1389](https://github.com/jcmoraisjr/haproxy-ingress/pull/1389) (ianroberts)
* fix response error if proto is not http or https [#1395](https://github.com/jcmoraisjr/haproxy-ingress/pull/1395) (jcmoraisjr)
* skip invalid route and proto ref in gateway api [#1394](https://github.com/jcmoraisjr/haproxy-ingress/pull/1394) (jcmoraisjr)
* fix cors default origin on ingress api [#1400](https://github.com/jcmoraisjr/haproxy-ingress/pull/1400) (jcmoraisjr)
* fix missing apis on gateway api standard channel [#1410](https://github.com/jcmoraisjr/haproxy-ingress/pull/1410) (jcmoraisjr)
* fix server state file generation [#1404](https://github.com/jcmoraisjr/haproxy-ingress/pull/1404) (jcmoraisjr)
* improve validation of PEM data [#1405](https://github.com/jcmoraisjr/haproxy-ingress/pull/1405) (jcmoraisjr)
* fix rewrite path handling for regex paths [#1411](https://github.com/jcmoraisjr/haproxy-ingress/pull/1411) (jcmoraisjr)
