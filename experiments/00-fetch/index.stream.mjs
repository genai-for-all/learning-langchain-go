

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
        stream: true
      })
    })

    const reader = response.body.getReader()

    let reading = true
    let content = []
    let responseText = ""
    let result = null
    while (reading) {
      const { done, value } = await reader.read()

      if (done) {
        reading = false
        result = content.slice(-1)[0]
        result.response = responseText
        
      }
      const decodedValue = new TextDecoder().decode(value)

      if (decodedValue !== "") {
        let jsondecodedValue = JSON.parse(decodedValue)
        content.push(jsondecodedValue)
        responseText += jsondecodedValue.response
        process.stdout.write(jsondecodedValue.response);
      }

    }

    //console.log(content)

    console.log(result)
    
    
} catch(error) {
    console.log(error)
}