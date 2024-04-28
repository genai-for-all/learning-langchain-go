import { html, render, Component } from '../js/preact-htm.js'
import Header  from './Header.js'
import Content  from './Content.js'
import Footer  from './Footer.js'


class App extends Component {
    ref = null
    setRef = (dom) => this.ref = dom
    
    constructor(props) {
        super()
        this.state = {}
    }

    render() {
        return html`
        <div ref=${this.setRef}>
            <div class="section">
                <div class="container">
                    <${Header} class="hero is-small is-primary"/>
                </div>
            </div>
            <div class="section">
                <div class="container">
                    <${Content} class="section" />
                </div>
            </div>
            <div class="section">
                <div class="container">
                    <${Footer} class="footer"/>
                </div>
            </div>
        </div>
        `
    }
}
//render(html`<${App}/>`, document.getElementById("app"));
//<${Title} subtitle="🚧 this is a wip"/>

export default App


