# Dotty

🌎 Simple dotfile manager

## Description

- OK
- OK_SLASH: ok, but symlink has an appended slash
- BROKEN_S: symlink exists, but points to wrong location
- ROGUE_F: file exists in location
- ROGUE_D: dir exists in location
- TOLINK_F: file does not exist in location. will be linked during reconciliation
- TOLINK_D: same for dir
- MISSING: file missing from both src and dest

## TODO:

- ensure valid link at any point in time (creation, pre-existance) points to valid file or folder. if resolved-to-location non-existant, prompt and create
