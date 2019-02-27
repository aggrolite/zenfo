import React from 'react'
import Moment from 'react-moment'
import fetch from 'isomorphic-fetch'
import Promise from 'es6-promise'
import { Link } from 'react-router-dom'
import API from '../constants/api'

//import Event from './Event'

//const Event = ({ match }) => (
//  <div>
//    <h2>{match.params.id}</h2>
//  </div>
//)

class Calendar extends React.Component {

  constructor(props) {
    super(props)
    this.onScroll = this.onScroll.bind(this)

    this.state = {
      events: [],
      venues: [],
      loading: false,
      error: null,
    }
  }

  onScroll() {

    const {
      loading,
      events,
    } = this.state;

    if (loading) return;
    var lastLi = document.querySelector("ul > li.event:last-child");
    var lastLiOffset = lastLi.offsetTop + lastLi.clientHeight;
    var pageOffset = window.pageYOffset + window.innerHeight;
    if (pageOffset > lastLiOffset) {
      let start = events.slice(-1)[0].start
      this.loadEvents(start)
    }

  }

  loadEvents(after) {

    let eventsURL = API.BASE_URL + "/events"
    if (after) {
      eventsURL = eventsURL + "?after=" + after
    }
      this.setState({loading: true});
      Promise.all([
        fetch(eventsURL).then((resp) => resp.json()),
        fetch(API.BASE_URL + "/venues").then((resp) => resp.json())
      ]).then(resp => {
        var newEvents = this.state.events.slice().concat(resp[0])
        this.setState({
          events: newEvents,
          venues: resp[1],
          loading: false
        });
      })
  }

  componentDidMount() {
    window.addEventListener('scroll', this.onScroll, false);
    this.loadEvents()
  }
  render() {

    const { events, venues, loading } = this.state;

    if (loading) {
      return <p>Loading...</p>
    }

	return (
      <div className="calendar">
        <div className="venues">
          {venues.map(venue =>
              <div className="venue" key={venue.id}>
                <input type="checkbox" checked={true}/>
                <label htmlFor={venue.id}>{venue.name}</label>
              </div>
          )}
        </div>
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
