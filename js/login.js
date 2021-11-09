function ReturnData(loginErr, passErr, login) {
  if (loginErr != "true") {
    document.getElementById("login").value = login
  }

  if (loginErr == "true") {
    document.getElementsByTagName("span").namedItem("login").innerHTML = "Account does not exist"
    document.getElementsByTagName("input").namedItem("login").style.boxShadow = "0 0 6px #d45252"
    document.getElementsByTagName("input").namedItem("login").style.borderColor = "#d45252"
  }

  if (passErr == "true") {
    console.log(document.getElementsByClassName("error-notification").value)
    document.getElementsByTagName("span").namedItem("password").innerHTML = "Password does not match"
    document.getElementsByTagName("input").namedItem("password").style.boxShadow = "0 0 6px #d45252"
    document.getElementsByTagName("input").namedItem("password").style.borderColor = "#d45252"
  }
}