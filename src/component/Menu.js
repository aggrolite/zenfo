import React from 'react'
import { BrowserRouter as Router, Route, Link, Switch } from 'react-router-dom'

import About from './About'
import Calendar from './Calendar'
import Event from './Event'

class Menu extends React.Component {
  render() {
    return (
      <Router>
        <div className="menu">
          <Link to="/">Home</Link>
          <Link to="/about">About</Link>

          <Switch>
            <Route path="/" exact={true} component={Calendar} />
            <Route path="/about" exact={true} component={About} />
            <Route path="/e/:id" component={Event} />
          </Switch>
        </div>
      </Router>
    )
  }
}
export default Menu
