# Reproducibility

CoyIM currently only supports reproducible builds on Linux with AMD64. This document describes both how to do this, but also how to verify the existing signatures. The CoyIM reproducibility process generates a file called build_file that contains the SHA256 sum of both the GUI and the CLI binaries of CoyIM. Anyone that generates the same file can then generate a detached armored signature and make that available for others to verify.

## Generating reproducible binaries

In order to generate reproducible binaries, you need to have docker installed. For some operating systems with SELinux you also need to mark the coyim checkout directory as being available from inside of Docker, using this command, where DIR is the coyim directory:

```sh
  chcon -Rt svirt_sandbox_file_t $DIR
```

In order to build CoyIM reproducibilly, you simply do

```sh
  make reproducible-linux-build
```

inside of the CoyIM directory. This will create a new Docker image and then use it to build CoyIM. At the end of the process, it will generate three files:

```sh
  bin/coyim
  bin/coyim-cli
  bin/build_info
```

If you want to sign the build\_info file using your default GPG key, you can simply run

```sh
  make sign-reproducible
```

This will generate

```sh
  bin/build_info.0xAAAAAAAAAAAAAAAA.asc
```

where `0xAAAAAAAAAAAAAAAA` is the long-form key ID of your GPG key.

If you have access to the `twstrike/coyim` Bintray account, you can put your bintray username in `~/.bintray_api_user` and your API key in `~/.bintray_api_key` and then simply run:

```sh
  make upload-reproducible-signature
```

If not, you can run

```sh
  make send-reproducible-signature
```

which will mail the signed `build_info` file to [coyim@thoughtworks.com](mailto:coyim@thoughtworks.com). You can also manually mail the this file of course.


## Verifying reproducible binaries

From v0.4 each release of CoyIM will have several signatures for build\_info files available. You can of course download and verify each one of those signatures manually, but we also provide a simple way of verifying it using a small Ruby script. It can be invoked like this:

```sh
  make check-reproducible-signatures
```

This will download everything necessary for the current tag (so you should first check out the tag you want to verify), and then verify that the coyim and coyim-cli binaries match the hashes inside of the `build_info` file, and then verify that each signature checks out for the `build_info` file.
