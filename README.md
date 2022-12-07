# wl-screenshot

wayland screenshot utility

## Dependencies

- wl-clipboard
- grim
- slurp
- swappy

## Install

```sh
make install
```

## Usage

```
Usage of wl-screenshot:
  -capture string
        Capture area or all (default "area")
  -edit
        Edit image with swappy
  -notify
        Show notification with notify-send
  -path string
        Path to save file

```

## Sway elegant binding

```
# default screenshot area to clipboard
bindsym --to-code Print exec wl-screenshot

# default screenshot and open swappy
bindsym --to-code Print exec wl-screenshot -edit

# scrap area and save to file
bindsym --to-code Print exec wl-screenshot -capture area -path ~/Pictures/
```
