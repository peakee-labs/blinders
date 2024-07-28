### Build State

Mark each build output with a commit hash

### A service is in `change state`:

Case 1:

- Having un-staged or staged git changes in its director
- Having git changes in local submodules

Could be redeploy but will not having commit hash (un committed deployment)

Case 2:

Have no change but `the current commit hash` != `the last deployment commit hash`

### Force deploy -> deploy all
