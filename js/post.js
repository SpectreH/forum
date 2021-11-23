function InitPostPage() {
  rateButtons = document.getElementsByClassName("rate-button")
  if (liked == "true") {
    document.getElementById("post-like").setAttribute("value", "liked")
  } else if (disLiked == "true") {
    document.getElementById("post-dislike").setAttribute("value", "disliked")
  }

  for (i = 0; i < rateButtons.length; i++) {
    var value = rateButtons[i].getAttribute("value")
    if (value == "liked") {
      rateButtons[i].setAttribute("style", "background-color: #87e98a;");
    } else if (value == "disliked") {
      rateButtons[i].setAttribute("style", "background-color: #f74c4c;");
    }
  }

  if (loggedIn == "false") {
    document.getElementById("new-comment").remove();
    document.getElementById("login-message").style.display = "block"
  } else {
    document.getElementById("new-comment").style.display = "grid"
    document.getElementById("login-message").remove();
  }
}

let code
function SetLiked(element, type) {
  if (loggedIn == "false") {
    GenerateAlertBox("NotLoggedIn", "Please login to rate the post!")
    return
  }

  var mirrorButton, id
  if (type == "comment") {
    id = element.getAttribute("comment-id")
    mirrorButton = document.getElementById(id + "-dislike")
  } else {
    id = postId
    mirrorButton = document.getElementById("post-dislike")
  }

  if (element.value != "liked") {
    code = 1
    element.setAttribute("value", "liked");
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
  if (loggedIn == "false") {
    GenerateAlertBox("NotLoggedIn", "Please login to rate the post!")
    return
  }

  var mirrorButton, id
  if (type == "comment") {
    id = element.getAttribute("comment-id")
    mirrorButton = document.getElementById(id + "-like")
  } else {
    id = postId
    mirrorButton = document.getElementById("post-like")
  }

  if (element.value != "disliked") {
    code = -1
    element.setAttribute("value", "disliked");
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
  console.log(loggedIn)

  if (loggedIn == "false") {
    GenerateAlertBox("NotLoggedIn", "Please login to add comment!")
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
  if (element.value != "") {
    element.childNodes[3].innerHTML -= 1;
  }

  element.setAttribute("value", "");
  element.setAttribute("style", "background-color: none;");
}

function SendPostRequest(code, type, id) {
  var url = "/" + postId;
  var data = type + ";" + code + ";" + id;
  var ajax = new XMLHttpRequest();
  ajax.open("POST", url, true);
  ajax.send(data);
}