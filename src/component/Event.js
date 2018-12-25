import React from 'react'
import fetch from 'isomorphic-fetch'
import PropTypes from 'prop-types'

class Event extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      events: [],
    };
  }
  componentDidMount() {
    const { id } = this.props.match.params
    fetch(`/api/events?id=${id}`)
      .then((events) => {
        this.setState(() => ({ events }))
      })
  }
  render() {

    const { events } = this.state;
 
    return (
      <div className="event">
        <h1>{events[0].name}</h1>
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
  history: PropTypes.string,
}

export default Event
