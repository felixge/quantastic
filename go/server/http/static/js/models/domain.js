(function(ns) {
  var domain = ns.model.domain = {};

  model.TimeEntry = function(properties) {
    this.start = properties.start;
    this.end = properties.end;
    this.category = properties.category;
  };

  model.TimeCategory = function(properties) {
    this.name = properties.name;
    this.parent = properties.parent;
    this.children = properties.children || [];
  };
})(Quantastic);
