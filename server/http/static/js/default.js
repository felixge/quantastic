function visualizeTime() {
  var list = d3.select('.time ol');
  var data = list.selectAll('li')[0].map(function(item) {
    return {
      duration: item.getAttribute('data-duration'),
    }
  });
  console.log(data);
  list.remove();
}

//visualizeTime();
