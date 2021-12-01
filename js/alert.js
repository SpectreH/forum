let ALERT_TIMER;
let TIME_OUT;
function GenerateAlertBox(type, message) {
  element = document.getElementsByTagName("div").namedItem("alert");

  if (element.style.display == "flex") {
    clearInterval(ALERT_TIMER);
    clearTimeout(TIME_OUT);
  }

  if (type != "") {
    element.style.opacity = 1;
    element = document.getElementsByTagName("div").namedItem("alert");
    document.getElementsByTagName("h1").namedItem("alert-text").innerHTML = message;
    element.style.display = "flex";

    if (type.toLowerCase() != "success") {
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