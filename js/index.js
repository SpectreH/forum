function GenerateAlertBox(registrationSuccess, alertMessage) {
  if (registrationSuccess != "") {
    element = document.getElementsByTagName("div").namedItem("alert");
    document.getElementsByTagName("h1").namedItem("alert-text").innerHTML = alertMessage;
    element.style.display = "block";

    if (registrationSuccess != "Login" && registrationSuccess  != "Register") {
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