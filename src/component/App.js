import React from 'react'
import { BrowserRouter as Router, Route } from 'react-router-dom'

import About from './About'
import Calendar from './Calendar'
import Event from './Event'
import Menu from './Menu'

class App extends React.Component {
  render() {
    return (
      <Router>
        <div>
          <Menu />

          <Route path="/" exact={true} component={Calendar} />
          <Route path="/about" exact={true} component={About} />
          <Route path="/e/:id" component={Event} />
        </div>
      </Router>
    )
  }
}
export default App
