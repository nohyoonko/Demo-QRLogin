/* QR Code 생성 -> text(랜덤값) 서버에서 생성해서 가지고 오기 */
var qrcode = new QRCode(document.getElementById("qrcode"), {
    text: "01234",
    width: 150,
    height: 150,
    colorDark: "#000000",
    colorLight: "#ffffff",
    correctLevel: QRCode.CorrectLevel.H
});

$("#qrcode > img").css({"margin": "5% auto"});

/* Modal Popup */
  var modal = document.getElementById('myModal');
  var btn = document.getElementById("myBtn");
  var span = document.getElementsByClassName("close")[0];                                          

  btn.onclick = function() {
      modal.style.display = "block";
  }
  span.onclick = function() {
      modal.style.display = "none";
  }
  // When the user clicks anywhere outside of the modal, close it
  window.onclick = function(event) {
      if (event.target == modal) {
          modal.style.display = "none";
      }
  }
