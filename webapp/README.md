


HTTP_PORT=8888 LLM=deepseek-coder:1.3b OLLAMA_BASE_URL=http://host.docker.internal:11434 docker compose --profile container up --build

HTTP_PORT=8888 LLM=deepseek-coder:1.3b-instruct OLLAMA_BASE_URL=http://host.docker.internal:11434 docker compose --profile container up --build