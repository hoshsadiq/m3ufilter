# m3ufilter

This is a utility that will allow you to cleanup your M3U/M3U8/M3U+ files. This can change feed titles, names, tvg attributes, add/remove additional entries and much more.

## How to run
Simply create a config (see example below), and the you can run the binary:
```yaml
m3ufilter -config /path/to/config.yaml
```

The command has the following arguments
```yaml
Usage of m3ufilter:
  -config string
        Config file location (default "~/.m3u.conf")
  -log-output string
        Where to output logs. Defaults to stderr
  -playlist-output string
        Where to output the playlist data. Ignored when using -server flag. Defaults to stdout
  -server
        Run a server to retrieve the playlist as a URL
```

### Example config
```yaml
core:
  listen: localhost:8080
  sync_title_name: true
providers:
  - uri: file:///path/to/m3u/playlist.m3u
    filters:
      - match(Attr["group-title"], "UK.*") && !match(Title, "^24/7")
      - match(Attr["tvg-id"], "3e.ie")
    replacements:
      name:
        - find: "[\\s\\:\\|]+"
          replace: " "
        - find: "^VIP "
          replace: ""
      attributes:
        tvg-name:
          - find: "[\\s\\:\\|]+"
            replace: " "
          - find: "^VIP "
            replace: ""
    setters:
      - name: replace(Title, "NEWS", "News")
        attributes:
          tvg-id: tvg_id(Title) + ".us"
        filters:
          - Title == "ABC News"
          - Title == "USA ABC NEWS HD"
          - Title == "USA CNN"
          - Title == "CNN"
          - Title == "CNN HD"
```

The following functions are available:

- `strlen(text string) int`

    Will return the length of the string
- `match(subject string, regexp string) bool`

    Will return `true` if the `subject` matches the regular expression
- `replace(subject string, find_regexp string, replace string) string`

    Will look for the regular express `find_regexp` and replace with the value of `replace` and return that.
- `tvg_id(text string)`

    Will try its best to turn text into a valid tvg-id attribute value. This does not include the usual country extension. The idea is that you pass the channel name into this, and it will spit out something that can be used as tvg-id.

    For example:
    ```
    tvg_id("CCN HD") > cnn
    ```

    The login behind this will be improved, but right now, all it does is simply remove SD/HD/FHD from the title and any character that isn't a-zA-Z0-9.

## Future plans
The idea behind this is to a be one stop shop for generating both xmltv and m3u files from any source.
This will eventually add support for xml, and will automatically try and match up channels and EPG data should this not exist.
Any other ideas you have? Feel free to raise a ticket.
