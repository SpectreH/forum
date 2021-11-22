function InitMainPage(loggedIn, messageType, alertMessage) {
  GenerateAlertBox(messageType, alertMessage)
  SetCategories()

  InitButtons(loggedIn)
}

function InitButtons(loggedIn) {
  document.getElementById("header-list").style.display = "flex";
  if (loggedIn == "true") {
    ChangeButtons()
  }
}

function ChangeButtons() {
  document.getElementsByTagName("span").namedItem("login/account").innerHTML = "Account"
  document.getElementById("login-link").href = "/account"
  document.getElementsByTagName("span").namedItem("registration").innerHTML = "New Post"
  document.getElementById("registration-link").href = "/new"

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

  const node = document.createTextNode("Log Out");
  span.appendChild(node);
  a.appendChild(span)
  li.appendChild(a)

  const headerListElement = document.getElementById("header-list");
  headerListElement.appendChild(li)
}

function SetCategories() {
  var categoriesSectionsLenght = document.getElementsByClassName("post-entry-categories").length;
  var categoriesSections = document.querySelectorAll("div.post-entry-categories");

  for (i = 0; i < categoriesSectionsLenght; i++) {
    var categories = categoriesSections[i].getElementsByTagName("div");
    var showCategoriesCounter = categories.length;
    if (categories.length > 3) {
      showCategoriesCounter = 3;
      CreateShowMoreCatButton(categoriesSections[i], categories.length - 3);
    }
    
    for (k = 0; k < showCategoriesCounter; k++) {
      categories[k].style.display = "block"
    }
  }
}

function CreateShowMoreCatButton(element, lenght) {
  var div = document.createElement("div");
  var a = document.createElement("a");
  var node = document.createTextNode(lenght + "+");

  div.classList.add("category-box");
  div.setAttribute('type', "show-more");

  a.classList.add("category-box-content");
  a.setAttribute('type', "category");
  a.setAttribute('rel', "nofollow");
  a.setAttribute('onclick', "ShowAllCategories(this)");
  a.setAttribute('href', "#");

  a.appendChild(node);
  div.appendChild(a);
  element.appendChild(div);
}

function ShowAllCategories(element) {
  var categoriesSectionLenght = element.parentElement.parentElement.getElementsByTagName("div").length;
  var categoriesSection = element.parentElement.parentElement.getElementsByTagName("div");
  for (i = 0; i < categoriesSectionLenght; i++) {
    categoriesSection[i].style.display = "block";
  }
  element.parentElement.remove();
}


let ALERT_TIMER;
let TIME_OUT;
function GenerateAlertBox(messageType, alertMessage) {
  element = document.getElementsByTagName("div").namedItem("alert");

  if (element.style.display == "flex") {
    clearInterval(ALERT_TIMER);
    clearTimeout(TIME_OUT);
  }

  if (messageType != "") {
    element.style.opacity = 1;
    element = document.getElementsByTagName("div").namedItem("alert");
    document.getElementsByTagName("h1").namedItem("alert-text").innerHTML = alertMessage;
    element.style.display = "flex";

    if (messageType != "Login" && messageType != "Register" && messageType != "Logout" && messageType != "newPost") {
      document.getElementsByTagName("h1").namedItem("alert-text").style.backgroundColor = "rgba(211, 10, 10, 0.85)";
      document.getElementsByTagName("h1").namedItem("alert-text").style.border = "2px solid rgba(211, 10, 10, 0.85)";
    }

    TIME_OUT = setTimeout(function () { AlerFadeOut(element); }, 3000);
  }
}

function ClearAlertMessage(element) {
  clearInterval(ALERT_TIMER);
  clearTimeout(TIME_OUT);
  element.style.display = "none";
}

function AlerFadeOut(element) {
  var op = 1;  // initial opacity
  ALERT_TIMER = setInterval(function () {

    if (op <= 0.01) {
      clearInterval(ALERT_TIMER);
      element.style.display = "none";
    }
    element.style.opacity = op;
    op -= 0.01;
  }, 25);
}