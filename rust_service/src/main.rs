use std::vec;
use std::thread;
use tungstenite::{connect, Message};

trait DataPacketInterface {
    fn get_packet(&self) -> &String;
    fn set_packet(&mut self, data: String);
    fn recv_msg(&self);
}


pub struct DataPacket {
    msg: String,
    size: u32,
}

impl DataPacketInterface for DataPacket {
    fn get_packet(&self) -> &String{
        return &self.msg;
    }

    fn set_packet(&mut self, data: String) {
       self.msg = String::from(data); 
    }

    fn recv_msg(&self) {
        let (mut socket, response) = connect("ws://localhost:8080/ws").expect("can't connect");

        println!("Connected to the server");
        println!("Response HTTP code: {}", response.status());
        println!("Response contains the following headers:");

        for (header, _) in response.headers() {
            println!("* {header}");
        }

        socket.send(Message::Text("Hello WebSocket".into())).unwrap();
        loop {
            let msg = socket.read().expect("Error reading message");
            println!("Received: {msg}");
        }
            
        // socket.close(None);
    }
}

fn main() {
    println!("Starting rust service....");

    let data = DataPacket { 
        msg: String::from(""), 
        size:0 
    };

    thread::spawn(move || {  
        data.recv_msg();
    }).join().expect("recv msg thread completed.")
    
}
