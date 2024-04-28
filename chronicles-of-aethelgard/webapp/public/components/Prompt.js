import { html, render, Component } from '../js/preact-htm.js'

const system_message = `You are the dungeon master, 
expert at interpreting and answering questions based on provided sources.
Using the provided context, answer the user's question 
to the best of your ability using only the resources provided. 
Be verbose!`

const human_message = `Who are the players of Chronicles of Aethelgard, give details for every player?`

class Prompt extends Component {
    //ref = null
    setRef = (dom) => this.ref = dom
    //refUploader = null
    setRefUploader = (dom) => this.refUploader = dom
    setRefUploadFile = (dom) => this.refUploadFile = dom

    constructor(props) {
        super()

        this.state = {
          aborter: new AbortController(),
          humanMessage: human_message,
          systemMessage: system_message
        }

        this.btnSubmitOnClick = this.btnSubmitOnClick.bind(this)
        this.btnStopOnClick = this.btnStopOnClick.bind(this) 
        /*   
        this.btnClearAnswerOnClick = this.btnClearAnswerOnClick.bind(this)
        this.btnPrintConversationOnClick = this.btnPrintConversationOnClick.bind(this)
        this.btnSelectOnClick = this.btnSelectOnClick.bind(this)
        this.btnUploadOnClick = this.btnUploadOnClick.bind(this)
        */

        this.txtSystemMessageOnChange = this.txtSystemMessageOnChange.bind(this)
        this.txtHumanMessageOnChange = this.txtHumanMessageOnChange.bind(this)
        
    }

    txtSystemMessageOnChange(e) {
        // triggered by the textarea of the system message
        this.setState({ systemMessage: e.target.value })
    }
    txtHumanMessageOnChange(e) {
        // triggered by the textarea of the human message
        this.setState({ humanMessage: e.target.value })
    }

    async btnSubmitOnClick() {
      console.log("🤓 Prompt: clicked")

      let responseText=""
      // 🫢 this is a hack
      var that = this

      try {
        let waitingTimer = setInterval(waitingMessage, 500)
        let waiting = true

        function waitingMessage() {
          const d = new Date()
          console.log("🤔 waiting", d.toLocaleTimeString())
          // 🫢 I use the hack here
          that.ref.dispatchEvent(
            new CustomEvent("waiting", {
              bubbles: true, 
              composed: true, 
              detail: { 
                text: "🤖 Answer: 🤔 computing " + d.toLocaleTimeString(),
                from: "prompt",  
              }
            })
          )
        }

        const response = await fetch("/prompt", {
          method: "POST",
          headers: {
            "Content-Type": "application/json;charset=utf-8",
          },
          body: JSON.stringify({
            question: this.state.humanMessage,
            systemMessage: this.state.systemMessage,
          }),
          signal: this.state.aborter.signal,
        })
  
        const reader = response.body.getReader()

        while (true) {
          const { done, value } = await reader.read()
  
          if (waiting) {
            clearInterval(waitingTimer)
            waiting = false
            // send message to the Content component
            // the Content component will then update the Answer component
            this.ref.dispatchEvent(
              new CustomEvent("waiting", {
                bubbles: true, 
                composed: true, 
                detail: { 
                  text: "🤖 Answer:",
                  from: "prompt",  
                }
              })
            )
          }
          
          if (done) {
            responseText = responseText + "\n"
            //this.ref_txtResponse.innerHTML = markdownit().render(responseText)
            this.ref.dispatchEvent(
              new CustomEvent("response", {
                bubbles: true, 
                composed: true, 
                detail: { 
                  text: markdownit().render(responseText),
                  from: "prompt",  
                }
              })
            )

            return
          }
          // Otherwise do something here to process current chunk
          const decodedValue = new TextDecoder().decode(value)
          console.log(decodedValue)
          responseText = responseText + decodedValue
          //this.ref_txtResponse.innerHTML = markdownit().render(responseText)
          this.ref.dispatchEvent(
            new CustomEvent("response", {
              bubbles: true, 
              composed: true, 
              detail: { 
                text: markdownit().render(responseText),
                from: "prompt",  
              }
            })
          )
        }
      } catch (error) {
        console.log("😡", error)
      }
    }

    btnStopOnClick() {
        console.log("🛑 Stop: clicked")
        this.state.aborter.abort()
    }
  
    render() {
        return html`
        <div ref=${this.setRef}>
          <div  class="field">
            <label class="label">
                <span class="tag is-info">System:</span>
            </label>
            <div class="control">
                <textarea id="system-message" 
                class="textarea is-family-code is-info" 
                rows="5"
                value=${this.state.systemMessage} 
                onInput=${this.txtSystemMessageOnChange}
                placeholder="System Message" />
            </div>

            <label class="label">
                <span class="tag is-primary">Prompt:</span>
            </label>
            <div class="control">
                <textarea id="human-message" 
                class="textarea is-family-code is-primary" 
                rows="3"
                value=${this.state.humanMessage} 
                onInput=${this.txtHumanMessageOnChange}
                placeholder="Ask me any question about this game"/>
            </div>
          </div> 

          <div class="content">
            <div class="field is-grouped">
                <div class="control">
                    <button class="button is-link is-small" onclick=${this.btnSubmitOnClick}>Submit Question</button>
                </div>
                <div class="control">
                    <button class="button is-link is-danger is-small" onclick=${this.btnStopOnClick}>Stop</button>
                </div>

            </div>
          </div>
        </div>
        `
    }
}

export default Prompt
