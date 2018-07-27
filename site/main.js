const drawButton = "draw-button"
const roomCodeEntry = "room-code-entry"
const createRoom = "create-room-button"
const roleOutput = "role-output"
const numOfPlayers = "num-of-players"

apiURL = "https://j2s9y2val3.execute-api.us-east-1.amazonaws.com/dev/"

var roomCode = ''
var playersNum = 5
var room = {}

document.getElementById(createRoom).addEventListener("click", function(){
    playersNum = document.getElementById(numOfPlayers).value
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
    console.log(roomCode)
    fetch(`${apiURL}/dealCard?code=${roomCode}`, {
        method: 'PUT'
    }).then((res) => {
        console.log(res.status)
        console.log(res.text().then((val) => {
            console.log(val)
            document.getElementById(roleOutput).innerHTML = `${val}`
        }))
    }).catch((err) => {
        console.log(err)
        document.getElementById(roleOutput).innerHTML = `Player limit reached.`
    })
});
