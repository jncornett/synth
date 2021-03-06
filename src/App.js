import React, { Component } from 'react'

import logo from './logo.svg';
import './App.css';

import { pentatonicScale } from "./music"
import Instrument from "./Instrument"


class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <h1 className="App-title">Welcome to React</h1>
        </header>
        <Instrument
          scale={pentatonicScale}
          octave={4}
          socketUrl="ws://octo.local:4998/cmd"
        />
      </div>
    );
  }
}

export default App;
