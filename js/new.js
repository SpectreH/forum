function CheckFileSize(element) {
  if (element.files.length > 0) {
    const fileSize = element.files.item(0).size;
    if (fileSize < 2048000) {
      document.getElementById("imageName").innerHTML = element.files.item(0).name
    } else {
      alert("File too big, please select a file less than 2 mb");
      element.value = null;
      document.getElementById("imageName").innerHTML = ""
    }
  }
}