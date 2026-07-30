# Known issues

Current gaps in the implementation.

## Custom command mode (option 2) is disabled

Currently, custom command invocation via cli after `./benchwire` is disabled.
Support will be introduced in the future.

## No resume support

If a batch is interrupted partway through, the partial yaml files stay in
`results/yaml/`. The next invocation deletes or overwrites them and starts over
from zero, there's no resume.

