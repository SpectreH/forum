$(document).ready(function () {
  if (document.getElementById("login-error") != null) {
    document.getElementsByTagName("input").namedItem("login").style.boxShadow = "0 0 6px #d45252"
    document.getElementsByTagName("input").namedItem("login").style.borderColor = "#d45252"
  }

  if (document.getElementById("password-error") != null) {
    document.getElementsByTagName("input").namedItem("pass").style.boxShadow = "0 0 6px #d45252"
    document.getElementsByTagName("input").namedItem("pass").style.borderColor = "#d45252"
  }
});