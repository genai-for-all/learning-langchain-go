import { html, render, Component } from '../js/preact-htm.js'
import Prompt  from './Prompt.js'
import Title from './Title.js'

class Application extends Component {
  render() {
    return html`
    <${Title}/>
    <hr></hr>
    <${Prompt}/>
    `
  }
}

export default Application
