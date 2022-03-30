## How to use

```go
promx.NewServe(
	viper.GetString("PROMETHEUS_NAME"),
	viper.GetString("PROMETHEUS_PATH"),
	viper.GetString("PROMETHEUS_PORT"), 
).Start()
```