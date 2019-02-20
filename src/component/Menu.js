import React from 'react'
import { Link } from 'react-router-dom'

class Menu extends React.Component {
  render() {
	return (
     <div className="menu">
       <ul>
         <li>zenfo.info</li>
         <li><Link to="/">home</Link></li>
         <li><Link to="/about">about</Link></li>
       </ul>
     </div>
    )
  }
}
export default Menu
