import domain = require('./domain');

export class Timeline{

}


export class Color{
  private r = 0;
  private g = 0;
  private b = 0;

  getR() { return this.r; }
  setR(val: number) { return this.r = val; }
  getG() { return this.r; }
  setG(val: number) { return this.g = val; }
  getB() { return this.b; }
  setB(val: number) { return this.b = val; }
  getHex(): string {
    return '#'+toHex(this.r)+toHex(this.g)+toHex(this.b);
  }
}

var toHex = (num: number): string => {
  var hex = num.toString(16);
  if (hex.length == 1) {
    hex = '0'+hex;
  }
  return hex;
}
