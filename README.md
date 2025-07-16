# rust_and_go_sandbox
using this sandbox  to create a multi-threaded backend rust with go receiving data packet.
reference: https://tutorialedge.net/projects/chat-system-in-go-and-react/part-3-designing-our-frontend/

links for sync wait:
https://stackoverflow.com/questions/67013517/golang-channel-missing-some-values

quick curl examples:
curl -X POST --url http://localhost:8080/alert -H 'Content-Type:application/json' -d '{"user_id":"8bb2b166-358d-4da9-8113-66242e0e0921", "method":"USER_INFO", "msg":"take me away!!!!"}' --next --url http://localhost:8080/alert -d '{"user_id":"8bb2b166-358d-4da9-8113-66242e0e0921", "method":"USER_INFO", "msg":"take me away2!!!!"}' --next --url http://localhost:8080/alert -d '{"user_id":"8bb2b166-358d-4da9-8113-66242e0e0921", "method":"USER_INFO", "msg":"take me away3!!!!"}' --next --url http://localhost:8080/alert -d '{"user_id":"8bb2b166-358d-4da9-8113-66242e0e0921", "method":"USER_INFO", "msg":"take me away4!!!!"}'