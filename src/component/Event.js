import React from 'react'
import fetch from 'isomorphic-fetch'
import PropTypes from 'prop-types'

class Event extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      event: {},
      loading: false,
      error: null,
    };
  }
  componentDidMount() {
    const { id } = this.props.match.params
    fetch(`https://zenfo.info/api/events?id=${id}`)
      .then(response => {
        if (response.ok) {
          return response.json()
        } else {
          throw new Error('Failed to fetch event')
        }
      })
      .then(data => this.setState({ event: data.length > 0 ? data[0] : {} }))
      .catch(error => this.setState({ error: error, loading: false }))
  }
  render() {

    const { event } = this.state;

    //console.log(events)

    return (
      <div className="event">
        <h1>{event.name}</h1>
        <div className="blurb">
          {event.blurb}
        </div>
        <div className="desc">
          {event.desc}
        </div>
      </div>
    )
  }
}

Event.propTypes = {
  match: PropTypes.shape({
    params: PropTypes.shape({
      id: PropTypes.string,
    })
  }),
  history: PropTypes.object,
}

export default Event
