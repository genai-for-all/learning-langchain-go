

const { LLM = 'gemma' } = process.env;

const ollmaBaseUrl = process.env["OLLAMA_BASE_URL"]

try {
    const response = await fetch(`${ollmaBaseUrl}/api/generate`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        model: LLM,
        prompt: "Who is James T Kirk?",
        stream: false
      })
    })

    const content = await response.json();

    console.log(content)
} catch(error) {
    console.log(error)
}