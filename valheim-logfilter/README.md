# valheim-logfilter
Filter Valheim dedicated server log output

Valheim server by default logs a lot of noise. These env variables allow users to remove unwanted lines from the log.

| Prefix | Default | Purpose |
|----------|----------|-------|
| `VALHEIM_LOG_FILTER_EMPTY` | `true` | Filter empty log lines |
| `VALHEIM_LOG_FILTER_MATCH` | ` ` | Filter log lines exactly matching |
| `VALHEIM_LOG_FILTER_STARTSWITH` | `(Filename:` | Filter log lines starting with |
| `VALHEIM_LOG_FILTER_ENDSWITH` |  | Filter log lines ending with |
| `VALHEIM_LOG_FILTER_CONTAINS` |  | Filter log lines containing |
| `VALHEIM_LOG_FILTER_REGEXP` |  | Filter log lines matching regexp |

All environment variables except for `VALHEIM_LOG_FILTER_EMPTY` are prefixes. Meaning you can define multiple matches like so:
```
-e VALHEIM_LOG_FILTER_STARTSWITH=foo \
-e VALHEIM_LOG_FILTER_STARTSWITH_BAR=bar \
-e VALHEIM_LOG_FILTER_STARTSWITH_SOMETHING_ELSE="some other filter"
```


# Usage
```
Usage of valheim-logfilter:
  -alsologtostderr
    	log to standard error as well as files
  -env-contains string
    	Valheim contains filter varname prefix (default "VALHEIM_LOG_FILTER_CONTAINS")
  -env-empty string
    	Valheim empty-line filter varname (default "VALHEIM_LOG_FILTER_EMPTY")
  -env-endswith string
    	Valheim ends-with filter env varname prefix (default "VALHEIM_LOG_FILTER_ENDSWITH")
  -env-match string
    	Valheim match filter env varname prefix (default "VALHEIM_LOG_FILTER_MATCH")
  -env-regexp string
    	Valheim regexp filter varname prefix (default "VALHEIM_LOG_FILTER_REGEXP")
  -env-startswith string
    	Valheim starts-with filter env varname prefix (default "VALHEIM_LOG_FILTER_STARTSWITH")
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
```