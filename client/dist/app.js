(function e(t,n,r){function s(o,u){if(!n[o]){if(!t[o]){var a=typeof require=="function"&&require;if(!u&&a)return a(o,!0);if(i)return i(o,!0);throw new Error("Cannot find module '"+o+"'")}var f=n[o]={exports:{}};t[o][0].call(f.exports,function(e){var n=t[o][1][e];return s(n?n:e)},f,f.exports,e,t,n,r)}return n[o].exports}var i=typeof require=="function"&&require;for(var o=0;o<r.length;o++)s(r[o]);return s})({1:[function(require,module,exports){
var view = require('./models/view');

var red = new view.Color();
red.setR(255);
console.log(red.getHex());

console.log('yes');

},{"./models/view":2}],2:[function(require,module,exports){
var Timeline = (function () {
    function Timeline() {
    }
    return Timeline;
})();
exports.Timeline = Timeline;

var Color = (function () {
    function Color() {
        this.r = 0;
        this.g = 0;
        this.b = 0;
    }
    Color.prototype.getR = function () {
        return this.r;
    };
    Color.prototype.setR = function (val) {
        return this.r = val;
    };
    Color.prototype.getG = function () {
        return this.r;
    };
    Color.prototype.setG = function (val) {
        return this.g = val;
    };
    Color.prototype.getB = function () {
        return this.b;
    };
    Color.prototype.setB = function (val) {
        return this.b = val;
    };
    Color.prototype.getHex = function () {
        return '#' + toHex(this.r) + toHex(this.g) + toHex(this.b);
    };
    return Color;
})();
exports.Color = Color;

var toHex = function (num) {
    var hex = num.toString(16);
    if (hex.length == 1) {
        hex = '0' + hex;
    }
    return hex;
};

},{}]},{},[1])