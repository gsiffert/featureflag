# featureflag

This package provides a simple way to create feature flags in Go.
Feature flags are a way to enable or disable features in your application without changing the code.
This can be useful for testing new features, or for enabling features for specific users.

This package has been designed to be modular, allowing you to only import the dependencies
that you need.

There are multiple packages completing this one, those packages allows to retrieve the Value of the feature flag from different sources:
- [ffaws](https://github.com/gsiffert/ffaws): Implements sources for AWS Secrets Manager, AppConfig, SSM Parameter Store.

If unfortunately, the source you need is not implemented, you can implement your own by implementing the `SourceReader` interface.
Feel free to open an issue if you wish to be mentioned in the list.

## Usage

```go
package main

import (
    "context"
    "time"
	
    "github.com/gsiffert/featureflag"
)

type Config struct {
    Feature1 bool `json:"feature1"`
    Feature2 bool `json:"feature2"`
}

func main() {
    fileSource := featureflag.NewFileReader("/etc/config.json")
    source := featureflag.MapSourceReader(fileSource, featureflag.MapJSON[Config])
 
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    
    // This featureFlag will be updated every 30 seconds with the latest secret from the file.
    featureFlag, err := featureflag.New(ctx, 30*time.Second, source)
    if err != nil {
        panic(err)
    }
}
```
