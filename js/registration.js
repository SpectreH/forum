function ReturnData() {
  if (nameError == "true" || emailError == "true") {
    document.getElementById("username").value = username
    document.getElementById("email").value = email
  }

  if (nameError == "true") {
    document.getElementsByTagName("span").namedItem("name").innerHTML = "Username is already taken"
    document.getElementsByTagName("input").namedItem("username").style.boxShadow = "0 0 6px #d45252"
    document.getElementsByTagName("input").namedItem("username").style.borderColor = "#d45252"
  }

  if (emailError == "true") {
    console.log(document.getElementsByClassName("error-notification").value)
    document.getElementsByTagName("span").namedItem("email").innerHTML = "Email is already registered"
    document.getElementsByTagName("input").namedItem("email").style.boxShadow = "0 0 6px #d45252"
    document.getElementsByTagName("input").namedItem("email").style.borderColor = "#d45252"
  }
}