window.onload = function(){
    imgs = [];
    
    streamWidget = document.getElementsByClassName("video-stream");
    
    for (var i = 0; i < streamWidget.length; i++) {
        imgs.push(streamWidget[i].getElementsByTagName("img")[0]);
    }
    streams = removeDuplicatesAndMakeObj(imgs);
    
    for (key in streams) {
        var res = initPanelElements(streams[key]);
        init(key, streams[key], res[0], res[1]);
    }
}

/**
 * The method for initialization bottom-panel and top-panel elements.
 *
 * @param imgarr The key value of streams object.
 * @return res Array of output and pauseBtns arrays.
 */
function initPanelElements(imgarr) {
    var output = [];
    var pauseBtns = [];
    
    for( var i = 0; i < imgarr.length; i++) {
        output.push(imgarr[i].parentNode.getElementsByClassName("info-lbl")[0]);
        pauseBtns.push(imgarr[i].parentNode.getElementsByClassName("pause-btn")[0]);
    }
    return [output, pauseBtns];
}

/**
 * The method for sockets's initialization.
 *
 * @param addr The key of streams object.
 * @param imgarr The key value of streams object.
 * @param output Array of span elements.
 * @param pauseBtns Array of input button elements.
 */
function init(addr, imgarr, output, pauseBtns) {
    var state = false;
    var socket = new WebSocket(addr);
    
    for (var i = 0; i < imgarr.length; i++) {
        pauseBtns[i].onclick = function() {
            if (state == false) {
                state = true;
                socket.close();
            }
            else {
                init(addr, imgarr, output, pauseBtns);
            }
        }
        
    }
    
    socket.onopen = function(event) {
        for (var i = 0; i < imgarr.length; i++) {
            output[i].style.color = "green";
            output[i].textContent = "[connected]";
            pauseBtns[i].value = '\u2590\u2590';
            
        }
    }
    
    socket.onmessage = function(event) {
        for (var i = 0; i < imgarr.length; i++)
            imgarr[i].src = event.data;
    }
    
    socket.onclose = function(event) {
        if(event.wasClean)
            for (var i = 0; i < imgarr.length; i++) {
                output[i].style.color = "red";
                output[i].textContent = "[disconnected]";
                pauseBtns[i].value = '\u25BA';
            }
        else
            for (var i = 0; i < imgarr.length; i++) {
                output[i].style.color = "red";
                output[i].textContent = "[server error: disconnected]";
                pauseBtns[i].value = '\u25B7';
            }
    }
    
    socket.onerror = function(event) {
    }
}

/**
 * The function removes duplicates addresses and makes the object (associative array).
 *
 * @param images Array of img elements.
 * @return obj The object {addr:[],...}, keys are addresses, values are arrays
 * of img elements.
 */
function removeDuplicatesAndMakeObj(images) {
    obj = {};
    
    for (var i = 0; i < images.length; i++) {
        var host = images[i].parentNode.getAttribute("host");
        var port = images[i].parentNode.getAttribute("port");
        var imgAddr = "ws://" + host + ":" + port+ "/ws";
        if (imgAddr in obj)
            obj[imgAddr].push(images[i]);
        else
            obj[imgAddr] = [images[i]];
    }
    return obj;
}

