import React from 'react'
import fetch from 'isomorphic-fetch'
import PropTypes from 'prop-types';

class Event extends React.Component {
  componentDidMount() {
    const { id } = this.props.match.params
    fetch(`/api/events?id=${id}`)
      .then((events) => {
        this.setState(() => ({ events }))
      })
  }
  render() {
    return (
      <div>`blabla{this.props.match.params.id}`</div>
    )
  }
}

Event.propTypes = {
  match: PropTypes.shape({
    params: PropTypes.shape({
      id: PropTypes.string,
    })
  }),
}

export default Event
