import { html, render, Component } from '../js/preact-htm.js'

class Answer extends Component {
    ref = null
    setRef = (dom) => this.ref = dom

    setRefContent = (dom) => this.refContent = dom
    
    constructor(props) {
        super()
        this.state = {
            MsgHeader: "🤖 Answer:",
            ResponseContent: "..."
        }

        this.changeMsgHeaderText = this.changeMsgHeaderText.bind(this)
        this.changeResponseContent = this.changeResponseContent.bind(this)

    }

    changeMsgHeaderText(text) {
        this.setState({MsgHeader: text})
    }
    changeResponseContent(text) {
        //this.setState({ResponseContent: text})
        this.refContent.innerHTML = text
    }

    render() {
        return html`
            <div ref=${this.setRef} class="content">
                <article class="message is-dark">
                    <div class="message-header">
                        <p id="msg_header">${this.state.MsgHeader}</p>
                    </div>

                    <div id="txt_response" class="message-body" ref=${this.setRefContent}>
                        ${this.state.ResponseContent}
                    </div>
                    <!--
                    <div class="is-family-primary">
                    </div>
                    -->

                    <div class="message-footer">
                    </div>
                </article>
            </div>
        `
    }
}

export default Answer

