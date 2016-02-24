# Gitlab Users Sync on Groups

## PREREQUISITE

Copy config.json.sample and edit it with your url and private key.

## USAGE

```
  -help
    	Show usage

  -id string
    	Specify repository id

  -ids string
    	A list of id separated by comma

  -m string
    	Specify method to retrieve groups infos, available methods:
    - -m users
    - -m users          
    - -m groups         -ids GROUP_IDS or -id GROUP_ID    -search PATTERN
    - -m team           -id GROUP_ID
    - -m new_member     -id GROUP_ID
    - -m sync_members   -ids GROUP_IDS or -id GROUP_ID    -search PATTERN

  -search string
    	Specify a pattern to search
```
