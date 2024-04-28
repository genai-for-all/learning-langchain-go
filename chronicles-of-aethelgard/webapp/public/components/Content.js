import { html, render, Component } from '../js/preact-htm.js'

import Prompt from './Prompt.js'
import Answer from './Answer.js'

class Content extends Component {
    refPrompt = null
    refAnswer = null
    ref = null
    setRefPrompt = (dom) => this.refPrompt = dom
    setRefAnswer = (dom) => this.refAnswer = dom
    setRef = (dom) => this.ref = dom

    constructor(props) {
        super()
        this.state = {}
    }

    componentDidMount() {
      //console.log(this.refPrompt.ref)
      //console.log(this.refAnswer.ref)

      this.ref.addEventListener('waiting', (e) => {
        if (e.detail.from === "prompt") {
          this.refAnswer.changeMsgHeaderText(e.detail.text)
        }
      })

      this.ref.addEventListener('response', (e) => {
        if (e.detail.from === "prompt") {
          this.refAnswer.changeResponseContent(e.detail.text)
        }
      })
    }

    render() {
        return html`
          <div ref=${this.setRef}>
            <${Prompt} ref=${this.setRefPrompt} message="this is a message"/>
            <hr></hr>
            <${Answer} ref=${this.setRefAnswer}/>
          </div>
        `
    }
}

export default Content
