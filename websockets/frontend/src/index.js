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

var dataInit = { 
  user_id: crypto.randomUUID(),
  method: "USER_INFO",
  msg: "",
  global: true,
};

let connect = cb => {
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
    socket.send(JSON.stringify(dataInit));
  };

  socket.onmessage = msg => {
    if (msg !== "") {
      dataInit.alert_msg = msg;
      cb(dataInit.alert_msg);
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
  console.log("sending msg: ", msg);
  if (msg !== "") {
    // let data =  {
    //   user_id: crypto.randomUUID(),
    //   method: "USER_UUID",
    //   alert_msg: msg
    // };
    socket.send(JSON.stringify(msg));
  }   
};

export { connect, sendMsg };
