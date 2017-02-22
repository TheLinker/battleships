"use strict";

var MyWebSocket = function(host) {
    this.host = host || document.location.host;

    this.callbacks = {
        "close": function(msg) {console.log(msg);},
        "error": function(msg) {console.log(msg);}
    };

    this.connect();
}

MyWebSocket.prototype.connect = function() {
    if (window["WebSocket"]) {
        var that = this;
        try {
            this.conn = new WebSocket("ws://" + this.host + "/ws");
        } catch(e) {
            console.log(e);
            return alse;
        }
        this.conn.onclose = function (evt) {
            that.callbacks['close']("Connection Closed");
        };
        this.conn.onmessage = function (evt) {
            var message = evt.data;
            var msg = JSON.parse(message);

            console.log(msg);

            if(typeof msg.Type !== 'undefined' && typeof that.callbacks[msg.Type] !== 'undefined')
                that.callbacks[msg.Type](msg.Msg);
        };
        this.conn.onerror = function (evt) {
            that.callbacks['error'](evt.message);
        };
    } else {
        throw("WebSocket not supported");
    }
}

MyWebSocket.prototype.close = function(code) {
    this.conn.close(code);
    this.conn = null;
}

MyWebSocket.prototype.registerCallback = function(type, func) {
    if (typeof func !== "function") throw("func not a function");
    this.callbacks[type] = func;

    return this;
}

MyWebSocket.prototype.removeCallback = function(type) {
    delete this.callbacks[type];

    return this;
}

MyWebSocket.prototype.sendData = function(type, msg) {
    try {
        if (this.conn && this.conn.readyState === 1) {
            this.conn.send(JSON.stringify({
                "type": type,
                "msg": msg
            }));

            return true;
        }

        throw("Disconnected from the WS server");
    } catch(e) {
        console.log(e);
        return false;
    }

}