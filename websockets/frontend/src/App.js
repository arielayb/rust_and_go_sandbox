// import logo from './logo.svg';
// import './App.css';

// App.js
import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg } from "./index"
import Header from './components/header/Header';

class App extends Component {
  constructor(props) {
    super(props);
    connect();
  }

  send() {
    console.log("hello");
    sendMsg("hello");
  }

render() {
  return (
    <div className="App">
      <Header />
      <button onClick={this.send}>Hit</button>
    </div>
  );
}
}

export default App;
