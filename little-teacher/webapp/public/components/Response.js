import { html, render, Component } from '../js/preact-htm.js'

class Response extends Component {
    setRefContent = (dom) => this.refContent = dom
    setRefHeader = (dom) => this.refHeader = dom

    constructor(props) {
        super()
        this.state = {}

        this.changeMsgHeaderText = this.changeMsgHeaderText.bind(this)
        this.changeResponseContent = this.changeResponseContent.bind(this)

    }

    changeMsgHeaderText(text) {
        this.refHeader.innerHTML = text
    }

    changeResponseContent(text) {
        //this.setState({ResponseContent: text})
        this.refContent.innerHTML = text
    }

    render() {
        return html`
        <div class="content">
            <article class="message is-dark">
                <div class="message-header">
                    <p id="msg_header" ref=${this.setRefHeader}>🤖 Answer:</p>
                </div>

                <div id="txt_response" class="message-body" ref=${this.setRefContent}>
                </div>

                <div class="message-footer">
                </div>
            </article>
        </div>
        <div class="content">
        </div>
        `
    }
}

export default Response
