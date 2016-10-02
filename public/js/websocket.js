var socket;

window.onload = function(){
init()
document.onkeyup = keymonitor
document.onclick = function() {
        f = document.getElementById("focus_for_keyboard");
        f.focus()
      };
}

/**
 * The method for keyCode catch and send.
 */
function keymonitor(e) {
    res = e.keyCode
    if (e.shiftKey && e.keyCode>=65 && e.keyCode<=90){ //A-Z
         res = e.keyCode
    }
    else if (e.keyCode>=65 && e.keyCode<=90){ //a-z
        res = e.keyCode + 32
    }
    else if (e.keyCode == 37){//Left arrow
        socket.send(27)//Escape
        socket.send(91)//[
        res = 68 //D
    }
    else if (e.keyCode == 39){//Right arrow
        socket.send(27)//Escape
        socket.send(91)//[
        res = 67 //C
    }
    else if (e.keyCode == 38){//Up arrow
        socket.send(27)//Escape
        socket.send(91)//[
        res = 65 //A
    }
    else if (e.keyCode == 40){//Down arrow
        socket.send(27)//Escape
        socket.send(91)//[
        res = 66 //B
    }
    else if (e.keyCode == 13){ //Enter
        socket.send(10) 
    }
    else if (e.shiftKey && e.keyCode==48) { //)
        res = 41
    }
    else if (e.shiftKey && e.keyCode==49){ //!
         res = 33
    }
    else if (e.shiftKey && e.keyCode==50){ //@
         res = 64
    }
    else if (e.shiftKey && e.keyCode==51){ //#
         res = 35
    }
    else if (e.shiftKey && e.keyCode==52){ //$
         res = 36
    }
    else if (e.shiftKey && e.keyCode==53){ //%
         res = 37
    }
    else if (e.shiftKey && e.keyCode==54){ //^
         res = 94
    }
    else if (e.shiftKey && e.keyCode==55){ //&
         res = 38
    }
    else if (e.shiftKey && e.keyCode==56){ //*
         res = 42
    }
    else if (e.shiftKey && e.keyCode==57){ //(
         res = 40
    }
    else if (e.shiftKey && e.keyCode==187){ //- Doesn't work in Firefox
         res = 43
    }
    else if (e.shiftKey && e.keyCode==189){ //+ Doesn't work in Firefox
         res = 95
    }
    else if (e.keyCode==187){ //= Doesn't work in Firefox
         res = 61
    }
    else if (e.keyCode==189){ //_ Doesn't work in Firefox
         res = 45
    }
    else if (e.shiftKey && e.keyCode == 192){ //~
       res = 126
    }
    else if (e.shiftKey && e.keyCode == 219){ //{
        res = 123
    }
    else if (e.shiftKey && e.keyCode == 221){ //}
        res = 125
    }
    else if (e.shiftKey && e.keyCode == 220){ //|
        res = 124
    }
    else if (e.shiftKey && e.keyCode == 222){ //"
         res = 34
    }
    else if (e.shiftKey && e.keyCode == 186){ //:
        res = 58
    }
    else if (e.shiftKey && e.keyCode == 191){ //?
        res = 63
    }
    else if (e.shiftKey && e.keyCode == 190){ //>
        res = 62
    }
    else if (e.shiftKey && e.keyCode == 188){ //<
        res = 60
    }
    else if (e.keyCode == 188){//,
        res = 44
    }
    else if  (e.keyCode == 190) { //.
        res = 46
    }
    else if (e.keyCode == 191){ ///
        res = 47
    }
    else if (e.keyCode == 186){ //; Doesn't work in Firefox
        res = 59
    }
    else if (e.keyCode == 222){ //'
        res = 39
    }
    else if (e.keyCode == 220){ //\
        res = 92
    }
    else if (e.keyCode == 221){ //]
        res = 93
    }
    else if (e.keyCode == 219){ //[
        res = 91
    }
    else if (e.keyCode == 192){ //`
        res = 96
    }
    socket.send(res)
}

/**
 * The method for socket initialization.
 */
function init() {
    term = document.getElementsByTagName("table")
    host = term[0].getAttribute("host")
    port = term[0].getAttribute("port")
    socket = new WebSocket("ws://" + host + ":" + port+ "/ws_terminal");
    socket.onopen = function(event) {
    }
    
    socket.onmessage = function(event) {
        term[0].innerHTML =  event.data
        window.scrollTo(0,document.body.scrollHeight)
    }
    
    socket.onclose = function(event) {
    }
    
    socket.onerror = function(event) {
    }
}

