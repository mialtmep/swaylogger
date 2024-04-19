# swaylogger
## This project is for educational purposes ONLY. I am not responsible for any direct or indirect damage and/or misuse. This technique is not unique and was publicly available at the time of implementation and publication.
Indented use of this particular implementation — a practical example to show that the QWERTY layout is ineffective and obsolete.
### Usage
1. get sway pid with `ps -axu | grep -v 'grep' | grep -iP "sway\$"`
2. get input id with `swaymsg -t get_inputs -r | jq "."`
3. in `/dev/input/by-id` find correct input id
4. in `/proc/%sway pid%/fd` find correct file descriptor symlinked to the input file
5. `./swaylogger %path to the symlinked file of sway process% %sway input id%`

Records symbols from English and Russian layouts and functional key, ignoring everything else.
