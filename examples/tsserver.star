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

# where to log data.
def sinks():
    return {
        "stderr": {},
    }