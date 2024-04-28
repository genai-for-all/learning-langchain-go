import { html, render, Component } from '../js/preact-htm.js'

class Header extends Component {
    constructor(props) {
        super()
        this.state = {}
    }
    render() {
        return html`
          <div class="hero-body">
            <p class="title is-3 has-text-centered">
              🩵 GoLang 🐳 Docker 🦙 GenAI Stack 🦜🔗 (📝 RAG)
            </p>
            <p class="title is-3 has-text-centered">
              🧙‍♂️ Chronicles of Aethelgard 🧌
            </p>
          </div>
        `
    }
}

export default Header

