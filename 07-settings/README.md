# Settings

## Like a guardrail
```
ollama.WithPredictRepeatLastN(30)
OLLAMA_BASE_URL=http://localhost:11434 LLM=deepseek-coder go run main.go
```

> default
```
ollama.WithPredictRepeatLastN(64)
LLM=deepseek-coder go run main.go
```

> deactivate (and truncate?)
```
ollama.WithPredictRepeatLastN(0)
LLM=deepseek-coder go run main.go
```

## Stop Words

```
llms.WithStopWords([]string{"console"})
LLM=deepseek-coder go run main.go
```

## Temperature

```
llms.WithTemperature(0.5)
LLM=deepseek-coder go run main.go
```
> -> 0 more "creative"