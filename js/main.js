function InitMainPage(loggedIn, messageType, alertMessage) {
  GenerateAlertBox(messageType, alertMessage)

  if (loggedIn == "true") {
    ChangeButtons()
  }
}

function ChangeButtons() {
  document.getElementsByTagName("span").namedItem("login/account").innerHTML = "Account"
  document.getElementById("login-link").href = "/account"
  document.getElementById("registration-element").remove();

  CreateLogOutButton()
}

function CreateLogOutButton() {
  const li = document.createElement("li");
  li.classList.add("header-list-element");

  const a = document.createElement("a");
  a.classList.add("header-list-element-link");
  a.setAttribute('href',"/logout");
  a.setAttribute('name',"logout-box");
  a.setAttribute('type',"logout-box");

  const span = document.createElement("span");
  span.classList.add("header-list-element-link-border");
  span.setAttribute('name',"logout");
  span.setAttribute('type',"logout");

  const node = document.createTextNode("Log out");
  span.appendChild(node);
  a.appendChild(span)
  li.appendChild(a)

  const headerListElement = document.getElementById("header-list");
  headerListElement.appendChild(li)
}

function GenerateAlertBox(messageType, alertMessage) {
  if (messageType != "") {
    element = document.getElementsByTagName("div").namedItem("alert");
    document.getElementsByTagName("h1").namedItem("alert-text").innerHTML = alertMessage;
    element.style.display = "block";

    if (messageType != "Login" && messageType  != "Register" && messageType  != "Logout") {
      document.getElementsByTagName("h1").namedItem("alert-text").style.backgroundColor = "rgba(211, 10, 10, 0.85)";
      document.getElementsByTagName("h1").namedItem("alert-text").style.border = "2px solid rgba(211, 10, 10, 0.85)";
    }

    setTimeout(function () { AlerFadeOut(element); }, 3000);
  }
}

function AlerFadeOut(element) {
  var op = 1;  // initial opacity
  element.style.display = 'block';
  var timer = setInterval(function () {
    if (op <= 0.01) {
      clearInterval(timer);
      element.style.display = "none";
    }
    element.style.opacity = op;
    op -= 0.01;
  }, 25);
}