# Terminal showcase for the README.
# Render from the repository root with:
#   vhs docs/demo.tape

Output docs/demo.gif

Set Shell "bash"
Set FontFamily "GoogleSansCode Nerd Font Mono"
Set FontSize 21
Set LineHeight 1.3
Set LetterSpacing 0.4
Set Width 1180
Set Height 520
Set Padding 16
Set Theme "Catppuccin Mocha"
Set WindowBar "Colorful"
Set WindowBarSize 40
Set BorderRadius 12
Set TypingSpeed 45ms
Set Framerate 30
Set PlaybackSpeed 1.0

Hide
Type "export PATH=$PWD:$PATH"
Enter
Type "export PS1='❯ '"
Enter
Type "clear"
Enter
Show
Sleep 500ms

Type "gator users"
Enter
Sleep 2s

Type "gator feeds"
Enter
Sleep 2.5s

Type "gator following"
Enter
Sleep 2s

Hide
Type "clear"
Enter
Show
Sleep 600ms

Type "gator browse 2"
Enter
Sleep 4.5s
