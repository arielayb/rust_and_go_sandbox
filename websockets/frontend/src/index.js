import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

var socket = new WebSocket("ws://localhost:8080/ws");

let connect = cb => {
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
  };

  socket.onmessage = msg => {
    if (msg !== "") {
      let data =  {
        user_uuid: crypto.randomUUID(),
        method: "USER_UUID",
        message: "hello",
      };
      cb(data.message);
    }
  };

  socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = error => {
    console.log("Socket Error: ", error);
  };
};

let sendMsg = msg => {
  console.log("we got this!!!!!", msg);
  console.log("sending msg: ", msg);
  if (msg !== "") {
    let data =  {
      user_uuid: crypto.randomUUID(),
      method: "USER_UUID",
      message: msg
    };
    socket.send(JSON.stringify(data));
  }   
};

export { connect, sendMsg };