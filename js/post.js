let LOGGEDIN
let POSTID

function InitPostPage() {
  LOGGEDIN = document.getElementById("logged-in")
  POSTID = document.getElementById("post-rate").getAttribute("post-id")

  rateButton = document.getElementsByClassName("rate-button")

  for (i = 0; i < rateButton.length; i++) {
    if (rateButton[i].getAttribute("value") == "true" && rateButton[i].getAttribute("type") == "like") {
      rateButton[i].setAttribute("style", "background-color: #87e98a;");
    }

    if (rateButton[i].getAttribute("value") == "true" && rateButton[i].getAttribute("type") == "dislike") {
      rateButton[i].setAttribute("style", "background-color: #f74c4c;");
    }
  }

  if (LOGGEDIN == null) {
    document.getElementById("new-comment").remove();
    document.getElementById("login-message").style.display = "block"
  } else {
    document.getElementById("new-comment").style.display = "grid"
    document.getElementById("login-message").remove();
  }
}

let code
function SetLiked(element, type) {
  if (LOGGEDIN == null) {
    GenerateAlertBox("Fail_NotLoggedIn", "Please login to rate the post!")
    return
  }

  var mirrorButton, id
  if (type == "comment") {
    id = element.getAttribute("comment-id")
    mirrorButton = document.getElementById(id + "-dislike")
  } else {
    id = POSTID
    mirrorButton = document.getElementById("post-dislike")
  }

  if (element.value == "false") {
    code = 1
    element.setAttribute("value", "true");
    element.setAttribute("style", "background-color: #87e98a;");

    var counter = element.childNodes[3].innerHTML;
    element.childNodes[3].innerHTML = parseInt(counter, 10) + 1;

    ClearRating(mirrorButton);
  } else {
    ClearRating(element);
    code = 2
  }

  SendPostRequest(code, type, id)
}

function SetDisLiked(element, type) {
  if (LOGGEDIN == null) {
    GenerateAlertBox("Fail_NotLoggedIn", "Please login to rate the post!")
    return
  }

  var mirrorButton, id
  if (type == "comment") {
    id = element.getAttribute("comment-id")
    mirrorButton = document.getElementById(id + "-like")
  } else {
    id = POSTID
    mirrorButton = document.getElementById("post-like")
  }

  if (element.value == "false") {
    code = -1
    element.setAttribute("value", "true");
    element.setAttribute("style", "background-color: #f74c4c;");

    var counter = element.childNodes[3].innerHTML;
    element.childNodes[3].innerHTML = parseInt(counter, 10) + 1;

    ClearRating(mirrorButton);
  } else {
    ClearRating(element);
    code = -2
  }

  SendPostRequest(code, type, id)
}

function SubmitForm() {
  if (LOGGEDIN == null) {
    GenerateAlertBox("Fail_NotLoggedIn", "Please login to add comment!")
    return
  }

  var form = document.getElementById("commentForm")

  for (i = 0; i < form.length; i++) {
    if (form[i].checkValidity() == false) {
      form[i].reportValidity();
      return
    }
  }

  form.submit();
}

function ClearRating(element) {
  if (element.value == "true") {
    element.childNodes[3].innerHTML -= 1;
  }

  element.setAttribute("value", "false");
  element.setAttribute("style", "background-color: none;");
}

function SendPostRequest(code, type, id) {
  var url = "/" + POSTID;
  var data = type + ";" + code + ";" + id;
  var ajax = new XMLHttpRequest();
  ajax.open("POST", url, true);
  ajax.send(data);
}