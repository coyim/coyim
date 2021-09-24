Instructions for creating a release
-----------------------------------

- Check and make sure that all the core translations are up to date
- Update the authors file:
```sh
    make authors
```
- Run all tests and check that the current version runs well, preferably on multiple platforms
- Write release notes, following the established standard
- Create a signed tag
```sh
    git tag -s v0.x.y -m "Release version 0.x.y"
```
- Push the signed tag
```sh
    git push origin v0.x.y
```
- Wait for Github Actions to finish building and publishing all binaries to the Github Release page
- Wait for Travis to finish building and publishing the macOS build to the Github Release page
- Add release notes to the Github Release page
- Create a new release note entry on the CoyIM website
- Create a new blog post about the release on the CoyIM website
- Add the new release version for download to the website
- Build and push reproducibility signatures from as many as possible (see [REPRODUCIBILITY](REPRODUCIBILITY.md) for instructions)
- Tweet from @coyproject about the new release

