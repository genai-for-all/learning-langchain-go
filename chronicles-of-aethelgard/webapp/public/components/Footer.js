import { html, render, Component } from '../js/preact-htm.js'

class Footer extends Component {
    constructor(props) {
        super()
        this.state = {}
    }
    
    render() {
        return html`
          <div class="content has-text-centered">
            <p>
              🚧 this is a work in progress
            </p>
          </div>
        `
    }
}

export default Footer
