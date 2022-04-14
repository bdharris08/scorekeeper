# Example
An example usage of ScoreKeeper

### Usage
```sh
cd example/memory
go build
./memory
```

Expected output:
```sh
./example
[{"action":"hop","avg":60.6},{"action":"skip","avg":61.2},{"action":"jump","avg":61.8}]
[{"action":"hop","avg":62.6551724137931},{"action":"skip","avg":70.3076923076923},{"action":"jump","avg":79.69565217391305}]
[{"action":"hop","avg":64.85714285714286},{"action":"skip","avg":73.04},{"action":"jump","avg":91.2}]
[{"action":"hop","avg":60.6},{"action":"skip","avg":65.42857142857143},{"action":"jump","avg":70.84615384615384}]
[{"action":"hop","avg":60.6},{"action":"skip","avg":63.241379310344826},{"action":"jump","avg":68.33333333333333}]
[{"action":"hop","avg":60.6},{"action":"skip","avg":61.2},{"action":"jump","avg":61.8}]
[{"action":"skip","avg":61.2},{"action":"jump","avg":61.8},{"action":"hop","avg":60.6}]
```
