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
  group_order:
    - Entertainment
    - Family/Kids
    - News
    - Drama
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
epg_providers:
  - url: file:///path/to/epg.xml
    channel_id_renames:
      replacement: find # key = what to replace it with, value = what to find
      bbc.uk: "BBC One"
```

#### The meaning of the config options are as follows:

##### Core config
- `core.server_listen` (`string`)

    If set, this will run as a server, rather a single run. If you want single runs, you can omit this option. See the arguments to specify the output.
    Default: disabled

- `core.update_schedule` (`string`)

    How often it should retrieve the latest channels. This is expressed in [cron syntax](https://github.com/mileusna/crontab#crontab-syntax-). It is highly recommended that you do not set this to a low interval. Set this to at least once every day. You can use the tool [crontab.guru](https://crontab.guru/) to figure out what interval you want.
    Default: `* */24 * * *` (that is, once every 24 hours)

- `core.auto_reload_config` (`true|false`)

    Whether or not to reload the config before every run. Please note, this will not affect `core.server_listen` and `core.update_schedule`
    Default: `true`

- `core.output` (`m3u|csv`)

    What to output. This can be either `csv` or `m3u`. CSV is useful for debugging and ensuring you've gone through all the channels, outside that, you generally want this to be `m3u`.
    Default: `m3u`

- `core.group_order` (`list` of `string`) (experimental)

    The order to put the categories in.

##### IPTV providers

- `providers`

    This is a list of providers of where to retrieve M3U lists. This is an array (see example above).

- `providers.url` (`string`)

    The URL of where to retrieve the M3U list. This can start with `file://` to retrieve a list from a local file system (this must be an absolute path).

- `providers.ignore_parse_errors` (`true|false`)

    If true, this will ignore errors when trying to parse an individual channel. E.g. `tvg-id="Channel "1"` would not be possible to parse, and will be ignored without any errors.
    Default: `false`

- `providers.check_streams.enabled` (`true|false`)

    If true, the stream URLs will be checked to see if they are alive before including them.
    Default: `false`

- `providers.check_streams.method` (`head|get`)

    How to validate whether a stream is available or not. Either using a `HEAD` request or `GET` request.
    
    The difference here is that `HEAD` is less likely to be correct. Often streams return that it is available when in reality is not. On the other hand, the GET method tries to actually retrieve zero bytes from the stream, thus this will give more accurate results. The problem with this is that, if your provider doesn't allow many connections at the same time, your stream will likely get cut off if the check is happening while you are watching a stream, and you'll have to wait until the update is finished before you can watch it without problems. If you time your updates to be when you're not watching any streams (e.g. when you're asleep for example), this shouldn't be a problem.
    
    Another thing to consider here is that, currently the stream available is only updated when the full update runs. Should the stream because available before the next run, as far as m3ufilter is concerned, this stream will stay unavailable until the next.
    
    Default: `head`

- `providers.check_streams.action` (`remove|noop`)

    The action to take if a stream is unavailable. Remove or don't take any action. `noop` is useful if you want to change the name of the channel to something else, e.g. prefix it with `[Unavailable]`.
    
    Default: `remove`

- `providers.filters` (`list` of `string`)

    A list of filters to apply to channels. This must return true or false (no strings or anything else). If it returns true, it will include the channel in the final list.
    You can use the functions and variables below to specify your logic.
    Default: `true` (meaning all channels are included)

- `providers.setters`

    A list of things to set on channels based on the filter.

- `providers.setters.name` (`string`)

    Set the name for this individual channel. This MUST return a string.

- `providers.setters.attributes`

    What to set any attribute too. This is go for setting logos where none exist, and/or enforcing `tvg-id` in case a channel does not have one but should. All attributes are listed below. This again, must return a string, and has the functions and variables below available.

    Example
    ```yaml
    providers:
      setters:
        tvg-id: mychannel.us
    ```

- `providers.setters.filters`

    A list of filters to limit this providers setter to. The same logic applies as the above filters method, and thus again, must return true/false. If true, it will run the setters.

##### EPG providers

- `epg_providers`

    This is a list of providers of where to retrieve EPG data from. Each entry in here must be able to retrieve _valid_ XMLTV data. If it unable to decode the full XML, it will skip over it. This is an array (see example above).

- `epg_providers.url` (`string`)

    The URL of where to retrieve the EPG data. This can start with `file://` to retrieve a list from a local file system (this must be an absolute path).

- `epg_providers.channel_id_renames` (`map`)

    This is a key value pair of channel IDs to rename within the XMLTV. This is useful in case the EPG data's channel IDs don't match the channel IDs in the M3U files. This will change the ID.

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

    The login behind this will be improved, but right now, all it does is simply remove SD/HD/FHD from the title and any character that isn't `a-zA-Z0-9`.

- `title(subject string) string`

    Will turn the text in `subject` into a title, by capitalising all words, and also ensures all letters in SD/HD/FHD are capitalised.

- `upper_words(subject string, word string...) string`

    Will turn the text `word` in `subject` into uppercase. Argument `word` can be repeated as multiple times.

- `starts_with(subject string, prefix string...) bool`

    Will return true if the text in `subject` starts with the text in `prefix`.

- `endss_with(subject string, suffix string...) bool`

    Will return true if the text in `subject` ends with the text in `suffix`.

##### Additionally, the following variables are available:

|variable|content|m3u tag mapping|
|--------|-------|-------|
|`ChNo`|The channel number|`tvg-chno`|
|`Id`|The ID to sync up with XMLTV|`tvg-id`|
|`Name`|This is the channel name|`tvg-name`|
|`Uri`|The URL for the stream|The URL (not a tag)|
|`Duration`|The duration of the stream, this is usually -1 due to Live TV being being.. well.. live.|The duration (not a tag)|
|`Logo`|The logo (can be either a url or a base64 data string)|`tvg-logo`|
|`Language`|The language of the stream|`tvg-language`|
|`Group`|The group category|`group-title`|

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

### Server endpoints
The following server endpoints are available for use:

- `GET /playlist.m3u`
  This will return the final filtered and updated playlist. This is what you would point your player to so it can get an up to date playlist.

- `POST /update`
  This is used to force the application to retrieve the latest version of all the providers. This is an asynchronous operation, and will return 204 on success.

## Future plans
The idea behind this is to be one stop shop for generating both xmltv and m3u files from any source.
This will eventually add support for xml, and will automatically try to match up channels and EPG data should this not exist.
Any other ideas you have? Feel free to raise a ticket.
