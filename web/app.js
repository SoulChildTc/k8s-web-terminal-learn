if (!window.location.origin) { // Some browsers (mainly IE) do not have this property, so we need to build it manually...
    window.location.origin = window.location.protocol + '//' + window.location.hostname + (window.location.port ? (':' + window.location.port) : '');
}

const origin = window.location.origin;

// options usage example
const options = {
    debug: true,
    devel: true,
    protocols_whitelist: ['websocket', 'xdr-streaming', 'xhr-streaming', 'iframe-eventsource', 'iframe-htmlfile', 'xdr-polling', 'xhr-polling', 'iframe-xhr-polling', 'jsonp-polling']
};

let SOCK;
let SESSIONID;

function InitSockJs() {
    SOCK = new SockJS(origin+'/api/sockjs?');
    SOCK.onopen = function() {
        //console.log('connection open');
        document.getElementById("status").innerHTML = "connected";
        document.getElementById("send").disabled=false;
        SOCK.send('{"Op":"bind","SessionID":"' + SESSIONID + '"}')
    };

    SOCK.onmessage = function(e) {
        data = JSON.parse(e.data)
        console.log(data.Data)
        document.getElementById("output").value += data.Data +"\n";
    };

    SOCK.onclose = function() {
        document.getElementById("status").innerHTML = "connection closed";
        //console.log('connection closed');
    };
}


function InitSession() {
    const httpRequest = new XMLHttpRequest();
    httpRequest.open('GET', '/api/v1/pod/default/busybox-deployment-dcb89bc87-hrspk/shell/busybox?shell=sh', true);
    httpRequest.send();

    httpRequest.onreadystatechange = function () {
        if (httpRequest.readyState === 4 && httpRequest.status === 200) {
            const data = JSON.parse(httpRequest.responseText)
            console.log(data.data)
            document.getElementById("sessionID").innerHTML = data.data.id
            SESSIONID = data.data.id
        }
    };
}


InitSession()
InitSockJs()


function send(e) {

    e.preventDefault();
    const t = document.getElementById("input").value;
    const cmd = '{"Op":"stdin","Data":"' + t + '\\r","Cols":164,"Rows":41}'
    SOCK.send(cmd); return false;
}