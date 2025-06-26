// import logo from './logo.svg';
// import './App.css';

// App.js
import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg } from "./index"
import Header from './components/header/Header';
import ChatHistory from './components/chatHistory/ChatHistory';
// import TextTable from './components/textTable/TextTable';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      chatHistory: []
    }
  }

  componentDidMount() {
    connect((msg) => {
      console.log("New Message")
      this.setState(prevState => ({
        chatHistory: [...this.state.chatHistory, msg]
      }))
      // this.setState(prevState => ({
      //   textTable: [...this.state.textTable, msg]
      // }))
      console.log(this.state);
    });
  }

  send() {
    console.log("hello");
    sendMsg("hello");
  }

render() {
  return (
    <div className="App">
      <Header />
      <ChatHistory chatHistory={this.state.chatHistory} />
      {/* <TextTable textTable={this.state.textTable}/> */}
      <button onClick={this.send}>Hit</button>
    </div>
  );
}
}

export default App;
