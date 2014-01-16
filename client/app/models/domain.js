var Date = (function () {
    function Date(year, month, day) {
        this.year = year;
        this.month = month;
        this.day = day;
    }
    Date.prototype.getYear = function () {
        return this.year;
    };
    Date.prototype.getMonth = function () {
        return this.month;
    };
    Date.prototype.getDay = function () {
        return this.day;
    };
    return Date;
})();
exports.Date = Date;

var Time = (function () {
    function Time(hour, minute, second) {
        this.hour = hour;
        this.minute = minute;
        this.second = second;
    }
    Time.prototype.getHour = function () {
        return this.hour;
    };
    Time.prototype.getMinute = function () {
        return this.minute;
    };
    Time.prototype.getSecond = function () {
        return this.second;
    };
    return Time;
})();
exports.Time = Time;

var DateTime = (function () {
    function DateTime(date, time) {
        this.date = date;
        this.time = time;
    }
    DateTime.prototype.getDate = function () {
        return this.date;
    };
    DateTime.prototype.getTime = function () {
        return this.time;
    };
    return DateTime;
})();
exports.DateTime = DateTime;

// @TODO Enforce start > end
var TimeRange = (function () {
    function TimeRange(start, end) {
        this.start = start;
        this.end = end;
    }
    TimeRange.prototype.getStart = function () {
        return this.start;
    };
    TimeRange.prototype.getEnd = function () {
        return this.end;
    };
    return TimeRange;
})();
exports.TimeRange = TimeRange;

var TimeEntry = (function () {
    function TimeEntry() {
    }
    return TimeEntry;
})();
exports.TimeEntry = TimeEntry;
