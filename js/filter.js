function ShowPostsByFilter() {
  var checkBoxFilterApplied = false
  $("#posts article").hide();

  $("#posts article").each(function () {
    if ($(this).attr("check-box-filter") == "true") {
      checkBoxFilterApplied = true
    }
  })

  if (!checkBoxFilterApplied) {
    $("#posts article").show();
  } else {
    $("#posts article").filter(function () {
      var tempVar = $(this).attr("check-box-filter");

      if (tempVar == "true") {
        return true
      } else {
        return false
      }
    }).show();
  }
}

$(document).ready(function () {
  $('div.sort-categories').delegate('input:checkbox', 'change', function () {
    var counter = 0

    $('input:checked').each(function () {
      counter = counter + 1
    });

    if (counter == 0) {
      $('.content > article').attr("check-box-filter", "false")
    }

    var selector = $('input:checked').map(function () {
      return $(this).attr('category');
    }).get();

    function FilterCategories(elem) {
      var elemCats = elem.attr('categories');

      if (elemCats) {
        elemCats = elemCats.split(';');
      } else {
        elemCats = Array();
      }

      for (var i = 0; i < selector.length; i++) {
        if (jQuery.inArray(selector[i], elemCats) != -1) {
          return true;
        }
      }

      return false;
    }
    $('.content > article').each(function (i, elem) {
      if (FilterCategories(jQuery(elem))) {
        $(elem).attr("check-box-filter", "true");
      } else {
        $(elem).attr("check-box-filter", "false");
      }
    });

    ShowPostsByFilter();

  }).find('input:checkbox').change();
});