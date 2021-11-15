let ALL_CATEGORIES = [];

function CheckFileSize(element) {
  if (element.files.length > 0) {
    const fileSize = element.files.item(0).size;
    if (fileSize < 2048000) {
      document.getElementById("imageName").innerHTML = element.files.item(0).name
    } else {
      GenerateAlertBox("FileSize", "File too big, please select a file less than 2 mb");
      element.value = null;
      document.getElementById("imageName").innerHTML = ""
    }
  }
}

function SubmitForm() {
  var mainForm = document.getElementById("mainForm");
  var categorieForm = document.getElementById("allNewCategories");

  for (i = 0; i < mainForm.length; i++) {
    if (mainForm[i].checkValidity() == false) {
      mainForm[i].reportValidity();
      return
    }
  }

  if (categorieForm.childElementCount < 1) {
    GenerateAlertBox("CategoryCounter", "Require at least 1 category!");
    return
  }

  mainForm.submit();
}

function DeleteNewCategorie(element) {
  for (i = 0; i < ALL_CATEGORIES.length; i++) {
    if (element.getAttribute("value") == ALL_CATEGORIES[i]) {
      ALL_CATEGORIES.splice(i);
    }
  }

  element.parentElement.remove();
}

function AddCategorie() {
  input = document.getElementById("newCategorie");

  if (input.value.length < 3) {
    GenerateAlertBox("CategorySize", "New category length must be at least 3 character long!");
    return;
  }
  
  if (CheckCategoryDublicates(input.value)) {
    GenerateAlertBox("CategoryDublicate", "You are already added this category!");
    return;
  }

  ALL_CATEGORIES.push(input.value);

  var div = document.createElement("div");
  var a = document.createElement("a");
  var span = document.createElement("span");
  var node = document.createTextNode(input.value);

  div.classList.add("category-box");

  a.classList.add("category-box-content");
  a.setAttribute('type', "category");
  a.setAttribute('rel', "nofollow");
  a.setAttribute('href', "#");
  a.setAttribute('value', input.value);
  a.setAttribute('onclick', "DeleteNewCategorie(this)");

  span.classList.add("close-categorie");

  span.appendChild(document.createTextNode("X"));
  a.appendChild(node);

  a.appendChild(span);
  div.appendChild(a);

  var parentElement = document.getElementById("allNewCategories");
  parentElement.appendChild(div);

  input.value = "";
}

function CheckCategoryDublicates(categorieName) {
  for (i = 0; i < ALL_CATEGORIES.length; i++) {
    if (categorieName == ALL_CATEGORIES[i]) {
      return true;
    }
  }
  return false;
}