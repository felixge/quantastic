export class Date{
  constructor(private year: number,
              private month: number,
              private day: number) {}
  getYear() {return this.year;}
  getMonth() {return this.month;}
  getDay() {return this.day;}
}

export class Time{
  constructor(private hour: number,
              private minute: number,
              private second: number) {}
  getHour() {return this.hour;}
  getMinute() {return this.minute;}
  getSecond() {return this.second;}
}

export class DateTime{
  constructor(private date: Date, private time: Time) {}
  getDate() {return this.date;}
  getTime() {return this.time;}
}

// @TODO Enforce start > end
export class TimeRange {
  constructor(private start: DateTime, private end: DateTime){}
  getStart() {return this.start;}
  getEnd() {return this.end;}
  //getDuration(): Time { return this.end - this.start; }
}

export class TimeEntry{
  private ranges: TimeRange[];
}
