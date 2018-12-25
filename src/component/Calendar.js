import React from 'react'
import Moment from 'react-moment'
import fetch from 'isomorphic-fetch'
import { Link } from 'react-router-dom'

//import Event from './Event'

//const Event = ({ match }) => (
//  <div>
//    <h2>{match.params.id}</h2>
//  </div>
//)

class Calendar extends React.Component {

  constructor(props) {
    super(props);

    this.state = {
      events: [],
      loading: false,
      error: null,
    };
  }

  componentDidMount() {
    fetch("/api/events")
		.then(response => {
           if (response.ok) {
             return response.json()
           } else {
             throw new Error('Failed to fetch events')
           }
        })
		.then(data => this.setState({ events: data }))
		.catch(error => this.setState({ error: error, loading: false }));
  }
  render() {

    const { events, loading } = this.state;

    if (loading) {
      return <p>Loading...</p>
    }

	return (
      <div>
        <ul>
          {events.map(event =>
            <li key={event.id}>
              <Moment format="D MMM YYYY">
                {event.start}
              </Moment>
              <Link to={"/e/" + event.id}>
                {event.name} @ {event.venue.name}
              </Link>
            </li>
          )}
        </ul>

      </div>
    )
  }
}

export default Calendar
