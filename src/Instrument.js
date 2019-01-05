import React, { Component } from 'react'
import { Grid, Row, Col, Button } from "react-bootstrap"
import Tone from "tone"

import './Instrument.css'

function Key({ note, attack, release }) {
  return (
    <Button onMouseDown={attack} onMouseUp={release} bsSize="large" block>
      {note}
    </Button>
  )
}

class Instrument extends Component {
  constructor(props) {
    super(props)
    this.synth = new Tone.Synth({
			"oscillator" : {
				"type" : "amtriangle",
				"harmonicity" : 0.5,
				"modulationType" : "sine"
			},
			"envelope" : {
				"attackCurve" : 'exponential',
				"attack" : 0.05,
				"decay" : 0.2,
				"sustain" : 0.2,
				"release" : 1.5,
			},
			"portamento" : 0.05
		}).toMaster()
  }

  attack = note => {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
      if (this.socket.readyState === WebSocket.CONNECTING) {
        return
      }
      this.socket = new WebSocket(this.props.socketUrl)
      this.socket.onmessage = e => {
        const d = JSON.parse(e.data)
        this.synth.triggerAttackRelease(`${d.note}${this.props.octave}`, "8n")
      }
    }
    this.socket.send(JSON.stringify({ note }))
  }

  release = () => {
    this.synth.triggerRelease()
  }

  componentDidMount() {
    this.socket = new WebSocket(this.props.socketUrl)
    this.socket.onmessage = e => {
      const d = JSON.parse(e.data)
      this.synth.triggerAttackRelease(`${d.note}${this.props.octave}`, "8n")
    }
  }

  componentWillUnmount() {
    if (this.socket.readyState !== WebSocket.CLOSED) {
      this.socket.close()
    }
    this.socket = null
  }

  render() {
    const buttons = this.props.scale.map(note => {
      return (
        <Col key={note} xs={12} md={6} lg={3}>
          <Key note={note} attack={() => this.attack(note)} release={this.release} />
        </Col>
      )
    })
    return (
      <Grid>
        <Row>
          {buttons}
        </Row>
      </Grid>
    )
  }
}

export default Instrument;