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

curl http://host.docker.internal:11434/api/generate -d '{
  "model": "gemma",
  "prompt": "Who is his best friend?",
  "stream": false,
  "context": [106,1645,108,6571,603,6110,584,29592,235336,107,108,106,2516,108,14443,584,235265,29592,603,476,60193,3285,578,573,7233,72554,576,573,3464,7081,44639,4100,235265,1315,603,476,8995,1651,31863,578,51018,3836,604,926,12881,235269,20803,3403,235269,578,192038,71502,235265,107,108]
}' | jq -c '{ response }'





#### Use streaming
```bash
curl http://host.docker.internal:11434/api/generate -d '{
  "model": "gemma",
  "prompt": "Who is James T Kirk?"
}' 
```
