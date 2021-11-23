function InitPostPage() {
  if (liked == "true") {
    var element = document.getElementById("like");
    element.setAttribute("value", "liked");
    element.setAttribute("style", "background-color: #87e98a;");
  } else if (disLiked == "true") {
    var element = document.getElementById("dislike");
    element.setAttribute("value", "disliked");
    element.setAttribute("style", "background-color: #f74c4c;");
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
function SetLiked(element) {
  if (loggedIn == "false") {
    GenerateAlertBox("NotLoggedIn", "Please login to rate the post!")
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

  SendPostRequest(code)
}

function SetDisLiked(element) {
  if (loggedIn == "false") {
    GenerateAlertBox("NotLoggedIn", "Please login to rate the post!")
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

  SendPostRequest(code, postId)
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

function SendPostRequest(code) {
  var id = "/" + postId
  var ajax = new XMLHttpRequest();
  ajax.open("POST", id, true);
  ajax.send(code);
}