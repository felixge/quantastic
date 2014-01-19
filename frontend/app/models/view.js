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
