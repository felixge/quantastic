module.exports = function() {
  var view = {};

  view.Color = function(r, g, b) {
    r = r || 0;
    g = g || 0;
    b = b || 0;

    var o = {
      getR: function() { return r; },
      setR: function(val) { r = val; },
      getG: function() { return g; },
      setG: function(val) { g = val; },
      getB: function() { return b; },
      setB: function(val) { b = val; },
    };

    return o.setR(r) || o.setG(g) || o.setB(b) || o;
  };

  view._validateColor = function(val) {
    if (typeof val != "number") {
      return '';
    }
  }

  return view;
};
