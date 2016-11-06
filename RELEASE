Instructions for creating a release
-----------------------------------

- Run all tests and check that the current version runs well, preferably on multiple platforms
- Write release notes, following the established standard
- Create a signed tag
    git tag -s v0.x.y -m "Release version 0.x.y"
- Push the signed tag
    git push origin v0.x.y
- Add release notes to the tag on Github
- Wait for Travis to finish building and publishing all binaries to Bintray
- Create a new blog post on coyim-pages with the release notes
- Update the config on coyim-pages to make the new version the current release
- Build and push reproducibility signatures from as many as possible (see REPRODUCIBILITY for instructions)
- Tweet from @coyproject about the new release
