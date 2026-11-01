# Release maintenance

This page preserves the repository's existing publishing procedure. It does not
define a language/documentation versioning policy.

## Publishing a release

Push your changes to `master` and wait for both **CI** and **Native vendor
packages** to succeed for the exact commit you want to ship. Then tag that
tested commit and push the tag:

```sh
git tag v0.2.1 TESTED_COMMIT_SHA
git push origin v0.2.1
```

The tag starts the lightweight **Release** workflow, not another full CI/native
test run. It verifies successful `master` push runs of both workflows for that
exact commit and downloads the tested SDL artifacts. It then rebuilds the
compiler with the release version embedded, packages all four platforms, and
smoke-tests the bundles outside the checkout. On macOS/Linux this runs the real
installer with local downloads, then builds through its symlink; Windows also
builds through the installed launcher. Checks compile/run hello-world and a
headless SDL program using std native sources, then relocate the application and
its deployed runtime files and hide the original installation before running it
with library-search overrides cleared. Only when all platforms pass are the
archives, installers, and `SHA256SUMS` published together. Tags containing `-`
(such as `v0.2.1-rc.1`) are published as prereleases. Published releases cannot
be overwritten.

Release stops if either workflow is unfinished or unsuccessful, or if its
matching native artifacts are missing, expired, or ambiguous. It never silently
rebuilds SDL or borrows artifacts from another commit, a PR, or a fork. Native
artifacts are retained for 30 days; rerun the original **Native vendor packages**
`master` push run for that commit if they expire, then rerun the failed release.
Fix a failing check and test the new commit before creating a new version tag.

Development bundles are also available for seven days from successful native
workflow runs under **Actions**.

Packaging discovers all packages under `base/std` and `base/vendor`, retains only
the selected target's `artifacts/` subtree, and validates every native manifest's
declared inputs. New vendor packages must support every advertised release target
or packaging fails. Non-system transitive runtime libraries must still be declared
in `runtime_files`; dependency discovery is not automatic. CI runner smoke checks
do not establish compatibility with every older OS or replace clean-machine tests.
