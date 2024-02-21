# Tesla Vehicle Command SDK

This is greenmissions forked version of the Tesla Vehicle Command SDK. The original repository can be found [here](https://github.com/teslamotors/vehicle-command), where documentation and examples can be found.

To publish a new version of this package, follow these steps:

1. Create a new release in the fork, following semantic versioning.
2. Then you need to update goproxy to include the new version.
```
GOPROXY=proxy.golang.org go list -m github.com/greenmission/vehicle-command@<version>
```

The new version should now be available for use.