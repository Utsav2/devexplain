def sources():
    return {
        "gopls": {
            "process": {
                "name": "gopls",
            }
        },
    }

# where to log data.
def sinks():
    return {
        "stderr": {},
        "datadog": {},
    }