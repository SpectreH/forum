function InitMainPage() {
  alertType = document.getElementById("alert").getAttribute("type")
  alertMessage = document.getElementById("alert-text").innerHTML

  if (alertType != "") {
    GenerateAlertBox(alertType, alertMessage)
  }

  SetCategories()
  ShowExtraFilters()
}

function ShowExtraFilters() {
  if (document.getElementById("logged-in") != null) {
    document.getElementById("sort-extra").style.display = "grid";
  } else {
    document.getElementById("sort-extra").remove();
  }
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
  a.setAttribute('onclick', "ShowAllCategories(this); return false;");
  a.setAttribute('href', "");

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