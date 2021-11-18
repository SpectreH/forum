function InitPostPage(liked, disliked) {
  if (liked == "true") {
    var element = document.getElementById("like");
    element.setAttribute("value", "liked");
    element.setAttribute("style", "background-color: #87e98a;");
  } else if (disliked == "true") {
    var element = document.getElementById("dislike");
    element.setAttribute("value", "disliked");
    element.setAttribute("style", "background-color: #f74c4c;");
  }
}


let code
function SetLiked(loggedIn, element) {
  if (loggedIn == "false") {
    GenerateAlertBox("NotLoggedIn", "Please log in to rate the post!")
    return
  }

  if (element.value != "liked") {
    code = 1
    element.setAttribute("value", "liked");
    element.setAttribute("style", "background-color: #87e98a;");

    var counter = element.childNodes[3].innerHTML;
    element.childNodes[3].innerHTML = parseInt(counter, 10) + 1;

    ClearRating(document.getElementById("dislike"));
  } else {
    ClearRating(element);
    code = 2
  }

  postfunction(code)
}

function SetDisLiked(loggedIn, element) {
  if (loggedIn == "false") {
    GenerateAlertBox("NotLoggedIn", "Please log in to rate the post!")
    return
  }

  if (element.value != "disliked") {
    code = -1
    element.setAttribute("value", "disliked");
    element.setAttribute("style", "background-color: #f74c4c;");

    var counter = element.childNodes[3].innerHTML;
    element.childNodes[3].innerHTML = parseInt(counter, 10) + 1;

    ClearRating(document.getElementById("like"));
  } else {
    ClearRating(element);
    code = -2
  }

  postfunction(code)
}

function ClearRating(element) {
  if (element.value != "") {
    element.childNodes[3].innerHTML -= 1;
  }

  element.setAttribute("value", "");
  element.setAttribute("style", "background-color: none;");
}

function postfunction(code) {
  var ajax = new XMLHttpRequest();
  ajax.open("POST", "/1", true);
  ajax.send(code);
}