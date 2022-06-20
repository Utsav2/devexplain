# Devexplain: help understand developer metrics

The goal of this project is to help understand performance and set up guardrails on developer tools in real codebases. 
For example, OOMs in Typescript server can cause slow and flaky behavior in VSCode. 

Devexplain tries to be a hub to help collect all sorts of interesting metrics and send it to various sinks. 

## Config
Devexplain can be configured at startup with startup flags.


### Cadence
```
devexplain -cadence_minutes 15
```
Configures the interval at which metrics are scraped and sent to the various sinks.

### Continuous metrics
Devexplain is configured to read non-startup configuration from Starlark files, and it reloads config for every run loop, which lets developers make independent updates to the configs and the binary.


## Sources and Sinks
Devexplain can be configured in Starlark.

### Sources


Define a function called `sources()` and that lets you specify a dictionary of sources to metadata about that source. In this example, we start measuring process performance of `gopls`, the Go language server.

```
def sources():
    return {
        "gopls": {
            "process": {
                "name": "gopls",
            }
        },
    }
```

Process names can also be defined as a function/lambda for more interesting cases. The input to the function is a single string with the process name.

The output should be a boolean indicating whether the given process name matches what the source desires to measure.

For example:

```
def sources():
    return {
        "tsserver": {
            "process": {
                "name": lambda t: "tsserver" in t
                and "--serverMode partialSemantic" not in t, 
                # vscode starts up two tsservers by default - one for quick partial results, and another complete one.
                # for process monitoring, we're interested in the complete tsserver, so specify a filter to
                # skip the partial one.
            }
        },
    }
```


###  Sinks
Sinks are configured by another function. Currently, we support `stderr` and `datadog`.

```
def sinks():
    return {
        "stderr": {},
        "datadog": {},
    }
```


## Vision
Devexplain has all the scaffolding to begin measuring workflows and send that data on a specified interval. The goal would be to understand developer tool metrics and traces directly, like Typescript trace files, edit refresh times, and send them over for better observability into developer workflow performance.

Devexplain does not aim to be a developer performance monitoring tool. It's intended for developer experience/platform teams interested in tool reliability, not individual developer performance.
