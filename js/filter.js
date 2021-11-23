let categoriesPosts = [];
let likedPosts = [];
let dislikedPosts = [];
let createdPosts = [];
let allFiltersArr = [categoriesPosts, likedPosts, dislikedPosts, createdPosts];

function ShowPostsByFilter() {
  ConnectArrays()
  $("#posts article").hide();

  if (categoriesPosts.length == 0 && likedPosts.length == 0 && dislikedPosts.length == 0 && createdPosts.length == 0 && ($('input:checked').length) == 0) {
    $("#posts article").show();
  } else {
    for (i = 0; i < allFiltersArr.length; i++) {
      $(allFiltersArr[i]).show();
    }
  }
}

$(document).ready(function () {
  // Category filter section
  $('div.sort-categories').delegate('input:checkbox', 'change', function () {
    categoriesPosts = ClearArr();

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
        AppendElementToArr(categoriesPosts, elem)
      }
    });

    ShowPostsByFilter();
  }).find('input:checkbox').change();

  // Extra filter section
  $('div.sort-extra').delegate('input:checkbox', 'change', function () {
    if ($(this).is(":checked")) {
      if ($(this).attr("value") == "liked") {
        FindPosts("liked", true, likedPosts)
      } else if ($(this).attr("value") == "disliked") {
        FindPosts("disliked", true, dislikedPosts)
      } else {
        FindPosts("author", username, createdPosts)
      }
    } else {
      if ($(this).attr("value") == "liked") {
        likedPosts = ClearArr();
      } else if ($(this).attr("value") == "disliked") {
        dislikedPosts = ClearArr();
      } else {
        createdPosts = ClearArr();
      }
    }
  
    ShowPostsByFilter();
  }).find('input:checkbox').change();
});

function FindPosts(data, value, arr) {
  $('.content > article').each(function (i, elem) {
    if ($(elem).data(data) == value) {
      AppendElementToArr(arr, elem)
    }
  });
}

function ClearArr() {
  return [];
}

function ConnectArrays() {
  allFiltersArr = [categoriesPosts, likedPosts, dislikedPosts, createdPosts];
}

function AppendElementToArr(arr, element) {
  for (i = 0; i < arr.length; i++) {
    if (arr[i] == element) {
      return
    }
  }
  arr.push(element)
}