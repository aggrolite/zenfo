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
    fetch("http://localhost:8081/api/events")
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
      <div className="calendar">
        <ul>
          {events.map(event =>
            <li className="event" key={event.id}>
              <div className="date">
                <div className="day">
                  <Moment format="D">{event.start}</Moment>
                </div>
                <div className="month">
                  <Moment format="MMM">{event.start}</Moment>
                </div>
                <div className="year">
                  <Moment format="YYYY">{event.start}</Moment>
                </div>
              </div>
              <div className="overview">
                <Link to={"/e/" + event.id} className="title">
                  {event.name}
                </Link>
                <div className="venue">{event.venue.name}</div>
              </div>
            </li>
          )}
        </ul>

      </div>
    )
  }
}

export default Calendar
