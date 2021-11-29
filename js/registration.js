$(document).ready(function () {
  if (document.getElementById("email-error") != null) {
    document.getElementsByTagName("input").namedItem("email").style.boxShadow = "0 0 6px #d45252"
    document.getElementsByTagName("input").namedItem("email").style.borderColor = "#d45252"
  }
    
  if (document.getElementById("name-error") != null) {
    document.getElementsByTagName("input").namedItem("username").style.boxShadow = "0 0 6px #d45252"
    document.getElementsByTagName("input").namedItem("username").style.borderColor = "#d45252"
  }
});