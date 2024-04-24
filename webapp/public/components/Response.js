import { html, render, Component } from '../js/preact-htm.js'

class Response extends Component {
    constructor(props) {
        super()
        this.state = {}
      }
    render() {
        return html`
        <div class="content">
            <article class="message is-dark">
                <div class="message-header">
                    <p id="msg_header">🤖 Answer:</p>
                </div>

                <div id="txt_response" class="message-body">
                </div>
                <!--
                <div class="is-family-primary">
                </div>
                -->

                <div class="message-footer">
                </div>
            </article>
        </div>
        <div class="content">
            🚧 this is a work in progress
        </div>
        `
    }
}

export default Response
