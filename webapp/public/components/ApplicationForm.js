import { html, render, Component } from '../js/preact-htm.js'
import Prompt  from './Prompt.js'
import Response  from './Response.js'
import ApplicationTitle from './ApplicationTitle.js'


class ApplicationForm extends Component {
  render() {
    return html`
    <${ApplicationTitle}/>
    <hr></hr>
    <${Prompt}/>
    <${Response}/>
    `
  }
}

export default ApplicationForm
