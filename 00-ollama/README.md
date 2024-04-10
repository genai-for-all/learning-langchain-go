# Ollama API

### Query Ollama

> https://github.com/ollama/ollama/blob/main/docs/api.md


```bash
curl http://host.docker.internal:11434/api/generate -d '{
  "model": "gemma",
  "prompt": "Who is James T Kirk?",
  "stream": false
}'
```

#### Use it with jq
```bash
curl http://host.docker.internal:11434/api/generate -d '{
  "model": "gemma",
  "prompt": "Who is James T Kirk?",
  "stream": false
}' | jq -c '{ response, context }'
```


> copy paste the context from the answer


```bash
curl http://host.docker.internal:11434/api/generate -d '{
  "model": "gemma",
  "prompt": "Who is his best friend?",
  "stream": false,
  "context": []
}' | jq -c '{ response }'
```


#### Use streaming
```bash
curl http://host.docker.internal:11434/api/generate -d '{
  "model": "gemma",
  "prompt": "Who is James T Kirk?"
}' 
```
