window.onload = function(){
    menu = document.getElementById('menu');
	menu.isshowed = false;
	workspace = document.getElementById('workspace');
	initGestures()
}

function menuShow() {
	if (!('isshowed'  in menu) || menu.isshowed == false) {
		menu.style.transform = "translateX(0px)";
		workspace.style.transform = "translateX(350px)";
		menu.isshowed = true;
	}
	
	else {
		menu.style.transform = "translateX(-300px)";
		workspace.style.transform = "translateX(50px)";
		menu.isshowed = false;
	}	
}

var timeout;
function startProgram() {
    p = document.getElementsByTagName("p");
    var i = 0, len = p.length;
    function f() {
        if(i>0) {
            p[i-1].style.backgroundColor = "#CFD8DC";
            p[i-1].style.color = "#455A64";
        }
        p[i].style.backgroundColor = "#009688"
        p[i].style.color = "#ffffff";
        i++;
        if( i < len ){
            timeout = setTimeout( f, 1000 );
        }
    }
    f();
}

function stop() {
    clearTimeout(timeout);
}

function initZoom(iframe, innerWrapper) {
    widget = iframe.contentDocument.getElementById(innerWrapper);
    var mc = new Hammer.Manager(widget);
    var pinch = new Hammer.Pinch();
    mc.add([pinch]);
    mc.on( 'pinch', function(e) {
        real_scale = e.scale * 0.5;
        scale_str = "scale(" + real_scale + ", " + real_scale +")";
        if (e.scale < 1.0) {
            // User moved fingers closer together
            iframe.style.webkitTransform = scale_str;
        } else if (e.scale > 1.0) {
            // User moved fingers further apart
            iframe.style.webkitTransform = scale_str;
        }
    });
    
}

function addAxisMovementsWidget() {
	iframeAxisMovements = document.createElement('iframe');
	iframeAxisMovements.src = "axis-movements-widget.html";
    iframeAxisMovements.style.webkitTransform = "scale(0.5, 0.5)";
    iframeAxisMovements.frameBorder = "0";
	iframeAxisMovements.width = "700";
	iframeAxisMovements.height = "700";
    //var xCoordDown = 0;
    //var yCoordDown = 0;
    //var zoomFlag = false;
	iframeAxisMovements.onload = function() {
		initDragAndDrop();
        initZoom(iframeAxisMovements, 'axis-movements-widget');
    }
 
   /* workspace.addEventListener('mousedown', function(e) {
        xCoordDown = e.offsetX;
        yCoordDown = e.offsetY;
        zoomFlag = true;
        iframeAxisMovements.style.cursor = "cell"
        //alert("Down");
    });
    workspace.addEventListener('mouseup', function(e) {
        if (zoomFlag == true) {
        iframeAxisMovements.style.cursor = "default"
        xdiff =  xCoordDown - e.offsetX;
        ydiff = yCoordDown - e.offsetY;
        distanse = Math.sqrt(Math.pow(xdiff, 2) + Math.pow(ydiff, 2));
        zoom =  distanse/700;
        alert(zoom);
        zoomFlag = false;
    }
    });*/
   // workspace.appendChild(wrapper);
	workspace.appendChild(iframeAxisMovements);	
}


function addStreamWidget() {
	iframeStream = document.createElement('iframe');
	iframeStream.src = "stream-widget.html" ;
	iframeStream.width = "300";
	iframeStream.height = "250";
	iframeStream.onload = function(){
		initDragAndDrop();
	}
	workspace.appendChild(iframeStream);	
}

function addCoordInfoWidget() {
	iframeInfo = document.createElement('iframe');
	iframeInfo.src = "coordinates-info.html";
    iframeInfo.frameBorder = "0";
	iframeInfo.width = "600";
	iframeInfo.height = "600";
	iframeInfo.onload = function(){
		initDragAndDrop();
        initZoom(iframeInfo, 'coord-info-widget');
	}
	workspace.appendChild(iframeInfo);	  
}

function initGestures() {
	var workspace = document.getElementById('workspace');
	var points = [];
	var wait=0;
	workspace.addEventListener('touchmove', function(event) {
	  // If there's exactly three fingers inside this element
	  if (event.targetTouches.length == 3) {
		  points.push(event.targetTouches[0]);
		  wait+=1;
		  if ( wait > 0 && menu.isshowed == false) {
			  //len = points.length;
			  /*vector1 = Math.sqrt(Math.pow(points[len-1].pageX-points[0].pageX, 2) + Math.pow(points[len-1].pageY-points[0].pageY, 2));
			  vector2 = Math.sqrt(Math.pow(points[len-1].pageX-points[0].pageX, 2));
			  cos = vector2/vector1;
			  alert(cos);*/
			  	menuShow()
		  }
	  }
	  
	}, false);
}

/**
 * Dragstart event handler. The function counts offsets, gets element id and joins them in one string.
 *
 * @param event event object.
 */
function dragStart(event) {
    var style = window.getComputedStyle(event.target, null);
	widgetWidth = document.getElementsByTagName('iframe')[0].getAttribute('width');
    event.dataTransfer.setData("text/plain", (parseInt(style.getPropertyValue("left"),10) - widgetWidth - 50 - event.clientX) + ','
		+ (parseInt(style.getPropertyValue("top"),10) - event.clientY) + ',' + event.target.elementId);
}

/**
Implementation for mobile devicies.
*/
function dragStartMobile() {
}

/**
 * Dragover event handler. The function cancels the event and allows to drop element.
 *
 * @param event event object.
 */
function dragOver(event) {
    event.preventDefault();
    return false;
}

/**
 * Drop event handler. The function reads data from dataTransfer object and places the element.
 *
 * @param event event object.
 * @return false.
 */
function drop(event) {
    var data = event.dataTransfer.getData("text/plain").split(',');
    var drag_element = document.getElementsByTagName('iframe')[0];
    drag_element.style.left = (event.clientX + parseInt(data[0],10)) + 'px';
    drag_element.style.top = (event.clientY + parseInt(data[1],10)) + 'px';
    event.preventDefault();
    return false;
}

/**
Implementation for mobile devicies.
*/
function dropMobile() {
}

/**
 * The method for initializing listeners for drag-and-drop events.
 *
 */
function initDragAndDrop(drag_elements) {
	var drag_elements = []
	drag_elements.push(window.frames[window.frames.length-1].document.getElementsByClassName('drag-drop')[0]);
    for (i = 0; i < drag_elements.length; ++i) {
        drag_elements[i].elementId = i;
        drag_elements[i].addEventListener('dragstart', dragStart, false);
        document.body.addEventListener('dragover', dragOver, false);
        document.body.addEventListener('drop', drop, false);
    }
}

//TODO: Поправить координаты для Drag&Drop.
//TODO: Реализовать Drad&Drop для мобильных устройств.


