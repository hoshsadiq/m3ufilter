# m3ufilter

This is a utility that will allow you to cleanup your M3U/M3U8/M3U+ files based on a CSV.

> Warning: right now, due to its rapid development, I cannot guarantee breaking changes.

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
#### Main config
```yaml
core:
  server_listen: localhost:8080
  update_schedule: "*/24 * * * *"
  group_order:
    - Entertainment
    - Family/Kids
    - News
    - Drama
providers:
  - uri: file:///path/to/m3u/playlist.m3u
    csv: file:///path/to/m3u/playlist.csv
epg_providers:
  - uri: file:///path/to/epg.xml
    channel_id_renames:
      - from: find
        to: replacement
      - from: BBC One
        to: bbc.uk
```

#### CSV config
```csv
search-name,chno,id,name,group,shift,logo
BBC One UK,1,bbc1.uk,BBC One,Entertainment,,
CNN UK,160,cnn.uk,CNN,News,,
``` 

#### The meaning of the config options are as follows:

##### Core config
- `core.server_listen` (`string`)

    If set, this will run as a server, rather a single run. If you want single runs, you can omit this option.
    Default: disabled

- `core.update_schedule` (`string`)

    How often it should retrieve the latest channels. This is expressed in [cron syntax](https://github.com/mileusna/crontab#crontab-syntax-). It is highly recommended that you do not set this to a low interval. Set this to at least once every day. You can use the tool [crontab.guru](https://crontab.guru/) to figure out what interval you want.
    Default: `0 */24 * * *` (that is, once every 24 hours)

- `core.auto_reload_config` (`true|false`)

    Whether or not to reload the config before every run. Please note, this will not affect `core.server_listen` and `core.update_schedule`
    Default: `true`

- `core.group_order` (`list` of `string`) (experimental)

    The order to put the categories in.

##### IPTV providers

- `providers`

    This is a list of providers of where to retrieve M3U lists. This is an array (see example above).

- `providers.uri` (`string`)

    The URL of where to retrieve the M3U list. This can start with `file://` to retrieve a list from a local file system (this must be an absolute path).

- `providers.csv` (`string`)

    The URL of where to retrieve the CSV data list. This can start with `file://` to retrieve a list from a local file system (this must be an absolute path). See below for more information on the format.
    
    If this is set, only entries in this file will be in the final CSV.

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

##### EPG providers

- `epg_providers`

    This is a list of providers of where to retrieve EPG data from. Each entry in here must be able to retrieve _valid_ XMLTV data. If it unable to decode the full XML, it will skip over it. This is an array (see example above).

- `epg_providers.uri` (`string`)

    The URL of where to retrieve the EPG data. This can start with `file://` to retrieve a list from a local file system (this must be an absolute path).

- `epg_providers.channel_id_renames` (`map`)

    This is a key value pair of channel IDs to rename within the XMLTV. This is useful in case the EPG data's channel IDs don't match the channel IDs in the M3U files. This will change the ID.

#### CSV configuration
In order to manage the m3u, you will need to create a CSV with the relevant values you want to set everything to. The CSV must have the following columns:

- `search-name`
- `chno`
- `tvg-id`
- `tvg-name`
- `group-title`
- `tvg-shift`
- `tvg-logo`

You can set these up as the first row, in which case, the order of columns does not matter (and you can also remove individual columns to ignore them), however, if you don't have that header row, it must be in the above order.

What will happen is each M3U stream name (the bit after the `,` on the `#EXTINF` line) will be matched up against the `search-name` column. If no match is found, the stream will not be included in the final playlist. If it is, the stream is included in the final playlist, in addition to the values in the other columns getting set.

For example, if you have a CSV as follows:

```csv
search-name,chno,id,name,group,shift,logo
BBC One UK,1,BBC1.uk,BBC One,Entertainment,,
```

With the following M3U:
```m3u
#EXTM3U
#EXTINF:-1 tvg-id="" tvg-name="BBC" tvg-logo="https://picon-13398.kxcdn.com/bbcone.jpg" group-title="UK",BBC One UK
http://somewhere.com/111
#EXTINF:-1 tvg-id="Channel4.uk" tvg-name="Channel 4 HD UK" tvg-logo="https://picon-13398.kxcdn.com/channel4hd.jpg" group-title="UK",Channel 4 HD UK
http://somewhere.com/222
```

The final playlist will be:
```m3u
#EXTM3U
#EXTINF:-1 tvg-id="BBC1" tvg-name="BBC One" tvg-logo="https://picon-13398.kxcdn.com/bbcone.jpg" group-title="Entertainment",BBC One
http://somewhere.com/111
```

Notice how the `search-name` column matches the `BBC One UK` in the source, and in the final, each of the values match the relevant columns in the CSV.
Notice, also, that Channel 4 is not in the final list, as it wasn't in the CSV.

You can generate an initial CSV file using the `-csv` argument. Please note, that this will output additional columns that are otherwise ignored. You can remove them, or leave them.

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
