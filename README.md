# anomalia

`anomalia` is zero-dependency Go library for time series data analysis.

It supports anomaly detection and correlation. The API is simple and configurable in the sense that you can choose which algorithm suits your needs for anomaly detection and correlation.

:construction: **The library is currently under development so things might move or change!**

## Installation

Installation is done using `go get`:

```shell
go get -u github.com/amrfaissal/anomalia
```

## Supported Go Versions

`anomalia` supports `Go >= 1.10`.

## Documentation

### Quick Start Guide

This is a simple example to get you up & running with the library:

```go
package main

import (
    "fmt"
    "github.com/amrfaissal/anomalia"
)

func main() {
    // Load the time series from an external source.
    // It returns an instance of TimeSeries struct which holds the timestamps and their values.
    timeSeries := loadTimeSeriesFromAnExternalSource()

    // Instantiate the default detector which uses a threshold to determines anomalies.
    // Anomalies are data points that have a score above the threshold (2.5 in this case).
    anomalies := anomalia.NewDetector().WithThreshold(2.5).GetAnomalies(timeSeries)

    // Iterate over detected anomalies and print their exact timestamp and value.
    for _, anomaly := range anomalies {
        fmt.Println(anomaly.Timestamp, ",", anomaly.Value)
    }
}
```

And another example to check if two time series have a relationship or correlated:

```go
package main

import "github.com/amrfaissal/anomalia"

func main() {
    a := loadTimeSeriesA()
    b := loadTimeSeriesB()

    // If the time series data points do not follow a certain distribution
    // we use the Spearman correlator.
    coefficient := anomalia.NewCorrelator().WithTimeSeries(a, b).
    	WithMethod(anomalia.SpearmanRank, nil).
    	Run()

    // If the coefficient is above a certain threshold (0.7 for example), we consider
    // the time series correlated.
    if coefficient < 0.7 {
        panic("no relationship between the two time series")
    }
}
```

If the correlation algorithm accepts any additional parameters (see different implementations), you can pass them as a
 `float64` slice to the `WithMethod(method, options)` method.

## TODO / Suggestions

- CLI tool for rapid experimentation
- More algorithms such as STL with LOESS
- Benchmarks
- Move to a separate organization

## Resources

TODO

## License

```text
Copyright 2019 Faissal Elamraoui

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
