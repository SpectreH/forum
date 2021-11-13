function SetLiked(element) {
  document.getElementById("dislike")

  if (element.value != "liked") {
    element.setAttribute("value", "liked");
    element.setAttribute("style", "background-color: #87e98a;");

    var counter = element.childNodes[3].innerHTML;
    element.childNodes[3].innerHTML = parseInt(counter, 10) + 1;

    ClearRating(document.getElementById("dislike"));
  } else {
    ClearRating(element);
  }
}

function SetDisLiked(element) {
  if (element.value != "disliked") {
    element.setAttribute("value", "disliked");
    element.setAttribute("style", "background-color: #f74c4c;");

    var counter = element.childNodes[3].innerHTML;
    element.childNodes[3].innerHTML = parseInt(counter, 10) + 1;

    ClearRating(document.getElementById("like"));
  } else {
    ClearRating(element);
  }
}

function ClearRating(element) {
  if (element.value != "") {
    element.childNodes[3].innerHTML -= 1;
  }

  element.setAttribute("value", "");
  element.setAttribute("style", "background-color: none;");
}

function postfunction() {
  var ajax = new XMLHttpRequest();
  ajax.open("POST", "/1", true);
  ajax.send("test");
}