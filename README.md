# anomalia

`anomalia` is zero-dependency Go library for time series data analysis.

It supports anomaly detection and correlation. The API is simple and configurable in the sense that you can choose which algorithm suits your needs for anomaly detection and correlation.

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
import (
    "fmt"
    "github.com/amrfaissal/anomalia"
)

func main() {
    // Load the time series from an external source.
    // It returns an instance of TimeSeries struct which holds the timestamps and their values.
    timeSeries := loadTimeSeriesFromAnExternalSource()

    // Instantiate the default detector which uses a threshold to determines anomalies.
    // Anomalies are data points that have a score above the threshold (2 in this case).
    anomalies := anomalia.NewDetector(2).GetAnomalies(timeSeries)

    // Iterate over detected anomalies and print their exact timestamp and value.
    for _, anomaly := range anomalies {
        fmt.Println(anomaly.Timestamp, ",", anomaly.Value)
    }
}
```

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
