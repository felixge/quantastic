import domain = require('./models/domain');
import view = require('./models/view');

var red = new view.Color();
red.setR(255);
console.log(red.getHex());

console.log('yes');
