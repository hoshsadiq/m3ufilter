# m3ufilter

This is a utility that will allow you to cleanup your M3U/M3U8/M3U+ files. This can change feed titles, names, tvg attributes, add/remove additional entries and much more.

> Warning: right now, due to it's rapid development, I cannot guarantee breaking changes.

## How to run
Simply create a config (see example below), and the you can run the binary:
```yaml
m3ufilter -config /path/to/config.yaml
```

The command has the following arguments
```
Usage of m3ufilter:
  -config string
        Config file location (default "~/.m3u.conf")
  -log string
        Where to output logs. Defaults to stderr
  -playlist string
        Where to output the playlist data. Ignored when in server mode. Defaults to stdout
```

### Example config
```yaml
core:
  server_listen: localhost:8080
  update_schedule: "*/24 * * * *"
  output: m3u
providers:
  - uri: file:///path/to/m3u/playlist.m3u
    filters:
      - match(Group, "UK.*") && !match(Name, "^24/7")
      - match(Id, "3e.ie")
    setters:
      - name: replace(Name, "[\\s\\:\\|]+", " ")
      - name: replace(Name, "^VIP ", "")
      - name: replace(Name, "USA", "")
        attributes:
          tvg-id: tvg_id(Name) + ".us"
        filters:
          - Name == "USA CNN"
          - Name == "CNN"
          - Name == "CNN HD"
```

##### The meaning of the config options are as follows:
- `core.server_listen`

    If set, this will run as a server, rather a single run. If you want single runs, you can omit this option. See the arguments to specify the output.
    Default: disabled

- `core.update_schedule`

    How often it should retrieve the latest channels. This is expressed in [cron syntax](https://github.com/mileusna/crontab#crontab-syntax-). It is highly recommended that you do not set this to a low interval. Set this to at least once every day.
    Default: `true`

- `core.output`

    What to output. This can be either `csv` or `m3u`. CSV is useful for debugging and ensuring you've gone through all the channels, outside that, you generally want this to be `m3u`.

- `providers`

    This is a list of providers of where to retrieve M3U lists. This is an array (see example above).

- `providers.url`

    The URL of where to retrieve the M3U list. This can start with `file://` to retrieve a list from a local file system (this must be an absolute path).

- `providers.filters`

    A list of filters to apply to channels. This must return true or false (no strings or anything else). If it returns true, it will include the channel in the final list.
    You can use the functions and variables below to specify your logic.
    Default: `true` (meaning everything is included)

- `providers.setters`

    A list of things to set on channels based on the filter.

- `providers.setters.name`

    Set the name for this individual channel. This MUST return a string.

- `providers.setters.attributes`

    What to set any attribute too. This is go for setting logos where none exist, and/or enforcing `tvg-id` in case a channel does not have one but should. All attributes are listen below under the `Attr` variable. This again, must return a string, and has the functions and variables below available.

    Example
    ```yaml
    providers:
      setters:
        tvg-id: mychannel.us
    ```

- `providers.setters.filters`

    A list of filters to limit this providers setter to. The same logic applies as the above filters method, and thus again, must return true/false. If true, it will run the setters.

##### For filters, name and setters, the following functions are available:

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

##### Additionally, the following variables are available:

|variable|content|
|--------|-------|
|`Id`|The ID to sync up with XMLTV|
|`Name`|This is the channel name|
|`Uri`|The URL for the stream|
|`Duration`|The duration of the stream, this is usually -1 due to Live TV being being.. well.. live.|
|`Logo`|The logo (can be either a url or a base64 data string)|
|`Group`|The group category|

##### Generic expression syntax

Within the filters, you may use [these syntax operations](https://github.com/maja42/goval#operators) to filter channels. Many of those are also available within the name and attribute section. As long as they return the expected data type.

#### gotchas
Due to the underlying library used for the logic parsing, setting a value to a generic string is not straight forward and must be double quoted, first with single quote, then double quote.

For example, if you want to set the for a channel to "My Channel", you have to do is as follows:
```yaml
setters:
  - name: '"My Channel"' # this works
    filters:
      - Name == "some criteria"
  - name: "My Channel" # this is invalid
    filters:
      - Name == "some criteria"
  - name: 'My Channel' # this is invalid
    filters:
      - Name == "some criteria"
  - name: My Channel # this is invalid
    filters:
      - Name == "some criteria"
```

In theory all of the above should be valid, but until a solution has been thought of, the workaround is to simply prefix it with an equals, e.g.:
```yaml
setters:
  - name: = My Channel
    filters:
      - Name == "some criteria"
```

Note that prefixing it with an equals marks the whole expression as literal string, excluding the equals. If you want a string with an equals in front of the text, you'll need to use two equals.

## Future plans
The idea behind this is to a be one stop shop for generating both xmltv and m3u files from any source.
This will eventually add support for xml, and will automatically try and match up channels and EPG data should this not exist.
Any other ideas you have? Feel free to raise a ticket.
