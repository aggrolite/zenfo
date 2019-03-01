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

  onScroll(e) {
    e.preventDefault()
    const {
      loading,
      events,
    } = this.state;


    if (loading) return;

    var scrollTop = (document.documentElement && document.documentElement.scrollTop) || document.body.scrollTop;
    var scrollHeight = (document.documentElement && document.documentElement.scrollHeight) || document.body.scrollHeight;
    var clientHeight = document.documentElement.clientHeight || window.innerHeight;
    var scrolledToBottom = Math.ceil(scrollTop + clientHeight) >= scrollHeight;

    let start = events.slice(-1)[0].start
    if (scrolledToBottom) {
      this.loadEvents(start);
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
        //var newEvents = this.state.events.slice().concat(resp[0])
        this.setState({
          events: [...this.state.events, ...resp[0] ],
          //venues: resp[1],
          loading: false,
        });
      })
  }

  componentDidMount() {
    window.addEventListener('scroll', this.onScroll)
    this.loadEvents()
  }

  componentWillUnmount() {
    window.removeEventListener('scroll', this.onScroll)
  }

  render() {

    const { events, loading } = this.state;

	return (
      <div className="calendar">
        <ul className="ul" >
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
          {loading &&
            <li>Loading...</li>
          }
        </ul>

      </div>
    )
  }
}

export default Calendar
