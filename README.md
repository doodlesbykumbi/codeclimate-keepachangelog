# Code Climate Keep A Changelog Engine

`codeclimate-keepachangelog` is a Code Climate engine that wraps [parse_a_changelog](https://github.com/cyberark/parse-a-changelog/). You can run it on your command line using the Code Climate CLI, or on our hosted analysis platform.

[parse_a_changelog](https://github.com/cyberark/parse-a-changelog) is a validator for changelogs using the Keep a Changelog standard (http://keepachangelog.com).

### Installation

1. If you haven't already, [install the Code Climate CLI](https://github.com/codeclimate/codeclimate).
2. Add the following to your Code Climate config:
  ```yaml
  plugins:
    keepachangelog:
      enabled: true
  ```
3. Run `codeclimate engines:install`
4. You're ready to analyze! Browse into your project's folder and run `codeclimate analyze`.

### Building

```console
make image
```

This will build a `codeclimate/codeclimate-keepachangelog` image locally.

### Updating

[parse_a_changelog](https://github.com/cyberark/parse-a-changelog) is a validator for changelogs using the Keep a Changelog standard (http://keepachangelog.com). Once in a
while a new version is released. In order to get the latest version
and force a new docker image build, please update the base image in the
`Dockerfile`. Please avoid any unstable tags such as `latest` and keep it
explicit.

### Need help?

For help with parse_a_changelog, [check out their documentation](https://github.com/cyberark/parse-a-changelog).

If you're running into a Code Climate issue, first look over this project's [GitHub Issues](https://github.com/doodlesbykumbi/codeclimate-keepachangelog/issues), as your question may have already been covered. If not, [go ahead and open a support ticket with us](https://codeclimate.com/help).
