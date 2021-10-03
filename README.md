# Medieival-Dynasty-save-editor
Basic CLI save editor for Medieival Dynasty written in Go.
![](https://i.imgur.com/45qFQYH.png)
[Windows binaries](https://github.com/Sorrow446/Medieival-Dynasty-save-editor/releases)

# Usage
**Backup your saves first. Made for GOG 1.0.0.6.**

Give 1000 coins and 500 dynasty reputation:   
`md_save_editor.exe G:\md_save_backups\Quicksave.sav -c 1000 -r 500`

You can use the batch files instead by dragging your save file onto them.

```
Usage: md_save_editor [--age AGE] [--coins COINS] [--reputation REPUTATION] PATH

Positional arguments:
  PATH

Options:
  --age AGE, -a AGE [default: -1]
  --coins COINS, -c COINS [default: -1]
  --reputation REPUTATION, -r REPUTATION [default: -1]
  --help, -h             display this help and exit
  ```

# Disclaimer       
- This has no partnership, sponsorship or endorsement with Render Cube or Toplitz Productions.
- I won't be responsible if you break your saves.
