(function(model) {
  var orange = new model.Color(0xFF, 0x80, 0x00, 1);
  var blue = new model.Color(0x00, 0x80, 0xFF, 1);
  var grey = new model.Color(0xDA, 0xDA, 0xDA, 1);

  var sleep = new model.TimeCategory({
    name: 'Sleep',
    color: orange,
  });
  var work = new model.TimeCategory({
    name: 'Work',
    color: blue,
  });
  var untracked = new model.TimeCategory({
    name: 'Untracked',
    color: grey,
  });

  var entries = [
    new model.TimeEntry({
      start: new Date('2014-01-06T22:00:00Z'),
      end: new Date('2014-01-07T07:00:00Z'),
      category: sleep,
    }),
    new model.TimeEntry({
      start: new Date('2014-01-07T07:00:01Z'),
      end: new Date('2014-01-07T13:00:00Z'),
      category: work,
    }),
    new model.TimeEntry({
      start: new Date('2014-01-07T13:00:00Z'),
      end: new Date('2014-01-07T22:00:00Z'),
      category: untracked,
    }),
  ];

  var width = 960;
  var height = 100;
  var barHeight = 20;

  var x = d3.time.scale().range([0, width]);

  var dates = function(entries) {
    return entries.reduce(function(results, entry) {
      return results.concat(entry.start, entry.end);
    }, [])
  };

  var color = function(entry) {
    var c = entry.category.color;
    return 'rgba('+c.r+','+c.g+','+c.b+','+c.a+')'; 
  };

  var svg = d3.select('svg')
    .attr('width', width)
    .attr('height', height);

  var brush = d3.svg.brush().x(x).on('brush', brushed);

  function update(entries) {
    var extent = d3.extent(dates(entries));
    x.domain(extent);
    var bar = svg.selectAll("g")
        .data(entries)
      .enter().append("g")
        .attr('transform', function(d, i) { return "translate("+x(d.start)+",0)"; });
    bar.append('rect')
        .attr('fill', color)
        .attr('width', function(d) {
          return x(d.end)-x(d.start);
        })
        .attr('height', barHeight);
  }

  update(entries);

  var xAxis = d3.svg.axis()
      .scale(x)
      .orient("bottom")
      .tickFormat(d3.time.format("%H:%M"));

  svg
    .append("g")
    .attr("class", "x axis")
    .attr("transform", "translate(0," + barHeight + ")")
    .call(xAxis);

  svg
    .append("g")
    .attr("class", "x brush")
    .call(brush)
    .selectAll("rect")
    .attr("y", -6)
    .attr("height", barHeight + 7);

    function brushed() {
      var min5 = 5 * 60 * 1000;
      var extent = brush.extent();
      var newExtent = extent.map(function(start) {
        return new Date(Math.round(start.getTime()/min5)*min5);
      });

      d3.select(this).call(brush.extent(newExtent));
      console.log(newExtent);
    }



  function old() {
    var width = 960;
    var height = 80;
    var margin2 = {top: 10, right: 0, bottom: 20, left: 0};

    var parseDate = d3.time.format("%Y-%m-%d %H:%M").parse;

    var x2 = d3.time.scale().range([0, width]);
    var y2 = d3.scale.linear().range([height, 0]);

    var xAxis2 = d3.svg.axis()
        .scale(x2)
        .orient("bottom")
        .tickFormat(d3.time.format("%H:%M"));

    var brush = d3.svg.brush()
        .x(x2)
        .on("brush", brushed);

    brush.clamp(true);

    var area = d3.svg.area()
        .interpolate("monotone")
        .x(function(d) { return x2(d.start); })
        .y0(height)
        .y1(function(d) { return y2(1); });

    var svg = d3.select("#d3")
        .append("svg")
        .attr("width", width)
        .attr("height", height +50);

    svg
      .append("defs")
      .append("clipPath")
      .attr("id", "clip")
      .append("rect")
      .attr("width", width)
      .attr("height", height);


    var context = svg.append("g")
        .attr("transform", "translate(" + margin2.left + "," + margin2.top + ")");


      x2.domain(d3.extent(entries.map(function(d) { return d.start; })));
      y2.domain([0, d3.max(entries.map(function(d) { return 1; }))]);

      context
        .append("path")
        .datum(entries)
        .attr("d", area);

      context
        .append("g")
        .attr("class", "x axis")
        .attr("transform", "translate(0," + height + ")")
        .call(xAxis2);

      context
        .append("g")
        .attr("class", "x brush")
        .call(brush)
        .selectAll("rect")
        .attr("y", -6)
        .attr("height", height + 7);

    function brushed() {
      var min5 = 5 * 60 * 1000;
      var extent = brush.extent();
      var newExtent = extent.map(function(start) {
        return new Date(Math.round(start.getTime()/min5)*min5);
      });

      d3.select(this).call(brush.extent(newExtent));
      console.log(newExtent);
    }
  }
})(Quantastic.model);
