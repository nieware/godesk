# wmctrl notes

## list all windows

`wmctrl -lp`

The format of the window list: `<window ID> <desktop ID> <process ID> <client machine> <window title>`

	0x02e00012 -1 123                   N/A Desktop — Plasma
	0x02e00016 -1 123                   N/A Desktop — Plasma
	0x02e0001a -1 123                   N/A Desktop — Plasma
	0x02e00048 -1 123                   N/A Plasma
	0x03600011  0 123 rnieren-Latitude-E7470 Inbox - r.nieren@sportradar.com - Mozilla Thunderbird
	0x05e00006  0 123 rnieren-Latitude-E7470 Slack - Sportradar
	0x00a0003e  0 123 rnieren-Latitude-E7470 SRC [~/projects/ITF_Media_Platform/SRC] - .../staging/scoreboard/solutions/widgets/solution/components/plugin/timeline.js [SRC] - PhpStorm
	0x00a001f7  1 123 rnieren-Latitude-E7470 fmp-handball [~/projects/FMP/fmp-handball/fmp-handball] - .../app/Models/RuleSet.php [fmp-handball] - PhpStorm
	0x06c00001  0 123 rnieren-Latitude-E7470 New Tab - Google Chrome
	0x06c00006  0 123 rnieren-Latitude-E7470 [FMPCORE-590] [FMP DHB/HBL] Allow to overwrite points of a match with flexible points for each team - Sportradar JIRA - Google Chrome
	0x06c00007  0 123 rnieren-Latitude-E7470 Live Scores - ITF Tennis - Pro Circuit - Google Chrome
	0x01e00003  0 123 rnieren-Latitude-E7470 ~/projects/ITF_Media_Platform/201902_DCQLS_FCWG_issues/FC WG WG2 notes.txt (201902_DCQLS_FCWG_issues) - Sublime Text (UNREGISTERED)
	0x01e0002e  1 123 rnieren-Latitude-E7470 ~/projects/FMP/FMP Notes.md - Sublime Text (UNREGISTERED)
	0x01e00031  0 123 rnieren-Latitude-E7470 ~/projects/ITF_Media_Platform/20190228_PbP-investigation/PbP_investigation.md • - Sublime Text (UNREGISTERED)

## list all desktops

`wmctrl -d`

The format of the desktop list: `<desktop ID> [-*] <geometry> <viewport> <workarea> <title>`

	0  * DG: 5760x1200  VP: 0,0  WA: 0,0 5760x1200  ITF                                     
	1  - DG: 5760x1200  VP: 0,0  WA: 0,0 5760x1200  FMP         	

## move window with id to given desktop

`wmctrl -i -r 0x03600011 -t 1`

## move window to all desktops

one of these may or may not work (see https://unix.stackexchange.com/questions/11893/command-to-move-a-window-to-all-desktops):

`wmctrl -i -r 0x03600011 -t -2` (not -1!)
`wmctrl -i -r 0x03600011 -b add,sticky`



