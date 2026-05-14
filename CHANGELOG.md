# Changelog

## [1.5.0](https://github.com/bluefunda/trm-cli/compare/v1.4.2...v1.5.0) (2026-05-14)


### Features

* add macOS notarization to brew cask installer ([243f4a4](https://github.com/bluefunda/trm-cli/commit/243f4a4b172ca2e11cf572fd249bbfd9afd63a46))

## [1.4.2](https://github.com/bluefunda/trm-cli/compare/v1.4.1...v1.4.2) (2026-05-13)


### Bug Fixes

* remove Docker Hub release, publish to ghcr only ([cdb3ce3](https://github.com/bluefunda/trm-cli/commit/cdb3ce3ecaef0e5f3e6cfebe5f005aff54954803))

## [1.4.1](https://github.com/bluefunda/trm-cli/compare/v1.4.0...v1.4.1) (2026-05-13)


### Bug Fixes

* use shasum on macOS for checksum verification ([5e99540](https://github.com/bluefunda/trm-cli/commit/5e995405c07a02331c895d3ba15e5345fe3813f7))

## [1.4.0](https://github.com/bluefunda/trm-cli/compare/v1.3.0...v1.4.0) (2026-05-13)


### Features

* **cli:** add shell completion support for all major shells ([#16](https://github.com/bluefunda/trm-cli/issues/16)) ([d6e9b27](https://github.com/bluefunda/trm-cli/commit/d6e9b27ea2c29157e43fb9e12f2740ce30b862f9))

## [1.3.0](https://github.com/bluefunda/trm-cli/compare/v1.2.1...v1.3.0) (2026-05-13)


### Features

* add change request and comment CLI commands ([#13](https://github.com/bluefunda/trm-cli/issues/13)) ([afaa78a](https://github.com/bluefunda/trm-cli/commit/afaa78a83df4d60ef8140183c0c8f9b59911dbb9))

## [1.2.1](https://github.com/bluefunda/trm-cli/compare/v1.2.0...v1.2.1) (2026-04-24)


### Bug Fixes

* use DOCKER_USERNAME/DOCKER_PASSWORD org secrets, add continue-on-error to description step ([962beb1](https://github.com/bluefunda/trm-cli/commit/962beb1e65be6ed4d71d11383943b42dd5e193ef))

## [1.2.0](https://github.com/bluefunda/trm-cli/compare/v1.1.1...v1.2.0) (2026-04-24)


### Features

* add Docker image publishing to Docker Hub and ghcr.io ([de65a64](https://github.com/bluefunda/trm-cli/commit/de65a641e17f0ce9b5ad054dc60eab0e5e48fce9))

## [1.1.1](https://github.com/bluefunda/trm-cli/compare/v1.1.0...v1.1.1) (2026-04-23)


### Bug Fixes

* update homebrew-patch to use requests.rb and GITHUB_TOKEN for public repo ([263532b](https://github.com/bluefunda/trm-cli/commit/263532b1f9bfe280edeec96b0fd03e4640299395))

## [1.1.0](https://github.com/bluefunda/trm-cli/compare/v1.0.2...v1.1.0) (2026-04-23)


### Features

* rename binary from trm to requests, add cmd entry point ([e3ca5c5](https://github.com/bluefunda/trm-cli/commit/e3ca5c5a88e2e453d7d68174e49f26f5b2ab4111))

## [1.0.2](https://github.com/bluefunda/trm-cli/compare/v1.0.1...v1.0.2) (2026-04-23)


### Bug Fixes

* discard conn.Close() error in user info command ([29e8d37](https://github.com/bluefunda/trm-cli/commit/29e8d37441c31afe643f0d700ae0f866ae6586d3))

## [1.0.1](https://github.com/bluefunda/trm-cli/compare/v1.0.0...v1.0.1) (2026-04-23)


### Bug Fixes

* add required version field to golangci.yml for v2 ([#5](https://github.com/bluefunda/trm-cli/issues/5)) ([ef56674](https://github.com/bluefunda/trm-cli/commit/ef56674ceb74abb7a804da4674925baab386fc6f))
* discard w.Flush() error in Table printer ([#7](https://github.com/bluefunda/trm-cli/issues/7)) ([39a1b50](https://github.com/bluefunda/trm-cli/commit/39a1b509ab75c0af06af804b833641c4642a5eca))

## 1.0.0 (2026-04-23)


### Features

* initial trm-cli scaffold for bluerequests platform ([#2](https://github.com/bluefunda/trm-cli/issues/2)) ([37db3ed](https://github.com/bluefunda/trm-cli/commit/37db3ede1d8b2c79b2d93d16169ce44395cdee6e))


### Bug Fixes

* resolve errcheck lint violations ([#4](https://github.com/bluefunda/trm-cli/issues/4)) ([8ec7318](https://github.com/bluefunda/trm-cli/commit/8ec73188c10f54c9b3b46c391d7a4d27bf864263))
