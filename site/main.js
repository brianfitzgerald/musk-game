const drawButton = "draw-button"
const roomCodeEntry = "room-code-entry"
const createRoom = "create-room-button"
const roleOutput = "role-output"

apiURL = "https://j2s9y2val3.execute-api.us-east-1.amazonaws.com/dev/"

var roomCode = ''
var playersNum = 5
var room = {}

document.getElementById(createRoom).addEventListener("click", function(){
    fetch(`${apiURL}/createRoom?players=${playersNum}`, {
        method: 'POST'
    }).then((res) => {
        console.log(res.json().then((val) => {
            console.log(val)
            room = val
            console.log(room.Code)
            document.getElementById(roleOutput).innerHTML = `Your Room Code: ${room.Code}`
        }))
    })
});

document.getElementById(drawButton).addEventListener("click", function(){
    roomCode = document.getElementById(roomCodeEntry).value
    fetch(`${apiURL}/dealCard?code=${roomCodeEntry}`, {
        method: 'PUT'
    }).then((res) => {
        console.log(res.json().then((val) => {
            console.log(val)
            document.getElementById(roleOutput).innerHTML = `${val}`
        }))
    })
});
